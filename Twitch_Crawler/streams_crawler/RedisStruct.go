package streams_crawler

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
