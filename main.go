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

func loadDotenv() (string, []int64) {
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
	pinger.Count = 5
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

func sendAlert(chatID int64, ip string, status bool) {
	msg := tgbotapi.NewMessage(chatID, "Сервер по адресу "+ip+" сменил статус на "+strconv.FormatBool(status))
	bot.Send(msg)
}

func updateStatusMapAndAlert(statusMap map[string]bool, watchers []int64) map[string]bool {
	// Проверяем все адреса из watchlist
	for ip, status := range statusMap {
		newStatus := checkICMP(ip)
		if status != newStatus {
			for _, watcherID := range watchers {
				sendAlert(watcherID, ip, status)
			}
			statusMap[ip] = newStatus
		}
	}

	saveStatusMap(statusMap)
	return statusMap
}

var bot *tgbotapi.BotAPI

func main() {
	// Загружаем .env
	token, watchers := loadDotenv()

	// Стартуем бота
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	var watchlist = []string{"192.168.1.1", "192.168.100.1", "192.168.100.2", "192.168.100.3"}

	// Загрузить/создать словарь статусов
	statusMap := loadStatusMap()
	if len(statusMap) == 0 {
		statusMap = generateStatusMap(watchlist)
	}

	// Сигнализируем о включении
	for _, watcherID := range watchers {
		msg := tgbotapi.NewMessage(watcherID,
			"Бот запущен!\n"+
				"Отслежываемые ip: "+strconv.Itoa(len(watchlist))+"\n"+
				"Отслеживающих: "+strconv.Itoa(len(watchers))+"\n")
		bot.Send(msg)
	}

	// Отправляем данные
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			updateStatusMapAndAlert(statusMap, watchers)
		}
	}

}
