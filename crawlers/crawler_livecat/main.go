package main

import (
	"crawler_livecat/crawler"
	"time"
)

func main() {
	go crawler.WatermelonInitial()
	go crawler.SeventeenInitial()
	youtubeList := []string{
		"https://www.youtube.com/playlist?list=PLU12uITxBEPFJGVb2zSgCaWvMBe7vHonB", // 直播
		//"https://www.youtube.com/playlist?list=PL3ZQ5CpNulQkLQffQzV1TgCMKCUtK8fuA", // 新聞
		"https://www.youtube.com/playlist?list=PLiCvVJzBupKlQ50jZqLas7SAztTMEYv1f", // 遊戲
		"https://www.youtube.com/playlist?list=PL8fVUTBmJhHJrxHg_uNTMyRmsWbFltuQV", // 運動
		"https://www.youtube.com/playlist?list=PL57quI9usf_sQMlMeQrUr5O_pCncZilx3", // 科技
		"https://www.youtube.com/playlist?list=PLIFqWCuxNyoiKKthaTBqjIH6m6A9INomt", // 動物
		"https://www.youtube.com/playlist?list=PLU12uITxBEPFnoOrc_w0oJL6CEeKRhpcb", // 行動應用程式
		"https://www.youtube.com/playlist?list=PLU12uITxBEPHvBRHoUt-fzKr8Iz1HpfUC", // 網路攝影機直播
	}
	for _, url := range youtubeList {
		go crawler.YouTubeInitial(url)
	}

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
			for _, url := range youtubeList {
				go crawler.YouTubeInitial(url)
			}
		}
	}
}
