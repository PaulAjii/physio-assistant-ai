package config

import (
	"os"
	"strings"
)

// Config holds all runtime configuration for backend-core. It is loaded once at
// startup from environment variables. New settings (DB DSN, JWT secret, S3
// credentials, ...) should be added here as later phases land.
type Config struct {
	Port           string
	AIBackendURI   string
	UploadDir      string
	AllowedOrigins []string
}

// Load reads configuration from the environment, applying sensible defaults so
// the service still runs with no env vars set (useful for local development).
func Load() Config {
	origins := strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "*"), ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return Config{
		Port:           getEnv("PORT", "8080"),
		AIBackendURI:   getEnv("AI_BACKEND_URI", "http://localhost:5000/ai/process-audio"),
		UploadDir:      getEnv("UPLOAD_DIR", "uploads"),
		AllowedOrigins: origins,
	}
}

// getEnv returns the value of the environment variable named by key, or
// fallback when the variable is unset or empty.
func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
