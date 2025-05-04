package chat

import (
	"average-watcher-bot/checker"
	"average-watcher-bot/data"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot *tgbotapi.BotAPI

// Запуск бота
func StartUp(token string) {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

}

// Уведомление о изменении статуса хоста
func SendAlert(chatID int64, ip string, status bool) {
	var msg tgbotapi.MessageConfig
	if status == true {
		msg = tgbotapi.NewMessage(chatID, "✅ Сервер по адресу "+ip+" теперь онлайн!")
	} else {
		msg = tgbotapi.NewMessage(chatID, "❌ Сервер по адресу "+ip+" не отвечает!")
	}

	bot.Send(msg)
}

// Уведомление о старте бота
func StartupNotify(watchers []int64, watchlist []string) {
	for _, watcherID := range watchers {
		msg := tgbotapi.NewMessage(watcherID,
			"*🟢 Бот запущен!*\n"+
				"🎯 Отслежываемые ip: "+strconv.Itoa(len(watchlist))+"\n"+
				"🕶 Отслеживающих: "+strconv.Itoa(len(watchers))+"\n")
		bot.Send(msg)
	}
}

// Проверяет статус хостов, при изменении обновляет его, уведомляет наблюдателей
func UpdateStatusMapAndAlert(statusMap map[string]bool, watchers []int64) map[string]bool {
	for ip, status := range statusMap {
		newStatus := checker.CheckICMP(ip)
		if status != newStatus {
			for _, watcherID := range watchers {
				statusMap[ip] = newStatus
				SendAlert(watcherID, ip, newStatus)
			}
		}
	}

	data.SaveStatusMap(statusMap)
	return statusMap
}
