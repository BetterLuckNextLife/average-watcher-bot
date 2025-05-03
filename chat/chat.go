package chat

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	msg := tgbotapi.NewMessage(chatID, "Сервер по адресу "+ip+" сменил статус на "+strconv.FormatBool(status))
	bot.Send(msg)
}

// Уведомление о старте бота
func StartupNotify(watchers []int64, watchlist []string) {
	for _, watcherID := range watchers {
		msg := tgbotapi.NewMessage(watcherID,
			"Бот запущен!\n"+
				"Отслежываемые ip: "+strconv.Itoa(len(watchlist))+"\n"+
				"Отслеживающих: "+strconv.Itoa(len(watchers))+"\n")
		bot.Send(msg)
	}
}
