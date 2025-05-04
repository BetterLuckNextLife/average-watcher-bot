package data

import (
	"encoding/json"
	"log"
	"os"
)

// Загружает данные из файла storage/watchers.json
func LoadWatchers() []int64 {
	watchers := []int64{}

	data, err := os.ReadFile("storage/watchers.json")
	if err != nil {
		log.Printf("Ошибка чтения watchers.json: %v", err)
		return watchers
	}

	JsonErr := json.Unmarshal(data, &watchers)
	if JsonErr != nil {
		log.Printf("Ошибка парсинга watchers.json: %v", JsonErr)
	}

	if len(watchers) == 0 {
		log.Print("Не загружено ни одного наблюдателя!")
	} else {
		log.Printf("Загружено наблюдателей: %v", watchers)
	}

	return watchers
}
