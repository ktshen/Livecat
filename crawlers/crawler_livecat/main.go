package main

import (
	"crawler_livecat/crawler"
	"time"
)

func main() {
	go crawler.WatermelonInitial()
	//go crawler.SeventeenInitial()
	go crawler.YouTubeInitial()

	freqYouTube := time.NewTicker(time.Minute * 20)
	freqSeventeen := time.NewTicker(time.Minute * 10)
	freqWatermelon := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-freqWatermelon.C:
			go crawler.WatermelonInitial()
		case <-freqSeventeen.C:
			//go crawler.SeventeenInitial()
		case <-freqYouTube.C:
			go crawler.YouTubeInitial()
		}
	}

}
