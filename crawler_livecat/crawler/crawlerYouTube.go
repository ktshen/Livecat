package crawler

import (
	"bytes"
	"crawler_livecat/mongogo"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

var catogory string

// YouTubeInitial ...
func YouTubeInitial() {
	log.Println("START [YouTube]")
	urlList := []string{
		"https://www.youtube.com/playlist?list=PLU12uITxBEPFJGVb2zSgCaWvMBe7vHonB", // 直播
		//"https://www.youtube.com/playlist?list=PL3ZQ5CpNulQkLQffQzV1TgCMKCUtK8fuA", // 新聞
		"https://www.youtube.com/playlist?list=PLiCvVJzBupKlQ50jZqLas7SAztTMEYv1f", // 遊戲
		"https://www.youtube.com/playlist?list=PL8fVUTBmJhHJrxHg_uNTMyRmsWbFltuQV", // 運動
		"https://www.youtube.com/playlist?list=PL57quI9usf_sQMlMeQrUr5O_pCncZilx3", // 科技
		"https://www.youtube.com/playlist?list=PLIFqWCuxNyoiKKthaTBqjIH6m6A9INomt", // 動物
		"https://www.youtube.com/playlist?list=PLU12uITxBEPFnoOrc_w0oJL6CEeKRhpcb", // 行動應用程式
		"https://www.youtube.com/playlist?list=PLU12uITxBEPHvBRHoUt-fzKr8Iz1HpfUC", // 網路攝影機直播
	}
	for _, url := range urlList {

		crawlerYouTube(url)
		log.Println(catogory)
	}
	log.Println("END [YouTube]")

}

func crawlerYouTube(url string) {
	service := getService(portYouTube)
	defer service.Stop()
	wd := getRemote(portYouTube)
	defer wd.Quit()

	doc := getYouTubeDoc(wd, url)
	getYouTubeDATA(wd, doc)

}
func getYouTubeDoc(wd selenium.WebDriver, url string) (doc *goquery.Document) {
	wd.Get(url)
	time.Sleep(1 * time.Second)

	var frameHTML string
	frameHTML, err := wd.PageSource()
	handleError(err, "[Youtube] PageSource err")
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
	handleError(err, "[Youtube New] NewDocumentFromReader err")

	return
}

func getYouTubeDATA(wd selenium.WebDriver, doc *goquery.Document) {
	// var db *mgo.Database
	catogory = doc.Find("ytd-page-manager#page-manager").Find("ytd-playlist-sidebar-renderer.style-scope").Find("h1#title").Find("a.yt-simple-endpoint").Text()
	doc.Find("ytd-item-section-renderer.style-scope").Find("ytd-playlist-video-list-renderer.style-scope").Find("ytd-playlist-video-renderer.style-scope").Each(func(i int, s *goquery.Selection) {
		//thumbnails, _ := s.Find("div#content").Find("yt-img-shadow.style-scope").Find("img#img").Attr("src")
		videourl, _ := s.Find("div#content").Find("a").Attr("href")

		title, _ := s.Find("div#content").Find("div#meta").Find("span#video-title").Attr("title")
		host := s.Find("div#content").Find("yt-formatted-string.style-scope").Find("a.yt-simple-endpoint").Text()

		localTime := getLocalTimeVersion2()
		streamer := getYouTubeDoc(wd, "https://www.youtube.com"+videourl)
		videoid, videoembedded, chatroomembedded := youtubeEmbedded(videourl)
		viewersStr := streamer.Find("ytd-page-manager#page-manager").Find("div#columns").Find("div#primary").Find("div#primary-inner").Find("div#info").Find("div#info-contents").Find("ytd-video-primary-info-renderer.style-scope").Find("div#info").Find("yt-view-count-renderer").Find("span.view-count").Text()
		viewers := youtubeViewFormat(viewersStr)
		thumbnails := youtubeThumbnailsFormat(videoid)
		mongoDBData := &mongogo.MongoDB{
			Title:            title,
			Description:      "",
			Platform:         "YouTube",
			VideoID:          videoid,
			Host:             host,
			Status:           "live",
			Thumbnails:       thumbnails,
			Published:        localTime,
			Tags:             "",
			GeneralTag:       "",
			Timestamp:        localTime,
			Language:         "",
			ViewCount:        viewers,
			Viewers:          viewers,
			VideoURL:         "https://www.youtube.com" + videourl,
			VideoEmbedded:    videoembedded,
			ChatRoomEmbedded: chatroomembedded,
			Channel:          videoid,
		}
		// if i%100 == 0 {
		// 	db = mongogo.GetService("Crawler")
		// }
		sendToElasticByHTTPPost(*mongoDBData)
		// mongogo.MongogoInitial(db, "Livestreams", *mongoDBData)
		//log.Println("i is", i)
		//log.Printf("%+v", *mongoDBData)
	})
}

func youtubeViewFormat(str string) (viewers int) {
	//log.Println(str)
	first := strings.Index(str, "：")
	end := strings.Index(str, "人")
	if first == -1 {
		first = strings.Index(str, ":")
	}

	if first == -1 && end == -1 {
		return 0
	}

	if first == -1 {
		first = 0
	} else {
		first = first + 3
	}
	if end == -1 {
		end = len(str) - 3
	} else {
		end--
	}
	str = str[first:end]
	str = strings.Replace(str, ",", "", -1)
	viewers, _ = strconv.Atoi(str)
	//log.Println(viewers)
	return
}

func youtubeEmbedded(url string) (videoid string, videoembedded string, chatroomembedded string) {
	url = url[strings.Index(url, "/watch?v=")+9:]
	videoid = url[:strings.Index(url, "&")]
	videoembedded = "https://www.youtube.com/embed/" + videoid + "?rel=0&autoplay=1"
	chatroomembedded = "https://www.youtube.com/live_chat?v=" + videoid + "&is_popout=1&embed_domain=dlivecat.com"
	return
}

func youtubeThumbnailsFormat(videoid string) (thumbnails string) {
	thumbnails = "https://i.ytimg.com/vi/" + videoid + "/hqdefault_live.jpg?sqp=CJCOqOUF-oaymwEZCPYBEIoBSFXyq4qpAwsIARUAAIhCGAFwAQ==&rs=AOn4CLBkQtCV6l_0tb7U2crBDOS60-4y1A"
	return
}
