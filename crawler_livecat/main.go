package main

import (
	"crawler_livecat/crawler"
	"time"
)

func main() {

	freqYouTube := time.NewTicker(time.Minute * 10)
	freqSeventeen := time.NewTicker(time.Minute * 10)
	freqWatermelon := time.NewTicker(time.Minute * 10)
	for {
		select {
		case <-freqWatermelon.C:
			go crawler.WatermelonInitial()
		case <-freqSeventeen.C:
			go crawler.SeventeenInitial()
		case <-freqYouTube.C:
			go crawler.YouTubeInitial()
		}
	}

}
