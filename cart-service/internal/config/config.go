package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	LogLevel      string
	JWTSecret     string
	RedisAddress  string
	GRPCAddress   string
	RedisDB       int
	RedisPassword string
	CartTTL       time.Duration
	HTTPPort      string
}

func InitConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	redisdb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}
	cartTTL, err := time.ParseDuration(os.Getenv("CART_TTL"))
	if err != nil {
		panic(err)
	}
	return Config{
		LogLevel:      os.Getenv("LOG_LEVEL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		RedisAddress:  os.Getenv("REDIS_ADDRESS"),
		GRPCAddress:   os.Getenv("GRPC_ADDRESS"),
		RedisDB:       redisdb,
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		CartTTL:       cartTTL,
		HTTPPort:      os.Getenv("HTTP_PORT"),
	}
}
