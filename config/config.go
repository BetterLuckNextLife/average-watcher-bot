package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Загружает данные из окружения
func LoadDotenv() (string, []int64) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Файл окружения не найден")
	}

	token := os.Getenv("TG_TOKEN")

	// TODO: Добавить поддержку множества ID
	chatID, _ := strconv.ParseInt(os.Getenv("TG_CHAT_ID"), 10, 64)
	var watchers []int64
	watchers = append(watchers, chatID)

	return token, watchers
}
