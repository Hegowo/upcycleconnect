package config

import "os"

type Config struct {
	AppKey        string
	DBHost        string
	DBPort        string
	DBDatabase    string
	DBUsername    string
	DBPassword    string
	LogLevel      string
	SocieteApiKey string
}

func Load() *Config {
	return &Config{
		AppKey:        getEnv("APP_KEY", ""),
		DBHost:        getEnv("DB_HOST", "mysql"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBDatabase:    getEnv("DB_DATABASE", "upcycleconnect"),
		DBUsername:    getEnv("DB_USERNAME", "ucadmin"),
		DBPassword:    getEnv("DB_PASSWORD", "ucpassword"),
		LogLevel:      getEnv("LOG_LEVEL", "error"),
		SocieteApiKey: getEnv("SOCIETE_API_KEY", "50567e30ca8931ff0defa56a5c5d7b61"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
