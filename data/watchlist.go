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

func AddToWatchList(newIp string) (bool, error) {
	watchlist := LoadWatchList()

	// Проверить ip на корректность
	if net.ParseIP(newIp) == nil {
		return false, nil
	}
	var check bool = true
	for _, ip := range watchlist {
		if ip == newIp {
			check = check && false
		}
	}

	if check {
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

	}
	return true, nil
}
