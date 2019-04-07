package twitchstreams_crawler

import (
	"crawler_livecat/mongodbStruct"
	"crawler_livecat/mongogogo"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	//"reflect"
	"strings"
	"time"
	//"unsafe"
)

const (
	IMAGE_SIZE_WIGTH        = "480" //寬
	IMAGE_SIZE_LENGTH       = "360" //高
	GameInformationJsonFile = "Twitch_Game_List.json"
	STREAM_URL              = "https://api.twitch.tv/kraken/streams/"
	Time_delay              = 1000
	Twitch_Url              = "https://www.twitch.tv/"
	Client_ID               = "kw1hxsbf7rc8fh8dojm6jjggy7p6zj"
	Initialurl              = "https://api.twitch.tv/helix/streams?first=100"
	AllUrl                  = "https://api.twitch.tv/kraken/streams/summary"
	Embed                   = "https://player.twitch.tv/?channel="
	ChatEmbed               = "https://www.twitch.tv/embed/" //"https://www.twitch.tv/embed/hebo/chat"
)

var (
	count                   = 1
	test                    = 1
	gameInformationArray    [GameArraySize]Data2
	gameInformationArrayTop = 0
)

func TwitchCrawlerInitial() {
	fmt.Println("Convert Json file into this crawler.")
	GameListToArray()
	fmt.Println("Twitch Crawler is running ... ")
	ReqCount := GetAllStreamsCount()
	go TwitchAPI(Initialurl, ReqCount)
	PartialUpdate()
}

func PartialUpdate() {
	for o := 0; o < 50; o++ {
		fmt.Println("-------> Now wait to partial update. <--------")
		time.Sleep(time.Minute * 20)
		TwitchAPI(Initialurl, 20) //partial update for top 2,000 streams
	}
}

func CheckOffLine() {
	StreamOffLine()
}

func GetAPIClient(url string) []byte {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	handleError(err, "")
	req.Header.Add("Client-ID", Client_ID)

	resp, err := client.Do(req)
	handleError(err, "")
	content, err := ioutil.ReadAll(resp.Body)
	handleError(err, "")
	return content
}

func GetAllStreamsCount() int {
	api := GetAPIClient(AllUrl)
	var resp_decode Summary
	json.Unmarshal(api, &resp_decode)

	ReqCount := (resp_decode.Channels / 50) + 1
	return ReqCount
}

func TwitchAPI(url string, ReqCount int) {
	counter := 1
	_, BeforeM := UpdateTime("checkM")
	_, BeforeS := UpdateTime("checkS")
	for counter < ReqCount {
		fmt.Println(url)
		api := GetAPIClient(url)
		var resp_decode TwitchJson
		json.Unmarshal(api, &resp_decode)

		MongoClient := mongogogo.MongoConnect()
		for _, DATA := range resp_decode.Data {
			mongojson, WebJSON := TwitchJsonToMongoJson(DATA)
			postMongoInitial(MongoClient, mongojson, WebJSON)
		}
		err := mongogogo.MongoDisconnect(MongoClient)
		handleError(err, "")

		parameter := "&after=" + resp_decode.Pagination.Cursor
		url = Initialurl + parameter
		counter++

		if counter%25 == 1 {
			_, AfterM := UpdateTime("checkM")
			_, AfterS := UpdateTime("checkS")
			ts := AfterS - BeforeS
			tm := AfterM - BeforeM
			if ts < 50 {
				if tm < 1 {
					time.Sleep(time.Second * 20)
				}
			}
			_, AfterS = UpdateTime("checkS")
			_, AfterM = UpdateTime("checkM")
		}
	}
}

func StreamOffLine() {
	// doc, err := goquery.NewDocument("https://www.twitch.tv/nightblue3")
	// handleError(err, "")
	// doc.Find("div.shell-nav__link shell-nav__link--browse").Find("div.player-streamstatus").Each(func(i int, s *goquery.Selection) {
	// 	fmt.Println("---->")
	// 	status, exist := s.Find("span.player-tip").Attr("data-tip")
	// 	fmt.Println(status, i)
	// 	if exist {
	// 		fmt.Println(status, i)
	// 	} else {
	// 		fmt.Println("Href dose not exist!!")
	// 	}
	// })
}

func GameListToArray() {
	file, err := os.Open(GameInformationJsonFile)
	handleError(err, "Open file.")
	JsonFile, err := ioutil.ReadAll(file)
	handleError(err, "ReadAll")

	var jsonFile = new(GameNameJson)
	err = json.Unmarshal(JsonFile, &jsonFile)
	handleError(err, "Unmarshal")
	for _, DATA := range jsonFile.DATA {
		gameInformationArray[gameInformationArrayTop] = DATA
		gameInformationArrayTop++
	}
}

func GetGameName(Game_ID string) (string, string) {
	for i := 0; i < gameInformationArrayTop; i++ {
		if Game_ID == gameInformationArray[i].Id {
			return gameInformationArray[i].Name, gameInformationArray[i].Box_art_url
		}
	}
	return " ", " "
}

func UpdateTime(format string) (string, int) {
	switch format {
	case "checkM":
		updateTimeM := time.Now().Minute()
		return "0", updateTimeM
	case "checkS":
		updateTimeS := time.Now().Second()
		return "0", updateTimeS
	case "update":
		updateTime := time.Now().Format("2006-01-02 15:04:05")
		return updateTime, 0
	}
	return "0", 0
}

func TwitchJsonToMongoJson(DATA Data) (mongodbStruct.MongoDB, mongodbStruct.Web) {
	r := strings.NewReplacer("{width}x{height}", IMAGE_SIZE_WIGTH+"x"+IMAGE_SIZE_LENGTH)
	t := strings.Replace(DATA.Started_At, "T", " ", -1)
	t_ := strings.Replace(t, "Z", "", -1)
	embed := Embed + strings.ToLower(DATA.User_Name)
	chatembed := ChatEmbed + strings.ToLower(DATA.User_Name) + "/chat"

	var Web = new(mongodbStruct.Web)
	var streams = new(mongodbStruct.MongoDB)
	updateTime, _ := UpdateTime("update")

	Web.Host = DATA.User_Name
	streams.Host = DATA.User_Name                                //Stream Host
	streams.Published = t_                                       //Started Time
	streams.Platform = "Twitch"                                  //Platform
	streams.Thumbnails = r.Replace(DATA.Thumbnail_URL)           //Streams Thumbnail
	streams.Viewers = DATA.Viewer_Count                          //Stream Viwers
	streams.VideoURL = string(Twitch_Url + DATA.User_Name)       //Twitch Play
	streams.Tags, streams.GeneralTag = GetGameName(DATA.Game_ID) //Game Name and Game's Thumbnail
	streams.Language = DATA.Language                             //Language
	streams.Title = DATA.Title                                   //Stream title
	streams.Status = DATA.Type                                   //Live or Offline
	streams.VideoId = DATA.ID                                    //Video ID
	streams.VideoEmbedded = embed                                //Stream Embedded URL
	streams.Channel = DATA.User_Name                             //Stream Channel Name
	streams.ChatRoomEmbedded = chatembed                         //Stream Chatroom Embedded URL
	streams.Timestamp = updateTime                               //Mongodb Timestamp
	return *streams, *Web
}
