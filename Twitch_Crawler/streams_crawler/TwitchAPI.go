package streams_crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const IMAGE_SIZE_WIGTH = "320"  //寬
const IMAGE_SIZE_LENGTH = "150" //高
const TWITCH_URL = "https://www.twitch.tv/"

var count = 1

type TwitchJson struct {
	Data       []Data
	Pagination Pagination
}
type Data struct {
	ID            string
	User_ID       string
	User_Name     string
	Game_ID       string
	Community_IDs []string
	Type          string
	Tital         string
	Viewer_Count  int
	Started_At    string
	Language      string
	Thumbnail_URL string
	Tag_IDs       []string
}

type Community_ids struct {
	id []string
}

type Pagination struct {
	Cursor string
}

type GAME_ID_JSON struct {
	DATA []DATA2
}
type DATA2 struct {
	Id          string
	Name        string
	Box_art_url string
}

const Twitch_Url = "https://www.twitch.tv/"
const Client_ID = "kw1hxsbf7rc8fh8dojm6jjggy7p6zj"

func TwitchAPI() {
	url := "https://api.twitch.tv/helix/streams?first=100"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	handleError(err, "")
	req.Header.Add("Client-ID", Client_ID)

	resp, err := client.Do(req)
	handleError(err, "")

	content, err := ioutil.ReadAll(resp.Body)
	handleError(err, "")

	var resp_decode TwitchJson
	json.Unmarshal(content, &resp_decode)
	//fmt.Println(string(content))

	for _, DATA := range resp_decode.Data {
		redisjson := TwitchJsonToRedisJson(DATA)
		//fmt.Println(redisjson)
		postRedisInitial(redisjson)

	}

	return
}

func TwitchJsonToRedisJson(DATA Data) Streamer {
	r := strings.NewReplacer("{width}x{height}", IMAGE_SIZE_WIGTH+"x"+IMAGE_SIZE_LENGTH)
	t := strings.Replace(DATA.Started_At, "T", " ", -1)
	t_ := strings.Replace(t, "Z", "", -1)

	fmt.Println(count)
	count = count + 1

	var streamer = new(Streamer)
	streamer.Host = DATA.User_Name
	streamer.Category = DATA.Game_ID
	streamer.Time = t_
	streamer.Platform = "Twitch"
	streamer.Photo = r.Replace(DATA.Thumbnail_URL)
	streamer.Views = DATA.Viewer_Count
	streamer.Link = string(Twitch_Url + DATA.User_Name)
	streamer.Category = GetGameName(streamer.Category)

	return *streamer
}

func GetGameName(Game_ID string) (GameName string) {
	var GetGameName string
	url := "https://api.twitch.tv/helix/games?id=game-id"
	r := strings.NewReplacer("game-id", Game_ID)
	findGameNameUrl := r.Replace(url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", findGameNameUrl, nil)
	handleError(err, "")
	req.Header.Add("Client-ID", Client_ID)
	resp, err := client.Do(req)
	handleError(err, "")

	content, err := ioutil.ReadAll(resp.Body)
	handleError(err, "")

	var resp_decode GAME_ID_JSON
	json.Unmarshal(content, &resp_decode)

	for _, DATA := range resp_decode.DATA {
		//fmt.Println(DATA.Name)
		GetGameName = DATA.Name
	}

	return GetGameName
}
