package main

import (
	"average-watcher-bot/chat"
	"average-watcher-bot/config"
	"average-watcher-bot/data"
	"time"
)

func main() {
	// Загружаем данные
	interval := 1

	token := config.LoadToken()
	watchlist := data.LoadWatchList()
	watchers := data.LoadWatchers()

	// Загрузить/создать словарь статусов
	statusMap := data.LoadStatusMap()
	if len(statusMap) == 0 {
		statusMap = data.GenerateStatusMap(watchlist)
	}

	// Стартуем бота
	chat.StartUp(token)

	// Сигнализируем о включении
	chat.StartupNotify(watchers, watchlist)
	go chat.ListenUpdates()
	// Проверяем статус хостов и уведомляем при изменениях
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			chat.UpdateStatusMapAndAlert(statusMap, watchers)
		}
	}

}
