package twitchstreams_crawler

import (
	"crawler_livecat/mongodbStruct"
	"encoding/json"

	"go.mongodb.org/mongo-driver/mongo"

	//"log"
	//"mongogo"
	"crawler_livecat/mongogogo"
)

const (
	MongoDBTestName           = "Stream"
	MongoDBTestCollectionName = "TwitchTest"
	MongoDBName               = "Crawler"
	MongoDBCollectionName     = "Livestreams"
)

func postMongoInitial(client *mongo.Client, streams mongodbStruct.MongoDB, Host mongodbStruct.Web) mongodbStruct.MongoDB {
	//Tramsform the DATA to Bson DATA
	//streams_bson_encode, _ := bson.Marshal(streams)
	//var decode_bson MongoDB
	//bson.Unmarshal(streams_bson_encode, &decode_bson)
	//mongogogo.MongoDBInsertOne(MongoDBName, MongoDBCollecttionName, decode_bson)

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
