package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Redis struct {
	RedisHost string
	RedisPort string
	RedisPass string
	RedisAuthDB string
}

type Config struct {
	Mode string
	Redis Redis
}

var (
	config_instance *Config
	once sync.Once
)

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	config_instance = &Config{
		Mode: getEnv("MODE", "development"),
		Redis: Redis{
			RedisHost:   getEnv("REDIS_HOST", "192.168.0.103"),
			RedisPort:   getEnv("REDIS_PORT", "6379"),
			RedisPass:   getEnv("REDIS_PASS", ""),
			RedisAuthDB: getEnv("REDIS_AUTH_DB", "0"),
		},
	}

	return config_instance, nil
}

func GetConfig() *Config {
	once.Do(func() {
		if _, err := LoadConfig(); err != nil {
			log.Fatalf("Error loading config: %v\n", err)
		}
	})
	return config_instance
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}


