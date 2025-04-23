package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Mode string
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


