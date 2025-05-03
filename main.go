package main

import (
	"average-watcher-bot/chat"
	"average-watcher-bot/config"
	"average-watcher-bot/data"
	"time"
)

func main() {
	// Загружаем .env
	token, watchers := config.LoadDotenv()
	interval := 1

	// Загружаем список наблюдаемых
	var watchlist = []string{"192.168.101.11", "192.168.101.12", "192.168.101.13"}

	// Загрузить/создать словарь статусов
	statusMap := data.LoadStatusMap()
	if len(statusMap) == 0 {
		statusMap = data.GenerateStatusMap(watchlist)
	}

	// Стартуем бота
	chat.StartUp(token)

	// Сигнализируем о включении
	chat.StartupNotify(watchers, watchlist)

	// Проверяем статус хостов и уведомляем при изменениях
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			data.UpdateStatusMapAndAlert(statusMap, watchers)
		}
	}

}
