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
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPassword  string
	MailFrom      string
	AppURL        string

	StripeSecret        string
	StripePublicKey     string
	StripeWebhookSecret string
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
		SMTPHost:      getEnv("SMTP_HOST", "ssl0.ovh.net"),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		SMTPUser:      getEnv("SMTP_USER", "noreply@upcycleconnect.xyz"),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
		MailFrom:      getEnv("MAIL_FROM", "noreply@upcycleconnect.xyz"),
		AppURL:        getEnv("APP_URL", "https://upcycleconnect.xyz"),

		StripeSecret:        getEnv("STRIPE_SECRET", ""),
		StripePublicKey:     getEnv("STRIPE_PUBLIC_KEY", ""),
		StripeWebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
