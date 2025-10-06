package cache

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// DIUBAH: Definisikan dan ekspor struct Config di sini, bukan di dalam fungsi.
// Ini menjadi "sumber kebenaran" untuk konfigurasi Redis.
type RedisConfig struct {
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
}

// DIUBAH: Ubah parameter dari 'interface{}' menjadi 'RedisConfig' yang spesifik.
func SetupRedisClient(cfg RedisConfig) *redis.Client {

	// DIHAPUS: Definisi 'type RedisConfig' dan konversi 'cfg.(RedisConfig)' tidak lagi diperlukan.

	db, err := strconv.Atoi(cfg.RedisDB)
	if err != nil {
		slog.Error("Invalid REDIS_DB environment variable", "value", cfg.RedisDB, "error", "must be an integer")
		os.Exit(1)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}

	return rdb
}
