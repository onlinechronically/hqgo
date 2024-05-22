package hqtrivia

import "fmt"

type Show struct {
	PrizeCurrency string
	PrizeTotal    int
	ShowID        int
	Type          string
	Live          bool
	SocketURL     string
	BroadcastID   int
}

func (u User) Schedule() (showArr []Show, _ error) {
	scheduleData, _, err := u.request("GET", fmt.Sprintf("%s/shows/schedule", apiURL), nil)
	if err != nil {
		return nil, err
	}
	for _, showData := range scheduleData["shows"].([]interface{}) {
		showData := showData.(map[string]interface{})
		liveData, isLive := showData["live"].(map[string]interface{})
		var currency, gameType, socketUrl string
		var broadcastId, showId, prizeCents int
		currencyRaw := showData["currency"]
		if currencyRaw != nil {
			currency = currencyRaw.(string)
		} else {
			currency = ""
		}
		gameTypeRaw := showData["gameType"]
		if gameTypeRaw != nil {
			gameType = gameTypeRaw.(string)
		} else {
			gameType = ""
		}
		if isLive {
			socketUrlRaw := liveData["socketUrl"]
			if socketUrlRaw != nil {
				socketUrl = socketUrlRaw.(string)
			} else {
				socketUrl = ""
			}
			broadcastIdRaw := liveData["broadcastId"]
			if broadcastIdRaw != nil {
				broadcastId = int(broadcastIdRaw.(float64))
			} else {
				broadcastId = -1
			}
		} else {
			socketUrl = ""
			broadcastId = -1
		}
		prizeCentsRaw := liveData["prizeCents"]
		if prizeCentsRaw != nil {
			prizeCents = int(prizeCentsRaw.(float64))
		} else {
			prizeCents = -1
		}
		showIdRaw := liveData["showId"]
		if showIdRaw != nil {
			showId = int(showIdRaw.(float64))
		} else {
			showId = -1
		}
		showArr = append(showArr, Show{PrizeCurrency: currency, PrizeTotal: prizeCents, ShowID: showId, Type: gameType, Live: isLive, SocketURL: socketUrl, BroadcastID: broadcastId})
	}
	return showArr, nil
}

func (u User) GetLiveShows() (allShows []Show, _ error) {
	allShows, err := u.Schedule()
	if err != nil {
		return nil, err
	}
	for _, show := range allShows {
		if show.Live {
			allShows = append(allShows, show)
		}
	}
	return allShows, nil
}
