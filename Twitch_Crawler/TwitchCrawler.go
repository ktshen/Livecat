package main

import (
	"fmt"
	//"mongogogo"
	//"douyustreams_crawler"
	"twitchstreams_crawler"
	//"time"
	//"mongogo"
)

func main() {
	//var web = new(mongogo.Web)
	//web.Host = "Twitch"
	//collection := "Platform"
	//mongogo.MongogoWebInitial(web, collection)
	if true {
		twitchstreams_crawler.TwitchCrawlerInitial()
	}
	//douyustreams_crawler.DouyuCrawlerInitial()
	//streams_crawler.GetAllStreams()
	//mongogogo.MongoDBInsertOne()
	//mongogogo.MongoDBInsertOne("AAA", "test")
	fmt.Println("%%%%%% Your motion is compeleted %%%%%%")
}
