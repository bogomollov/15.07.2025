package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string
}

func Load() *Config {
	
	err := godotenv.Load()
  	if err != nil {
    	log.Fatal("Ошибка загрузки .env файла")
  	}
	
	config := Config{PORT: os.Getenv("PORT")}

	return &config
}