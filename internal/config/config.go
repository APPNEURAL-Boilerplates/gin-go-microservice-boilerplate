package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServiceName            string
	Environment            string
	Host                   string
	Port                   string
	GinMode                string
	LogLevelName           string
	ReadTimeoutSeconds     int
	WriteTimeoutSeconds    int
	ShutdownTimeoutSeconds int
	TrustedProxies         []string
}

func Load() Config {
	environment := getEnv("ENVIRONMENT", "development")

	return Config{
		ServiceName:            getEnv("SERVICE_NAME", "gin-microservice"),
		Environment:            environment,
		Host:                   getEnv("HOST", "0.0.0.0"),
		Port:                   getEnv("PORT", "8080"),
		GinMode:                getEnv("GIN_MODE", defaultGinMode(environment)),
		LogLevelName:           getEnv("LOG_LEVEL", "info"),
		ReadTimeoutSeconds:     getEnvAsInt("READ_TIMEOUT_SECONDS", 10),
		WriteTimeoutSeconds:    getEnvAsInt("WRITE_TIMEOUT_SECONDS", 10),
		ShutdownTimeoutSeconds: getEnvAsInt("SHUTDOWN_TIMEOUT_SECONDS", 10),
		TrustedProxies:         splitCSV(getEnv("TRUSTED_PROXIES", "")),
	}
}

func (c Config) Addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func (c Config) SlogLevel() slog.Level {
	switch strings.ToLower(strings.TrimSpace(c.LogLevelName)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists || strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func getEnvAsInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func defaultGinMode(environment string) string {
	if strings.EqualFold(environment, "development") || strings.EqualFold(environment, "local") {
		return "debug"
	}
	return "release"
}

func splitCSV(value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}
