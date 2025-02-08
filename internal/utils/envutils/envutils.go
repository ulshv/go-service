package envutils

import (
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvFiles(filenames ...string) {
	if err := godotenv.Load(filenames...); err != nil {
		slog.Warn("No .env file found")
	}
}

func RequireEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
