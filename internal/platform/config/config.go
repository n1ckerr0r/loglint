package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Addr        string
	DatabaseURL string
	RedisURL    string
	RabbitMQURL string
	QueueName   string
	Shutdown    time.Duration
}

func FromEnv() Config {
	return Config{
		Addr:        env("LOGLINT_ADDR", ":8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
		RabbitMQURL: os.Getenv("RABBITMQ_URL"),
		QueueName:   env("LOGLINT_QUEUE", "loglint.runs"),
		Shutdown:    durationEnv("LOGLINT_SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}
