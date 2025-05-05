package data

import (
	"encoding/json"
	"log"
	"net"
	"os"
)

func LoadWatchList() []string {
	watchlist := []string{}

	var data, err = os.ReadFile("storage/watchlist.json")
	if err == nil {
		JsonErr := json.Unmarshal(data, &watchlist)
		if JsonErr != nil {
			log.Printf(JsonErr.Error())
		}
	}
	return watchlist
}

func RemoveFromWatchList(targetIP string) (bool, error) {
	watchlist := LoadWatchList()
	updated := []string{}
	found := false

	for _, ip := range watchlist {
		if ip == targetIP {
			found = true
			continue
		}
		updated = append(updated, ip)
	}

	if !found {
		return false, nil
	}

	jsonData, err := json.Marshal(updated)
	if err != nil {
		log.Println("Ошибка сериализации при удалении IP:", err)
		return false, err
	}

	err = os.WriteFile("storage/watchlist.json", jsonData, 0644)
	if err != nil {
		log.Println("Ошибка записи файла при удалении IP:", err)
		return false, err
	}

	return true, nil
}

func AddToWatchList(newIp string) (bool, error) {
	watchlist := LoadWatchList()

	// Проверить ip на корректность
	if net.ParseIP(newIp) == nil {
		return false, nil
	}
	for _, ip := range watchlist {
		if ip == newIp {
			return false, nil
		}
	}

	watchlist = append(watchlist, newIp)

	jsonData, err := json.Marshal(watchlist)
	if err != nil {
		log.Println("Ошибка обработки списка наблюдаемых!")
		return false, err
	}

	err = os.WriteFile("storage/watchlist.json", jsonData, 0644)
	if err != nil {
		log.Println("Ошибка записи списка наблюдаемых!" + err.Error())
		return false, err
	}

	return true, nil
}
