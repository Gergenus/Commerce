package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresURL string
	LogLevel    string
	HTTPPort    string
	JWTSecret   string
	AccessTTL   time.Duration
	RefreshTTl  time.Duration
}

func InitConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	AccessTTL, err := time.ParseDuration(os.Getenv("AccessTTL"))
	if err != nil {
		panic(err)
	}
	RefreshTTl, err := time.ParseDuration(os.Getenv("RefreshTTl"))
	if err != nil {
		panic(err)
	}
	return Config{
		PostgresURL: os.Getenv("PostgresURL"),
		LogLevel:    os.Getenv("LogLevel"),
		HTTPPort:    os.Getenv("HTTPPort"),
		JWTSecret:   os.Getenv("JWTSecret"),
		AccessTTL:   AccessTTL,
		RefreshTTl:  RefreshTTl,
	}
}
