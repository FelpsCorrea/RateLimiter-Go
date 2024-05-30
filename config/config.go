package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	RateLimitPerIP    int
	RateLimitPerToken int
	BlockDuration     time.Duration
	RedisAddr         string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	blockDuration, err := time.ParseDuration(os.Getenv("BLOCK_DURATION"))
	if err != nil {
		log.Fatal("Invalid block duration format")
	}

	return &Config{
		RateLimitPerIP:    getIntFromEnv("RATE_LIMIT_PER_IP", 10),
		RateLimitPerToken: getIntFromEnv("RATE_LIMIT_PER_TOKEN", 100),
		BlockDuration:     blockDuration,
		RedisAddr:         os.Getenv("REDIS_ADDR"),
	}
}

func getIntFromEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Invalid value for %s, using default %d", key, defaultValue)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}
