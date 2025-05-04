package checker

import (
	"log"
	"time"

	"github.com/go-ping/ping"
)

// Проверяет статус хоста, возвращает true, если хост ответил хотя бы 1 раз
func CheckICMP(ip string) bool {
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
