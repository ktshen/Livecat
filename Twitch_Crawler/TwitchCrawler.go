package main

import (
	"fmt"
	//"mongogogo"
	//"douyustreams_crawler"
	"Twitch_Crawler/twitchstreams_crawler"
	//"time"
	//"mongogo"
)

func main() {

	fmt.Println("Start process ! ")

	for {
		twitchstreams_crawler.TwitchCrawlerInitial()
	}

}
