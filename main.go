package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-ping/ping"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func loadDotenv() (string, int64) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Файл окружения не найден")
	}

	token := os.Getenv("TG_TOKEN")
	chatID, _ := strconv.ParseInt(os.Getenv("TG_CHAT_ID"), 10, 64)

	return token, chatID
}

func loadStatusMap() map[string]bool {
	var data, err = os.ReadFile("statusMap.json")
	var statusMap = map[string]bool{}

	if err == nil {
		JsonErr := json.Unmarshal(data, &statusMap)
		if JsonErr != nil {
			fmt.Printf(JsonErr.Error())
		}
	}
	return statusMap
}

func checkICMP(ip string) bool {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		log.Println("Ошибка создания пингера:", err)
		return false
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 2

	err = pinger.Run()
	if err != nil {
		log.Println("Ошибка выполнения пинга:", err)
		return false
	}

	stats := pinger.Statistics()
	return stats.PacketsRecv > 0
}

func saveStatusMap(statusMap map[string]bool) {
	jsonStatusMap, err := json.Marshal(statusMap)
	if err != nil {
		log.Fatal(err)
	}
	os.WriteFile("statusMap.json", jsonStatusMap, 0666)
}

func generateStatusMap(watchlist []string) map[string]bool {
	statusMap := make(map[string]bool)

	// Проверяем все адреса из watchlist
	for _, ip := range watchlist {
		statusMap[ip] = checkICMP(ip)
	}

	saveStatusMap(statusMap)
	return statusMap
}

func main() {
	// Загружаем .env
	token, chatID := loadDotenv()

	// Стартуем бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	var watchlist = []string{"192.168.1.1", "192.168.100.1", "192.168.100.2", "192.168.100.3"}

	// Загрузить/создать словарь статусов
	statusMap := loadStatusMap()
	if len(statusMap) == 0 {
		statusMap = generateStatusMap(watchlist)
	}

	// Отправляем данные
	document := tgbotapi.NewDocument(chatID, tgbotapi.FilePath("statusMap.json"))
	bot.Send(document)
	bot.Debug = true
}
