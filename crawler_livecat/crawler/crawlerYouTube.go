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
	"gopkg.in/mgo.v2"
)

var catogory string

// YouTubeInitial ...
func YouTubeInitial() {
	log.Println("START [YouTube]")
	urlList := []string{"https://www.youtube.com/playlist?list=PL3ZQ5CpNulQkLQffQzV1TgCMKCUtK8fuA", //新聞
		"https://www.youtube.com/playlist?list=PLiCvVJzBupKlQ50jZqLas7SAztTMEYv1f", //遊戲
		"https://www.youtube.com/playlist?list=PL8fVUTBmJhHJrxHg_uNTMyRmsWbFltuQV", //運動
		"https://www.youtube.com/playlist?list=PL57quI9usf_sQMlMeQrUr5O_pCncZilx3", //科技
		"https://www.youtube.com/playlist?list=PLIFqWCuxNyoiKKthaTBqjIH6m6A9INomt", //動物
		"https://www.youtube.com/playlist?list=PLU12uITxBEPFnoOrc_w0oJL6CEeKRhpcb", //行動應用程式
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
	handleError(err, "")
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
	handleError(err, "")

	return
}

func getYouTubeDATA(wd selenium.WebDriver, doc *goquery.Document) {
	var db *mgo.Database
	catogory = doc.Find("ytd-page-manager#page-manager").Find("ytd-playlist-sidebar-renderer.style-scope").Find("h1#title").Find("a.yt-simple-endpoint").Text()
	doc.Find("ytd-item-section-renderer.style-scope").Find("ytd-playlist-video-list-renderer.style-scope").Find("ytd-playlist-video-renderer.style-scope").Each(func(i int, s *goquery.Selection) {
		thumbnails, _ := s.Find("div#content").Find("yt-img-shadow.style-scope").Find("img#img").Attr("src")
		videourl, _ := s.Find("div#content").Find("a").Attr("href")
		title, _ := s.Find("div#content").Find("yt-formatted-string.style-scope").Attr("title")
		host := s.Find("div#content").Find("yt-formatted-string.style-scope").Find("a.yt-simple-endpoint").Text()
		localTime := getLocalTime()
		streamer := getYouTubeDoc(wd, "https://www.youtube.com"+videourl)
		videoid, videoembedded, chatroomembedded := youtubeEmbedded(videourl)
		viewersStr := streamer.Find("ytd-page-manager#page-manager").Find("div#columns").Find("div#primary").Find("div#primary-inner").Find("div#info").Find("div#info-contents").Find("ytd-video-primary-info-renderer.style-scope").Find("div#info").Find("yt-view-count-renderer").Find("span.view-count").Text()
		viewers := youtubeViewFormat(viewersStr)
		mongoDBData := &mongogo.MongoDB{
			Title:            title,
			Description:      "",
			Platform:         "YouTube",
			VideoId:          videoid,
			Host:             host,
			Status:           "live",
			Thumbnails:       thumbnails,
			Published:        "",
			Tags:             "",
			GeneralTag:       "",
			Timestamp:        localTime.String(),
			Language:         "",
			ViewCount:        viewers,
			Viewers:          viewers,
			VideoURL:         "https://www.youtube.com" + videourl,
			VideoEmbedded:    videoembedded,
			ChatRoomEmbedded: chatroomembedded,
			Channel:          videoid,
		}

		if i%100 == 0 {
			db = mongogo.GetService("Crawler")
		}
		mongogo.MongogoInitial(db, "Livestreams", *mongoDBData)
		log.Println("i is", i)
		log.Printf("%+v", *mongoDBData)
	})
}

func youtubeViewFormat(str string) (viewers int) {
	str = str[strings.Index(str, ":")+1 : strings.Index(str, "人")-1]
	str = strings.Replace(str, ",", "", -1)
	viewers, _ = strconv.Atoi(str)
	return
}

func youtubeEmbedded(url string) (videoid string, videoembedded string, chatroomembedded string) {
	url = url[strings.Index(url, "/watch?v=")+9:]
	videoid = url[:strings.Index(url, "&")]
	videoembedded = "https://www.youtube.com/embed/" + videoid + "?rel=0&autoplay=1"
	chatroomembedded = "https://www.youtube.com/live_chat?v=" + videoid + "&is_popout=1&embed_domain=dlivecat.com"
	return
}
