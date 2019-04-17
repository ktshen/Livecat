package twitchstreams_crawler

import (
	"bytes"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"mongodbStruct"
	"mongogogo"
	"net/http"
	"time"
)

const (
	MongoDBNameTest       = "Stream"
	MongoDBName           = "Crawler"
	MongoDBCollectionName = "Livestreams"
)

func postMongoInitial(client *mongo.Client, streams mongodbStruct.MongoDB, Host mongodbStruct.Web) mongodbStruct.MongoDB {
	//transform the DATA to JSON DATA
	streams_encode, _ := json.Marshal(streams)
	var decode mongodbStruct.MongoDB
	json.Unmarshal(streams_encode, &decode)
	mongogogo.MongoInsertOne(client, MongoDBName, MongoDBCollectionName, decode)
	//mongogo.MongogoInitial(decode, MongoDBName2, "Livestreams")
	//web_encode, _ := json.Marshal(Host)
	//var _decode mongogo.Web
	//json.Unmarshal(web_encode, &_decode)
	//mongogo.MongogoWebInitial(_decode, "Host")
	return decode
}

func postESInitial(ESurl string, streams mongodbStruct.MongoDB) {
	counter := 1
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	data, _ := json.Marshal(streams)
	body := bytes.NewBuffer(data)
	req, err := http.NewRequest("POST", ESurl, body)
	req.Header.Set("Content-Type", "application/json")
	handleError(err, "Request Post")
	_, err = client.Do(req)
	if err != nil {
		log.Printf("Post to Elestic Search Error %d times !\n", counter)
		counter++
		time.Sleep(time.Second * 8)
		_, err = client.Do(req)
		if err != nil {
			log.Printf("Post to Elestic Search Error %d times !\n", counter)
			counter++
			time.Sleep(time.Second * 8)
			_, err = client.Do(req)
			if err != nil {
				log.Printf("Post to Elestic Search Error %d times !\n", counter)
				return
			}
		}
	}
}

//func dataTransform(streams mongodbStruct.MongoDB) url.Values {
// data := url.Values{}
// data.Set("Method", "post")
// data.Add("Host", streams.Host) //Stream Host
// //data.Add("Published", streams.Published)               //Started time
// data.Add("Thumbnails", streams.Thumbnails)             //Thumbnails
// data.Add("platform", "Twitch")                         //Stream Platform
// data.Add("Viewers", string(streams.Viewers))           //Stream number of viwers
// data.Add("VideoURL", streams.VideoURL)                 //Stream player
// data.Add("Tags", streams.Tags)                         //Game name
// data.Add("Language", streams.Language)                 //Stream Language
// data.Add("Title", streams.Title)                       //Stream title
// data.Add("Status", streams.Status)                     //Stream status (Live or Offline)
// data.Add("VideoId", streams.VideoID)                   //Video ID
// data.Add("GeneralTag", streams.GeneralTag)             //Game's thumbnail
// data.Add("VideoEmbedded", streams.VideoEmbedded)       //Stream embedded URL
// data.Add("Channel", streams.Channel)                   //Stream channel
// data.Add("ChatRoomEmbedded", streams.ChatRoomEmbedded) //Stream chatroom embedded URL
//data.Add("Timestamp", string(streams.Timestamp))       //Timestamp
// streams.Host = DATA.User_Name                                //Stream Host
// streams.Published = t_                                       //Started Time
// streams.Platform = "Twitch"                                  //Platform
// streams.Thumbnails = r.Replace(DATA.Thumbnail_URL)           //Streams Thumbnail
// streams.Viewers = DATA.Viewer_Count                          //Stream Viwers
// streams.VideoURL = string(Twitch_Url + DATA.User_Name)       //Twitch Play
// streams.Tags, streams.GeneralTag = GetGameName(DATA.Game_ID) //Game Name and Game's Thumbnail
// streams.Language = DATA.Language                             //Language
// streams.Title = DATA.Title                                   //Stream title
// streams.Status = DATA.Type                                   //Live or Offline
// streams.VideoId = DATA.ID                                    //Video ID
// streams.VideoEmbedded = embed                                //Stream Embedded URL
// streams.Channel = DATA.User_Name                             //Stream Channel Name
// streams.ChatRoomEmbedded = chatembed                         //Stream Chatroom Embedded URL
// streams.Timestamp = updateTime

//return data
//}

////////////////////////////////////////////////////////////////
/*package twitchstreams_crawler

import (
	"encoding/json"
	//"log"
	"redigogo"
	"strconv"
)

const LOCALIP = "127.0.0.1"
const SERVERIP = "140.115.153.185"

type Streamer struct {
	Host     string
	Category string
	Time     string
	Platform string
	Link     string
	Photo    string
	Views    int
}

func postRedisInitial(streamer Streamer) {
	//log.Println(streamer)
	streamer_encode, _ := json.Marshal(streamer)
	redisChannel0(streamer)
	redigogo.Lpush(SERVERIP, "1", streamer.Category, streamer_encode)
	redigogo.Lpush(SERVERIP, "2", streamer.Host, streamer_encode)
	redigogo.Lpush(SERVERIP, "3", streamer.Platform, streamer_encode)
	redigogo.Lpush(SERVERIP, "4", strconv.Itoa(streamer.Views), streamer_encode)
	//redigogo.LpushString(SERVERIP, "5", "test", streamer.Photo)

	var decode Streamer
	json.Unmarshal(streamer_encode, &decode)
	//log.Println(decode)
}

func redisChannel0(streamer Streamer) {
	streamer_encode, _ := json.Marshal(streamer)
	redigogo.Sadd(SERVERIP, "0", "Host", streamer.Host)
	redigogo.Sadd(SERVERIP, "0", "Category", streamer.Category)
	redigogo.Sadd(SERVERIP, "0", "Platform", streamer.Platform)
	redigogo.Sadd(SERVERIP, "0", "Link", streamer.Link)
	redigogo.Lpush(SERVERIP, "0", "All", streamer_encode)
	redigogo.LpushInt(SERVERIP, "0", "Views", streamer.Views)
}
*/
//Redis Struct
