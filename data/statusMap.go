package data

import (
	"average-watcher-bot/checker"
	"encoding/json"
	"log"
	"os"
)

// Загружает данные из файла storage/statusMap.json
func LoadStatusMap() map[string]bool {
	var data, err = os.ReadFile("storage/statusMap.json")
	var statusMap = map[string]bool{}

	if err == nil {
		JsonErr := json.Unmarshal(data, &statusMap)
		if JsonErr != nil {
			log.Printf(JsonErr.Error())
		}
	}
	return statusMap
}

// Проверяет статусы всех хостов и создаёт сохраняет их в файл
func GenerateStatusMap(watchlist []string) map[string]bool {
	statusMap := make(map[string]bool)

	// Проверяем все адреса из watchlist
	for _, ip := range watchlist {
		statusMap[ip] = checker.CheckICMP(ip)
	}

	SaveStatusMap(statusMap)
	return statusMap
}

// Сохраняет данные в фалй storage/statusMap.json
func SaveStatusMap(statusMap map[string]bool) {
	jsonStatusMap, err := json.Marshal(statusMap)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("storage/statusMap.json", jsonStatusMap, 0666)
}
