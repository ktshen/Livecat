package twitchstreams_crawler

import (
	"encoding/json"
	"fmt"
	//"mongodbStruct"
	//"net/http"
	"os"
	"strings"
	"time"
)

const (
	GameURL       = "https://api.twitch.tv/helix/games/top"
	GameArraySize = 1000000
	//TwitchFileJson = "Twitch.Json"
)

var (
	gameName     = new(GameNameFileJson)
	gameArray    [GameArraySize]Data2
	arrayCounter = 0
)

func createFile(ListName string) os.File {
	file, err := os.Create(ListName)
	handleError(err, "")
	return *file
}

func CloseFile(file os.File) {
	file.Close()
}

func WriteGameList(file os.File, gameInformation *GameNameFileJson) {
	file_encode, _ := json.Marshal(gameInformation)
	file.Write(file_encode)
}

func GetAllGameInformation() {
	f := createFile(GameInformationJsonFile)
	defer CloseFile(f)
	fmt.Println("----> Initialize the game list. <----")
	url := GameURL
	var i = 1

	_, BeforeM := UpdateTime("checkM")
	_, BeforeS := UpdateTime("checkS")

	for i < 8000 {
		fmt.Println("------->", i)
		api := GetAPIClient(url)
		var resp_decode GameNameJson
		json.Unmarshal(api, &resp_decode)
		for _, DATA := range resp_decode.DATA {
			DATA = GameInformationJson(DATA)
			gameArray[arrayCounter] = DATA
			arrayCounter++
			//fmt.Println(DATA)
		}
		parameter := "?after=" + resp_decode.Pagination.Cursor
		url = GameURL + parameter
		i++
		if i%29 == 0 {
			_, AfterM := UpdateTime("checkM")
			_, AfterS := UpdateTime("checkS")
			ts := AfterS - BeforeS
			tm := AfterM - BeforeM
			if ts < 50 {
				if tm < 1 {
					ts = 60 - ts
					time.Sleep(time.Duration(ts) * time.Second)
				}
			}
			_, BeforeS = UpdateTime("checkS")
			_, BeforeM = UpdateTime("checkM")
		}
	}
	//time.Sleep(time.Second * 50)
	gameName.DATA = gameArray
	WriteGameList(f, gameName)
	fmt.Println("------> Initialize completed. <------")
}

func GameInformationJson(DATA Data2) Data2 {
	r := strings.NewReplacer("{width}x{height}", "480"+"x"+"480")
	DATA.Box_art_url = r.Replace(DATA.Box_art_url)
	return DATA
}
