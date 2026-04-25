package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret string

	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string

	AppEnv  string
	AppPort string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "golang_api"),

		JWTSecret: getEnv("JWT_SECRET", ""),

		SMTPHost: getEnv("SMTP_HOST", "localhost"),
		SMTPPort: getEnv("SMTP_PORT", "1025"),
		SMTPUser: getEnv("SMTP_USER", ""),
		SMTPPass: getEnv("SMTP_PASS", ""),

		AppEnv:  getEnv("APP_ENV", "development"),
		AppPort: getEnv("APP_PORT", "8080"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET é obrigatório")
	}

	return cfg, nil
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}