package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	APP_URL  string
	APP_PORT string
}

var (
	instance *Config
	once     sync.Once
)

func Load() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Ошибка загрузки .env файла")
		}

		instance = &Config{APP_URL: os.Getenv("APP_URL"), APP_PORT: os.Getenv("APP_PORT")}
	})
	return instance
}

func GetURL() string {
	return instance.APP_URL
}

func GetPort() string {
	return instance.APP_PORT
}
