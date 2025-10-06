package config

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

// Config holds all necessary environment variables for the application.
type Config struct {
	// HTTP
	Port string
	// DB
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	DBDatabase string
	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       string
	// MQTT
	MQTTURL      string
	MQTTUsername string
	MQTTPassword string
	// BARU: Tambahkan JWT Secret
	JWTSecret string
}

// LoadConfig loads all required environment variables strictly.
// The application will exit if any mandatory variable is missing.
func LoadConfig() Config {
	return Config{
		// HTTP
		Port: GetRequiredEnv("PORT"),
		// DB
		DBHost:     GetRequiredEnv("DB_HOST"),
		DBPort:     GetRequiredEnv("DB_PORT"),
		DBUsername: GetRequiredEnv("DB_USERNAME"),
		DBPassword: GetRequiredEnv("DB_PASSWORD"),
		DBDatabase: GetRequiredEnv("DB_DATABASE"),
		// Redis
		RedisHost:     GetRequiredEnv("REDIS_HOST"),
		RedisPort:     GetRequiredEnv("REDIS_PORT"),
		RedisPassword: GetRequiredEnv("REDIS_PASSWORD"),
		RedisDB:       GetRequiredEnv("REDIS_DB"),
		// MQTT
		MQTTURL:      GetRequiredEnv("MQTT_URL"),
		MQTTUsername: GetRequiredEnv("MQTT_USERNAME"),
		MQTTPassword: GetRequiredEnv("MQTT_PASSWORD"),
		// BARU: Muat JWT Secret dari environment
		JWTSecret: GetRequiredEnv("JWT_SECRET"),
	}
}

// SetupLogger allows LOG_LEVEL to be optional, defaulting to "info".
func SetupLogger() *slog.Logger {
	level := slog.LevelInfo
	if strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") {
		level = slog.LevelDebug
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}

// GetRequiredEnv retrieves the value for the given environment key.
// If the variable is not set or is empty, the application exits with a fatal error.
func GetRequiredEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Fatalf("FATAL ERROR: Required environment variable %s is not set.", key)
	return "" // Should not be reached
}
