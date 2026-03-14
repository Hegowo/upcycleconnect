package config

import "os"

type Config struct {
	AppKey     string
	DBHost     string
	DBPort     string
	DBDatabase string
	DBUsername string
	DBPassword string
	LogLevel   string
}

func Load() *Config {
	return &Config{
		AppKey:     getEnv("APP_KEY", ""),
		DBHost:     getEnv("DB_HOST", "mysql"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBDatabase: getEnv("DB_DATABASE", "upcycleconnect"),
		DBUsername:  getEnv("DB_USERNAME", "ucadmin"),
		DBPassword: getEnv("DB_PASSWORD", "ucpassword"),
		LogLevel:   getEnv("LOG_LEVEL", "error"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
