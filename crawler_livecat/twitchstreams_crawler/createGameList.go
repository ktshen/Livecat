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
	f := createFile("Test.Json")
	defer CloseFile(f)
	fmt.Println("----> Initialize the game list. <----")
	url := GameURL
	var i = 1

	for i < 5001 {
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
			time.Sleep(time.Second * 58)
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

	//fmt.Println(gameInformation.data)
	return DATA
}
