package envutils

import (
	"log"
	"log/slog"
	"os"
	"path"

	"github.com/joho/godotenv"
)

func LoadEnvFiles(filenames ...string) {
	if len(filenames) == 0 {
		rootEnv := path.Join(os.Getenv("PWD"), "../../../.env")
		filenames = []string{rootEnv}
	}
	if err := godotenv.Load(filenames...); err != nil {
		slog.Warn("No .env file found")
	}
}

func RequireEnv(key string) string {
	value := os.Getenv(key)
	slog.Info("RequireEnv", "key", key, "value", value)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}
