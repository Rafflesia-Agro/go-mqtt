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

// HAPUS STRUCT INI: StatePayload tidak lagi digunakan.
// type StatePayload struct {
// 	HardwareStates []HardwareState `json:"hardware_state"`
// }

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

		r.Post("/coops/{id}/state", func(w http.ResponseWriter, r *http.Request) {
			coopID := chi.URLParam(r, "id")

			var command HardwareState
			if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
				HTTPError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON for state command: %w", err))
				return
			}
			if command.Hardware == "" {
				HTTPError(w, http.StatusBadRequest, errors.New("field 'hardware' is required"))
				return
			}
			if command.Meta != nil && len(command.Meta.Schedules) > 0 {
				if err := validateSchedules(command.Meta.Schedules); err != nil {
					HTTPError(w, http.StatusBadRequest, fmt.Errorf("invalid schedules for hardware '%s': %w", command.Hardware, err))
					return
				}
			}

			payload, err := json.Marshal(command)
			if err != nil {
				HTTPError(w, http.StatusInternalServerError, fmt.Errorf("failed to marshal state command: %w", err))
				return
			}

			redisKey := fmt.Sprintf("hardware_state_%s_%s", coopID, command.Hardware)
			err = redisClient.Set(r.Context(), redisKey, payload, 0).Err()
			if err != nil {
				slog.Error("Failed to save state to Redis", "error", err, "key", redisKey)
			} else {
				slog.Info("Coop hardware state saved to Redis", "key", redisKey)
			}

			topic := fmt.Sprintf("coops/%s/state", coopID)
			if err := broker.MqttPublish(client, topic, payload, 1, false); err != nil {
				HTTPError(w, http.StatusInternalServerError, err)
				return
			}

			slog.Info("Coop state published to MQTT", "coop_id", coopID, "hardware", command.Hardware, "topic", topic)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(map[string]any{"ok": true, "message": "state command for specific hardware sent"})
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
