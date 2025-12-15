package config

import (
	"log"
	"os"
)

type Config struct {
	DBDSN string
	Port  string
}

func Load() Config {
	cfg := Config{
		DBDSN: envOrDefault("DB_DSN", "postgres://darts_user:darts_pass@localhost:5432/darts?sslmode=disable"),
		Port:  envOrDefault("APP_PORT", "8080"),
	}

	if cfg.DBDSN == "" {
		log.Fatal("DB_DSN must be set")
	}

	return cfg
}

func envOrDefault(key, def string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return def
}
