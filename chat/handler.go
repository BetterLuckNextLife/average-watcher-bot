package chat

import (
	"average-watcher-bot/data"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ListenUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		text := update.Message.Text

		if strings.HasPrefix(text, "/add ") {
			ip := strings.TrimSpace(strings.TrimPrefix(text, "/add "))
			ok, err := data.AddToWatchList(ip)
			if ok == false && err == nil {
				send(update.Message.Chat.ID, "⚠️ Некорректный IP!")
			}
			if err != nil {
				send(update.Message.Chat.ID, "⚠️ Ошибка при добавлении")
			} else if ok == false {
				send(update.Message.Chat.ID, "⚠️ IP уже отслеживается")
			} else {
				send(update.Message.Chat.ID, "✔️ IP успешно добавлен: "+ip)
			}
		}

		if text == "/list" {
			list := data.LoadWatchList()
			msg := "🎯 Отслеживаемые IP:\n"
			for _, ip := range list {
				msg += "• " + ip + "\n"
			}
			send(update.Message.Chat.ID, msg)
		}

		if strings.HasPrefix(text, "/remove ") {
			ip := strings.TrimSpace(strings.TrimPrefix(text, "/remove "))
			ok, err := data.RemoveFromWatchList(ip)
			if err != nil {
				send(update.Message.Chat.ID, "⚠️ Ошибка при удалении")
			} else if !ok {
				send(update.Message.Chat.ID, "⚠️ Такого IP нет в списке")
			} else {
				send(update.Message.Chat.ID, "🗑️ IP удалён: "+ip)
			}
		}

		if strings.HasPrefix(text, "/status") {
			statusMap := data.LoadStatusMap()
			msg := "🖥 Статус серверов сейчас: \n"
			for ip, status := range statusMap {
				var statusText string
				if status {
					statusText = "🟢 Онлайн"
				} else {
					statusText = "🔴 Оффлайн"
				}
				msg += "• " + ip + ": " + statusText + "\n"
			}
			send(update.Message.Chat.ID, msg)
		}

	}
}

func send(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	bot.Send(msg)
}
