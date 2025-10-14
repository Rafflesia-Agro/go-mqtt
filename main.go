package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/Anjasfedo/go-mqtt/src/broker"
	"github.com/Anjasfedo/go-mqtt/src/cache"
	"github.com/Anjasfedo/go-mqtt/src/config"
	"github.com/Anjasfedo/go-mqtt/src/database"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5" // BARU: Impor jwtauth
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// --- Structs for Hardware State Command (HTTP Input) ---
type TimeSchedule struct {
	Order   int    `json:"order"`
	TimeOn  string `json:"time_on"`
	TimeOff string `json:"time_off"`
}

type HardwareMeta struct {
	Mode           string         `json:"mode"`
	TemperatureMin *float64       `json:"temperature_min,omitempty"`
	TemperatureMax *float64       `json:"temperature_max,omitempty"`
	Schedules      []TimeSchedule `json:"schedules,omitempty"`
}

type HardwareState struct {
	Hardware string        `json:"hardware"`
	State    bool          `json:"state"`
	Meta     *HardwareMeta `json:"meta"`
}

// Partial update structs for flexible requests
type PartialHardwareMeta struct {
	Mode           *string          `json:"mode,omitempty"`
	TemperatureMin *float64         `json:"temperature_min,omitempty"`
	TemperatureMax *float64         `json:"temperature_max,omitempty"`
	Schedules      *[]TimeSchedule  `json:"schedules,omitempty"`
}

type PartialHardwareState struct {
	Hardware string                `json:"hardware"`
	State    *bool                 `json:"state,omitempty"`
	Meta     *PartialHardwareMeta  `json:"meta,omitempty"`
}

// HAPUS STRUCT INI: StatePayload tidak lagi digunakan.
// type StatePayload struct {
// 	HardwareStates []HardwareState `json:"hardware_state"`
// }

// mergePartialUpdates merges partial updates into an existing hardware state
func mergePartialUpdates(existing HardwareState, partial PartialHardwareState) HardwareState {
	result := existing

	// Update hardware if provided (though this should match the URL param)
	if partial.Hardware != "" {
		result.Hardware = partial.Hardware
	}

	// Update state if provided
	if partial.State != nil {
		result.State = *partial.State
	}

	// Handle meta updates
	if partial.Meta != nil {
		if result.Meta == nil {
			result.Meta = &HardwareMeta{}
		}

		// Update mode if provided
		if partial.Meta.Mode != nil {
			result.Meta.Mode = *partial.Meta.Mode
		}

		// Update temperature min if provided
		if partial.Meta.TemperatureMin != nil {
			result.Meta.TemperatureMin = partial.Meta.TemperatureMin
		}

		// Update temperature max if provided
		if partial.Meta.TemperatureMax != nil {
			result.Meta.TemperatureMax = partial.Meta.TemperatureMax
		}

		// Update schedules if provided
		if partial.Meta.Schedules != nil {
			result.Meta.Schedules = *partial.Meta.Schedules
		}
	}

	return result
}

func validateSchedules(schedules []TimeSchedule) error {
	if len(schedules) <= 1 {
		return nil
	}
	sort.Slice(schedules, func(i, j int) bool {
		return schedules[i].Order < schedules[j].Order
	})
	const timeLayout = "15:04"
	orderSeen := make(map[int]bool)
	var previousTimeOff time.Time
	for i, s := range schedules {
		if orderSeen[s.Order] {
			return fmt.Errorf("duplicate order number %d found", s.Order)
		}
		orderSeen[s.Order] = true
		timeOn, err := time.Parse(timeLayout, s.TimeOn)
		if err != nil {
			return fmt.Errorf("invalid time_on format for order %d: '%s'", s.Order, s.TimeOn)
		}
		timeOff, err := time.Parse(timeLayout, s.TimeOff)
		if err != nil {
			return fmt.Errorf("invalid time_off format for order %d: '%s'", s.Order, s.TimeOff)
		}
		if !timeOff.After(timeOn) {
			return fmt.Errorf("time_off ('%s') must be after time_on ('%s') for order %d", s.TimeOff, s.TimeOn, s.Order)
		}
		if i > 0 {
			if timeOn.Before(previousTimeOff) {
				return fmt.Errorf("collision detected: order %d starts at %s before the previous schedule ends at %s", s.Order, s.TimeOn, previousTimeOff.Format(timeLayout))
			}
		}
		previousTimeOff = timeOff
	}
	return nil
}

// --- Main Orchestration and HTTP Logic ---

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat menemukan atau memuat file .env")
	}
	cfg := config.LoadConfig()
	logger := config.SetupLogger()
	slog.SetDefault(logger)

	// BARU: Inisialisasi authenticator JWT
	tokenAuth := jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	// --- Redis Connection ---
	redisClient := cache.SetupRedisClient(cache.RedisConfig{
		RedisHost:     cfg.RedisHost,
		RedisPort:     cfg.RedisPort,
		RedisPassword: cfg.RedisPassword,
		RedisDB:       cfg.RedisDB,
	})
	defer redisClient.Close()
	slog.Info("✅ Redis connection successful")

	// --- High-Performance Setup ---
	jobChan := make(chan database.SensorReading, database.JOB_BUFFER_SIZE)

	// --- Database Connection ---
	store, err := database.NewPostgresStore(jobChan, database.DBConfig{
		DBHost:     cfg.DBHost,
		DBPort:     cfg.DBPort,
		DBUsername: cfg.DBUsername,
		DBPassword: cfg.DBPassword,
		DBDatabase: cfg.DBDatabase,
	})
	if err != nil {
		slog.Error("Failed to initialize data store", "error", err)
		os.Exit(1)
	}
	defer store.Close()
	slog.Info("✅ Database connection successful")

	// --- Sensor Mappings & Worker Pool ---
	sensorTypes, err := database.LoadSensorTypes(store.Pool)
	if err != nil {
		slog.Error("Failed to load sensor types", "error", err)
		os.Exit(1)
	}
	slog.Info("✅ Sensor types loaded into memory", "count", len(sensorTypes))
	var wg sync.WaitGroup
	for i := 0; i < database.NUM_DB_WORKERS; i++ {
		wg.Add(1)
		go database.DBWorker(i, &wg, store.Pool, store.JobChan, sensorTypes)
	}
	slog.Info("🚀 Database worker pool started", "workers", database.NUM_DB_WORKERS)

	// --- Service Setup ---
	mqttClient := broker.SetupMQTTClient(store, sensorTypes, broker.MQTTConfig{
		MQTTURL:      cfg.MQTTURL,
		MQTTUsername: cfg.MQTTUsername,
		MQTTPassword: cfg.MQTTPassword,
	})

	// DIUBAH: Teruskan tokenAuth ke SetupRouter
	router := SetupRouter(mqttClient, redisClient, tokenAuth, store)

	// --- Graceful Shutdown ---
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	srv := StartServer(router, cfg)
	<-stop

	slog.Info("⏳ Shutting down services...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "err", err)
	}
	mqttClient.Disconnect(250)
	slog.Info("MQTT client disconnected")
	close(store.JobChan)
	wg.Wait()
	slog.Info("All database workers have finished.")
	slog.Info("✅ Server gracefully stopped.")
}

func CoopAccessMiddleware(store *database.PostgresStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Dapatkan input dari URL dan JWT
			coopIDStr := chi.URLParam(r, "id")
			coopID, err := strconv.ParseInt(coopIDStr, 10, 64)
			if err != nil {
				HTTPError(w, http.StatusBadRequest, errors.New("invalid coop ID format in URL"))
				return
			}

			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				HTTPError(w, http.StatusUnauthorized, err)
				return
			}

			// Ambil 'sub' (user id) dari token
			userIDStr, ok := claims["sub"].(string)
			if !ok {
				HTTPError(w, http.StatusForbidden, errors.New("invalid 'sub' claim in token"))
				return
			}
			userID, _ := strconv.ParseInt(userIDStr, 10, 64)

			// DIHAPUS: Kita tidak lagi mengambil 'role' dari token.
			// role, ok := claims["role"].(string)

			// BARU: Ambil peran pengguna langsung dari database.
			role, err := store.GetUserRole(userID)
			if err != nil {
				// Jika pengguna atau perannya tidak ditemukan, akses ditolak.
				HTTPError(w, http.StatusForbidden, err)
				return
			}

			// Bypass otorisasi untuk admin
			if role == "admin" {
				slog.Info("Admin access granted, bypassing ownership checks.", "user_id", userID, "coop_id", coopIDStr)
				next.ServeHTTP(w, r)
				return
			}

			// 2. Query database untuk mendapatkan pemilik & ABK kandang (hanya untuk non-admin)
			farmerID, abkID, err := store.GetCoopOwnerAndABK(coopID)
			if err != nil {
				if err == pgx.ErrNoRows {
					HTTPError(w, http.StatusNotFound, errors.New("coop not found"))
					return
				}
				HTTPError(w, http.StatusInternalServerError, err)
				return
			}

			// 3. Terapkan logika otorisasi untuk non-admin menggunakan peran dari DB
			isFarmOwner := (role == "farmer" && farmerID == userID && farmerID != 0)
			isAssignedAbk := (role == "abk" && abkID == userID && abkID != 0)

			if !isFarmOwner && !isAssignedAbk {
				HTTPError(w, http.StatusForbidden, fmt.Errorf("you are not authorized to access coop '%s'", coopIDStr))
				return
			}

			// Jika lolos, lanjutkan ke handler utama
			next.ServeHTTP(w, r)
		})
	}
}

func SetupRouter(client mqtt.Client, redisClient *redis.Client, tokenAuth *jwtauth.JWTAuth, store *database.PostgresStore) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)
	r.Use(CORSMiddleware())

	// --- Rute Publik ---
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"status": "ok", "mqtt_online": client.IsConnected()})
	})

	// --- Rute Terproteksi ---
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		// DIUBAH: Panggil Authenticator sebagai fungsi dengan tokenAuth
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Use(CoopAccessMiddleware(store))

		r.Get("/coops/{id}/state", func(w http.ResponseWriter, r *http.Request) {
			coopID := chi.URLParam(r, "id")

			// Use Redis SCAN to find all hardware state keys for this coop
			pattern := fmt.Sprintf("hardware_state_%s_*", coopID)
			var hardwareStates []HardwareState
			var keys []string

			iter := redisClient.Scan(r.Context(), 0, pattern, 0).Iterator()
			for iter.Next(r.Context()) {
				keys = append(keys, iter.Val())
			}
			if err := iter.Err(); err != nil {
				slog.Error("Failed to scan Redis keys", "error", err, "pattern", pattern)
				HTTPError(w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve hardware states from Redis"))
				return
			}

			// Retrieve values for all found keys
			for _, key := range keys {
				val, err := redisClient.Get(r.Context(), key).Result()
				if err == redis.Nil {
					continue // Skip if key doesn't exist (might have expired)
				} else if err != nil {
					slog.Error("Failed to get Redis value", "error", err, "key", key)
					continue
				}

				var state HardwareState
				if err := json.Unmarshal([]byte(val), &state); err != nil {
					slog.Error("Failed to unmarshal hardware state", "error", err, "key", key)
					continue
				}
				hardwareStates = append(hardwareStates, state)
			}

			slog.Info("Retrieved hardware states from Redis", "coop_id", coopID, "count", len(hardwareStates))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"data": map[string]any{
					"coop_id":         coopID,
					"hardware_states": hardwareStates,
					"count":           len(hardwareStates),
				},
			})
		})

		r.Post("/coops/{id}/state", func(w http.ResponseWriter, r *http.Request) {
			coopID := chi.URLParam(r, "id")

			// Set body size limit
			r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit
			decoder := json.NewDecoder(r.Body)

			// First attempt to decode as partial update
			var partialUpdate PartialHardwareState
			partialErr := decoder.Decode(&partialUpdate)

			var finalState HardwareState
			var hardwareName string

			if partialErr == nil && partialUpdate.Hardware != "" {
				// This is a partial update - fetch existing state and merge
				hardwareName = partialUpdate.Hardware
				redisKey := fmt.Sprintf("hardware_state_%s_%s", coopID, hardwareName)

				// Get existing state from Redis
				existingVal, err := redisClient.Get(r.Context(), redisKey).Result()
				if err == redis.Nil {
					// No existing state, create new one with defaults
					finalState = HardwareState{
						Hardware: hardwareName,
						State:    false, // default
						Meta:     nil,
					}
				} else if err != nil {
					HTTPError(w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve existing state from Redis: %w", err))
					return
				} else {
					// Unmarshal existing state
					if err := json.Unmarshal([]byte(existingVal), &finalState); err != nil {
						slog.Error("Failed to unmarshal existing hardware state", "error", err, "key", redisKey)
						// Start with a fresh state if corrupted
						finalState = HardwareState{
							Hardware: hardwareName,
							State:    false,
							Meta:     nil,
						}
					}
				}

				// Validate schedules if provided in partial update
				if partialUpdate.Meta != nil && partialUpdate.Meta.Schedules != nil {
					if err := validateSchedules(*partialUpdate.Meta.Schedules); err != nil {
						HTTPError(w, http.StatusBadRequest, fmt.Errorf("invalid schedules for hardware '%s': %w", hardwareName, err))
						return
					}
				}

				// Merge partial updates
				finalState = mergePartialUpdates(finalState, partialUpdate)

			} else {
				// This is a full state update (backward compatibility)
				decoder = json.NewDecoder(r.Body) // Reset decoder
				var fullCommand HardwareState
				if err := decoder.Decode(&fullCommand); err != nil {
					HTTPError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON for state command: %w", err))
					return
				}
				if fullCommand.Hardware == "" {
					HTTPError(w, http.StatusBadRequest, errors.New("field 'hardware' is required"))
					return
				}
				if fullCommand.Meta != nil && len(fullCommand.Meta.Schedules) > 0 {
					if err := validateSchedules(fullCommand.Meta.Schedules); err != nil {
						HTTPError(w, http.StatusBadRequest, fmt.Errorf("invalid schedules for hardware '%s': %w", fullCommand.Hardware, err))
						return
					}
				}
				finalState = fullCommand
				hardwareName = fullCommand.Hardware
			}

			// Save the final merged state to Redis
			payload, err := json.Marshal(finalState)
			if err != nil {
				HTTPError(w, http.StatusInternalServerError, fmt.Errorf("failed to marshal state command: %w", err))
				return
			}

			redisKey := fmt.Sprintf("hardware_state_%s_%s", coopID, hardwareName)
			err = redisClient.Set(r.Context(), redisKey, payload, 0).Err()
			if err != nil {
				slog.Error("Failed to save state to Redis", "error", err, "key", redisKey)
				HTTPError(w, http.StatusInternalServerError, fmt.Errorf("failed to save state to Redis"))
				return
			}

			slog.Info("Coop hardware state saved to Redis", "key", redisKey, "hardware", hardwareName)

			// Publish to MQTT
			topic := fmt.Sprintf("coops/%s/state", coopID)
			if err := broker.MqttPublish(client, topic, payload, 1, false); err != nil {
				HTTPError(w, http.StatusInternalServerError, err)
				return
			}

			slog.Info("Coop state published to MQTT", "coop_id", coopID, "hardware", hardwareName, "topic", topic)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"message": "hardware state updated successfully",
				"data": map[string]any{
					"hardware": hardwareName,
					"coop_id":  coopID,
				},
			})
		})

		// Dummy endpoint untuk testing sensor data tanpa MQTT
		r.Post("/coops/{id}/test-sensor", func(w http.ResponseWriter, r *http.Request) {
			coopIDStr := chi.URLParam(r, "id")
			_, err := strconv.ParseInt(coopIDStr, 10, 64)
			if err != nil {
				HTTPError(w, http.StatusBadRequest, fmt.Errorf("invalid coop ID format: %w", err))
				return
			}

			// Parse JSON body untuk sensor data
			r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit
			decoder := json.NewDecoder(r.Body)

			var sensorData map[string]interface{}
			if err := decoder.Decode(&sensorData); err != nil {
				HTTPError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON payload: %w", err))
				return
			}

			// Buat dummy telemetry message
			telemetry := broker.TelemetryMessage{
				CoopID:    coopIDStr,
				Data:      sensorData,
				Timestamp: func() time.Time {
			jakartaLocation, err := time.LoadLocation("Asia/Jakarta")
			if err != nil {
				slog.Error("Failed to load Asia/Jakarta timezone, using UTC", "error", err)
				return time.Now().UTC()
			}
			return time.Now().In(jakartaLocation)
		}(),
			}

			// Load sensor types untuk validasi
			currentSensorTypes, err := database.LoadSensorTypes(store.Pool)
			if err != nil {
				slog.Error("Failed to load sensor types", "error", err)
				HTTPError(w, http.StatusInternalServerError, fmt.Errorf("failed to load sensor types: %w", err))
				return
			}

			// Simpan ke database menggunakan method yang sama dengan MQTT
			if err := store.Save(telemetry, currentSensorTypes); err != nil {
				slog.Error("Failed to save test sensor data", "error", err, "coop_id", coopIDStr)
				HTTPError(w, http.StatusInternalServerError, fmt.Errorf("failed to save sensor data: %w", err))
				return
			}

			slog.Info("Test sensor data saved successfully", "coop_id", coopIDStr, "data_count", len(sensorData))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{
				"ok": true,
				"message": "test sensor data saved successfully",
				"data": map[string]any{
					"coop_id":      coopIDStr,
					"sensor_count": len(sensorData),
					"timestamp":    telemetry.Timestamp.Format(time.RFC3339),
				},
			})
		})
	})

	return r
}

// ... (Sisa fungsi: StartServer, HTTPError, CORSMiddleware tidak berubah) ...
func StartServer(router http.Handler, cfg config.Config) *http.Server {
	port := cfg.Port
	srv := &http.Server{Addr: ":" + port, Handler: router}
	go func() {
		slog.Info("✅ HTTP server listening", "url", fmt.Sprintf("http://localhost:%s", port))
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()
	return srv
}
func HTTPError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]any{"ok": false, "error": err.Error()})
}
func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
