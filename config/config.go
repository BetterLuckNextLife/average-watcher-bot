package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Загружает данные из окружения
func LoadToken() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Файл окружения не найден")
	}

	token := os.Getenv("TG_TOKEN")

	return token
}
