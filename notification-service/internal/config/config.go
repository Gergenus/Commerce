package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel          string
	JWTSecret         string
	TokenTTL          time.Duration
	FromEmail         string
	FromEmailPassword string
	FromEmailSMTP     string
	SMTPAddr          string
}

func InitConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	TokenTTL, err := time.ParseDuration(os.Getenv("TokenTTL"))
	if err != nil {
		panic(err)
	}

	return Config{
		LogLevel:          os.Getenv("LogLevel"),
		JWTSecret:         os.Getenv("JWTSecret"),
		TokenTTL:          TokenTTL,
		FromEmail:         os.Getenv("FROM_EMAIL"),
		FromEmailPassword: os.Getenv("FROM_EMAIL_PASSWORD"),
		FromEmailSMTP:     os.Getenv("FROM_EMAIL_SMTP"),
		SMTPAddr:          os.Getenv("SMTP_ADDR"),
	}
}
