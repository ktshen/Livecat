package twitchstreams_crawler

import (
	"Twitch_Crawler/mongodbStruct"
	"Twitch_Crawler/mongogogo"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	//TTL                     = int32(3000) //ExpireAfterSeconds
	IMAGE_SIZE_WIGTH        = "480"                   //width
	IMAGE_SIZE_LENGTH       = "360"                   //Lenth
	GameInformationJsonFile = "Twitch_Game_List.json" //File name
	STREAM_URL              = "https://api.twitch.tv/kraken/streams/"
	Twitch_Url              = "https://www.twitch.tv/"                        //For generating the stream's URL
	Client_ID               = "a63i69mhmzr6hdvxy3pthw8d2nrlc0"                //Client-ID for developing Twitch program
	SecondClient_ID         = "kw1hxsbf7rc8fh8dojm6jjggy7p6zj"                //The second Client-ID for updating the gameInformation
	Initialurl              = "https://api.twitch.tv/helix/streams?first=100" //For TwitchAPI function
	AllUrl                  = "https://api.twitch.tv/kraken/streams/summary"  //For GetAllStreamsCount function
	Embed                   = "https://player.twitch.tv/?channel="            //For generating embeddeed streams' URL
	ESURL                   = "http://120.126.16.88:17777/add"                //For posting data to the Elestic Search server
	ChatEmbed               = "https://www.twitch.tv/embed/"                  //"https://www.twitch.tv/embed/hebo/chat"
)

var (
	ttl                  = int32(1800)
	count                = 1
	test                 = 1
	gameInformationArray []Data2
)

// Run the main program to start the twitch crawler
// GameListToArray function converts Twitch_Game_List.json to slice
// GetAllStreamsCount function gets all streams on Twitch and count how many times we cursor the API
// PartialUpdate function starts to update the top 1,500 streams on Twitch
// TwitchAPI function get the top 100 streams on Twitch and cursors it
// CreateIndex function set the expiration on MongoDB
func TwitchCrawlerInitial() {
	fmt.Println("Convert Json file into this crawler.")
	GameListToArray()
	fmt.Println("Twitch Crawler is running ... ")
	ReqCount := GetAllStreamsCount()
	//CreateIndex() //Set expire on MongoDB
	defer ClearSlice()
	go PartialUpdate()
	TwitchAPI(Initialurl, ReqCount, "Initial")
}

func CreateIndex() {
	MongoClient := mongogogo.MongoConnect()
	mongogogo.PopulateIndex(MongoDBName, MongoDBCollectionName, MongoClient, ttl)
	err := mongogogo.MongoDisconnect(MongoClient)
	handleError(err, "")
}

func PartialUpdate() {
	url := Initialurl
	for o := 0; o < 8; o++ {
		time.Sleep(time.Minute * 5)
		fmt.Printf("---> Now start partial update. <--- %d times !", o)
		url = TwitchAPI(url, 15, "Update") //partial update for top 2,000 streams
	}
}

func GetAPIClient(url string, client_id string) []byte {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest("GET", url, nil)
	handleError(err, "")
	req.Header.Add("Client-ID", client_id)

	resp, err := client.Do(req)
	for err != nil {
		time.Sleep(time.Second * 8)
		resp, err = client.Do(req)
	}
	content, err := ioutil.ReadAll(resp.Body)
	handleError(err, "Read resp content")
	return content
}

func GetAllStreamsCount() int {
	api := GetAPIClient(AllUrl, Client_ID)
	var resp_decode Summary
	json.Unmarshal(api, &resp_decode)
	//fmt.Println(resp_decode.Channels)
	ReqCount := (resp_decode.Channels / 100) + 1
	return ReqCount
}

func TwitchAPI(url string, ReqCount int, selectIndex string) string {
	counter := 1
	_, BeforeM := UpdateTime("checkM")
	_, BeforeS := UpdateTime("checkS")
	for counter < ReqCount {
		switch selectIndex {
		case "Initial":
			ReqCount = GetAllStreamsCount()
			fmt.Printf("%10s    %3d   %3d\n", selectIndex, counter, ReqCount)
		case "Update":
			ReqCount = 15
			fmt.Printf("%10s    %3d   %3d\n", selectIndex, counter, ReqCount)
		}
		//fmt.Println(url)
		api := GetAPIClient(url, Client_ID)
		var resp_decode TwitchJson
		json.Unmarshal(api, &resp_decode)

		//MongoClient := mongogogo.MongoConnect()

		for _, DATA := range resp_decode.Data {
			json, _ := TwitchJsonToMongoJson(DATA)
			//postMongoInitial(MongoClient, json, WebJSON)
			postESInitial(ESURL, json)
		}
		//err := mongogogo.MongoDisconnect(MongoClient)
		//handleError(err, "")

		parameter := "&after=" + resp_decode.Pagination.Cursor
		url = Initialurl + parameter
		counter++

		if counter%15 == 1 {
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
	return url
}

func GameListToArray() {
	len := 0
	gameInformationArray = make([]Data2, len)
	file, err := os.Open(GameInformationJsonFile)
	handleError(err, "Open file.")
	JsonFile, err := ioutil.ReadAll(file)
	handleError(err, "ReadAll")

	var jsonFile = new(GameNameJson)
	err = json.Unmarshal(JsonFile, &jsonFile)
	handleError(err, "Unmarshal")
	for _, DATA := range jsonFile.DATA {
		gameInformationArray = append(gameInformationArray, DATA)
	}
}

func ClearSlice() {
	gameInformationArray = nil
}

func GetGameName(Game_ID string) (string, string) {
	for _, i := range gameInformationArray {
		if Game_ID == i.Id {
			return i.Name, i.Box_art_url
		}
	}
	return " ", " "
}

func UpdateTime(format string) (time.Time, int) {
	switch format {
	case "checkM":
		updateTimeM := time.Now().Minute()
		return time.Now(), updateTimeM
	case "checkS":
		updateTimeS := time.Now().Second()
		return time.Now(), updateTimeS
	case "update":
		updateTime := time.Now()
		local, err := time.LoadLocation("Local")
		handleError(err, "Get local time err")

		return updateTime.In(local), 0
	}
	return time.Now(), 0
}

func TwitchJsonToMongoJson(DATA Data) (mongodbStruct.MongoDB, mongodbStruct.Web) {
	r := strings.NewReplacer("{width}x{height}", IMAGE_SIZE_WIGTH+"x"+IMAGE_SIZE_LENGTH)
	//t := strings.Replace(DATA.Started_At, "T", " ", -1)
	//t_ := strings.Replace(t, "Z", "", -1)
	embed := Embed + strings.ToLower(DATA.User_Name)
	chatembed := ChatEmbed + strings.ToLower(DATA.User_Name) + "/chat"

	var Web = new(mongodbStruct.Web)
	var streams = new(mongodbStruct.MongoDB)
	updateTime, _ := UpdateTime("update")

	Web.Host = DATA.User_Name
	streams.Host = DATA.User_Name                                //Stream Host
	streams.Published = DATA.Started_At                          //Started Time
	streams.Platform = "Twitch"                                  //Platform
	streams.Thumbnails = r.Replace(DATA.Thumbnail_URL)           //Streams Thumbnail
	streams.Viewers = DATA.Viewer_Count                          //Stream Viwers
	streams.VideoURL = string(Twitch_Url + DATA.User_Name)       //Twitch Play
	streams.Tags, streams.GeneralTag = GetGameName(DATA.Game_ID) //Game Name and Game's Thumbnail
	streams.Language = DATA.Language                             //Language
	streams.Title = DATA.Title                                   //Stream title
	streams.Status = DATA.Type                                   //Live or Offline
	streams.VideoID = DATA.ID                                    //Video ID
	streams.VideoEmbedded = embed                                //Stream Embedded URL
	streams.Channel = DATA.User_Name                             //Stream Channel Name
	streams.ChatRoomEmbedded = chatembed                         //Stream Chatroom Embedded URL
	streams.Timestamp = updateTime                               //Mongodb Timestamp
	streams.CreatedAt = updateTime                               //Expire Data after 25 minutes from now
	return *streams, *Web
}
