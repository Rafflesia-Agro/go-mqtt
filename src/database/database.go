package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	// PERBAIKAN: Impor paket broker untuk menggunakan tipe TelemetryMessage
	"github.com/Anjasfedo/go-mqtt/src/broker"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Configuration constants for our batching workers (these are logical constants)
const (
	JOB_BUFFER_SIZE = 20000
	MAX_BATCH_SIZE  = 1000
	BATCH_TIMEOUT   = 1 * time.Second
	NUM_DB_WORKERS  = 10
)

// SensorReading represents a single, normalized row to be inserted into the database.
type SensorReading struct {
	CoopID       int64
	SensorTypeID int32
	Value        float64
	Timestamp    time.Time
}

// GetUTCTime returns current time in UTC timezone
func GetUTCTime() time.Time {
	return time.Now().UTC()
}

// DBConfig defines the required fields for DB connection.
type DBConfig struct {
	DBHost, DBPort, DBUsername, DBPassword, DBDatabase string
}

// PostgresStore holds the DB pool, Redis client, and the job channel.
type PostgresStore struct {
	Pool       *pgxpool.Pool
	Redis      *redis.Client
	JobChan    chan SensorReading
}

// NewPostgresStore initializes the database connection pool and the store.
func NewPostgresStore(jobChan chan SensorReading, cfg DBConfig, redisClient *redis.Client) (*PostgresStore, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBDatabase,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set session timezone to Asia/Jakarta
	if _, err := pool.Exec(context.Background(), "SET TIME ZONE 'Asia/Jakarta'"); err != nil {
		return nil, fmt.Errorf("failed to set timezone: %w", err)
	}

	return &PostgresStore{Pool: pool, Redis: redisClient, JobChan: jobChan}, nil
}

// PERBAIKAN: struct TelemetryMessage dihapus dari sini

// Save converts the message and pushes jobs onto the buffered channel.
// PERBAIKAN: Method ini sekarang menerima 'broker.TelemetryMessage' agar sesuai dengan interface.
func (s *PostgresStore) Save(msg broker.TelemetryMessage, sensorTypes map[string]int32) error {
	coopID, err := strconv.ParseInt(msg.CoopID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid coop_id format: %s", msg.CoopID)
	}

	// Store the last sensor data timestamp in Redis
	key := fmt.Sprintf("last_sensor_stored_%s", msg.CoopID)
	currentTime := msg.Timestamp.Format(time.RFC3339)
	if err := s.Redis.Set(context.Background(), key, currentTime, 0).Err(); err != nil {
		slog.Error("Failed to store last sensor timestamp", "error", err, "coop_id", msg.CoopID)
		// Continue anyway since this is not critical for database storage
	}

	for key, val := range msg.Data {
		sensorName := strings.ToLower(key)
		sensorTypeID, ok := sensorTypes[sensorName]
		if !ok {
			slog.Warn("Unknown sensor type received, skipping", "type", sensorName, "coop_id", msg.CoopID)
			continue
		}

		value, ok := val.(float64)
		if !ok {
			slog.Warn("Invalid value type for sensor, skipping", "type", sensorName, "value", val)
			continue
		}

		reading := SensorReading{
			CoopID:       coopID,
			SensorTypeID: sensorTypeID,
			Value:        value,
			Timestamp:    msg.Timestamp,
		}

		s.JobChan <- reading
	}
	return nil
}

func (s *PostgresStore) GetCoopOwnerAndABK(coopID int64) (int64, int64, error) {
	query := `
		SELECT f.farmer_id, c.abk_id
		FROM coops c
		LEFT JOIN farms f ON c.farm_id = f.id
		WHERE c.id = $1
	`
	var farmerID sql.NullInt64
	var abkID sql.NullInt64

	err := s.Pool.QueryRow(context.Background(), query, coopID).Scan(&farmerID, &abkID)

	if err != nil {
		if err == pgx.ErrNoRows {
			// Mengembalikan error spesifik jika kandang tidak ditemukan
			return 0, 0, pgx.ErrNoRows
		}
		return 0, 0, err
	}

	// Mengonversi sql.NullInt64 ke int64, mengembalikan 0 jika NULL
	return farmerID.Int64, abkID.Int64, nil
}

func (s *PostgresStore) GetUserRole(userID int64) (string, error) {
	query := `
		SELECT r.name
		FROM users u
		JOIN roles r ON u.current_role_id = r.id
		WHERE u.id = $1
	`
	var roleName string
	err := s.Pool.QueryRow(context.Background(), query, userID).Scan(&roleName)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Mengembalikan error yang jelas jika pengguna atau perannya tidak ditemukan
			return "", fmt.Errorf("user with id %d or their role not found", userID)
		}
		return "", err
	}
	return roleName, nil
}

// Close gracefully shuts down the database connection pool.
func (s *PostgresStore) Close() {
	s.Pool.Close()
}

// DBWorker runs in its own goroutine, processing jobs from the channel.
func DBWorker(id int, wg *sync.WaitGroup, pool *pgxpool.Pool, jobs <-chan SensorReading, sensorTypes map[string]int32) {
	defer wg.Done()
	slog.Info("DB worker started", "id", id)

	batch := make([]SensorReading, 0, MAX_BATCH_SIZE)
	ticker := time.NewTicker(BATCH_TIMEOUT)
	defer ticker.Stop()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				if len(batch) > 0 {
					slog.Info("Channel closed, processing final batch", "worker_id", id, "batch_size", len(batch))
					if err := batchInsert(pool, batch); err != nil {
						slog.Error("Final batch insert failed", "err", err, "worker_id", id)
					}
				}
				slog.Info("DB worker shutting down", "id", id)
				return
			}
			batch = append(batch, job)
			if len(batch) >= MAX_BATCH_SIZE {
				if err := batchInsert(pool, batch); err != nil {
					slog.Error("Batch insert failed", "err", err, "worker_id", id)
				}
				batch = make([]SensorReading, 0, MAX_BATCH_SIZE)
				ticker.Reset(BATCH_TIMEOUT)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				if err := batchInsert(pool, batch); err != nil {
					slog.Error("Batch insert failed on timeout", "err", err, "worker_id", id)
				}
				batch = make([]SensorReading, 0, MAX_BATCH_SIZE)
			}
		}
	}
}

// batchInsert performs a bulk insert using PostgreSQL's highly efficient COPY protocol.
func batchInsert(pool *pgxpool.Pool, readings []SensorReading) error {
	if len(readings) == 0 {
		return nil
	}

	rows := make([][]interface{}, len(readings))
	for i, r := range readings {
		rows[i] = []interface{}{r.Value, r.SensorTypeID, r.CoopID, r.Timestamp, r.Timestamp}
	}

	_, err := pool.CopyFrom(
		context.Background(),
		pgx.Identifier{"rec_sensor_coops"},
		[]string{"value", "sensor_type_id", "coop_id", "created_at", "updated_at"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		slog.Error("COPY From failed", "error", err)
		return err
	}

	slog.Info("Successfully inserted batch", "rows", len(readings), "timezone", "Asia/Jakarta")
	return nil
}

// LoadSensorTypes queries the DB on startup to map sensor names to their IDs.
func LoadSensorTypes(pool *pgxpool.Pool) (map[string]int32, error) {
	rows, err := pool.Query(context.Background(), "SELECT id, name FROM sensor_types")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	types := make(map[string]int32)
	for rows.Next() {
		var id int32
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		types[strings.ToLower(name)] = id
	}
	return types, nil
}
