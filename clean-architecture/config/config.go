package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port     string
	GinMode  string
	Timeout  time.Duration
	Shutdown time.Duration
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type LoggerConfig struct {
	Level  string
	Format string
}

func Load() *Config {
	return &Config{
		Server:   loadServerConfig(),
		Database: loadDatabaseConfig(),
		Logger:   loadLoggerConfig(),
	}
}

func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port:     getEnv("PORT", "8080"),
		GinMode:  getEnv("GIN_MODE", "release"),
		Timeout:  time.Duration(getEnvAsInt("SERVER_TIMEOUT", 30)) * time.Second,
		Shutdown: time.Duration(getEnvAsInt("SERVER_SHUTDOWN_TIMEOUT", 10)) * time.Second,
	}
}

func loadDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnvAsInt("DB_PORT", 5432),
		User:            getEnv("DB_USER", "api_user"),
		Password:        getEnv("DB_PASSWORD", "api_password"),
		Database:        getEnv("DB_NAME", "api_db"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME", 3600)) * time.Second,
	}
}

func loadLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:  getEnv("LOG_LEVEL", "info"),
		Format: getEnv("LOG_FORMAT", "json"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
