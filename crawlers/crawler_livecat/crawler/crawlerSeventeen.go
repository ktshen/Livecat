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

// SeventeenInitial is initial to crawler 17live
func SeventeenInitial() {
	log.Println("START [17]")
	crawlerSeventeen("https://17.live/suggested")
	log.Println("END [17]")

}

func crawlerSeventeen(url string) {
	service := getService(portSeventeen)
	defer service.Stop()
	wd := getRemote(portSeventeen)
	defer wd.Quit()

	doc := getSeventeenDoc(wd, url)
	getSeventeenDATA(doc)

}

func getSeventeenDoc(wd selenium.WebDriver, url string) (doc *goquery.Document) {

	wd.Get(url)
	wd.MaximizeWindow("")
	time.Sleep(5 * time.Second)
	wd.ExecuteScript("window.scrollTo(0,5000)", nil)
	time.Sleep(10 * time.Second)

	var frameHTML string
	frameHTML, err := wd.PageSource()
	handleError(err, "[17] PageSource err")

	doc, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
	handleError(err, "[17] NewDocumentFromReader err")
	return
}

func getSeventeenDATA(doc *goquery.Document) {
	// var db *mgo.Database
	doc.Find("div.Grid__Row-eThVWD").Find("div.LiveStreamBlock__Block-bdJheI").Each(func(i int, s *goquery.Selection) {
		title, _ := s.Find("a").Attr("title")
		localTime := getLocalTimeVersion2()
		videourl, _ := s.Find("a").Attr("href")
		viewers := s.Find("a").Find("span.Msg-eAZFzz").Text()
		viewers = strings.Replace(viewers, ",", "", -1)
		viewersInt, _ := strconv.Atoi(viewers)
		thumbnails, _ := s.Find("a").Find("div.LiveStreamBlock__AvatarWrapper-CSDSj").Attr("style")

		thumbnails = seventeenThumbnailsFormat(thumbnails)

		mongoDBData := &mongogo.MongoDB{
			Title:            title,
			Description:      "",
			Platform:         "17直播",
			VideoID:          "",
			Host:             title,
			Status:           "live",
			Thumbnails:       thumbnails,
			Published:        localTime,
			Tags:             "",
			GeneralTag:       "",
			Timestamp:        localTime,
			Language:         "",
			ViewCount:        viewersInt,
			Viewers:          viewersInt,
			VideoURL:         "https://17.live" + videourl,
			VideoEmbedded:    "",
			ChatRoomEmbedded: "",
			Channel:          "",
		}
		// if i%100 == 0 {
		// 	db = mongogo.GetService("Crawler")
		// }
		// mongogo.MongogoInitial(db, "Livestreams", *mongoDBData)
		sendToElasticByHTTPPost(*mongoDBData)

		//log.Println("i is", i)
		//log.Printf("%+v", *mongoDBData)

	})
}

func seventeenThumbnailsFormat(thumbnails string) string {
	thumbnails = thumbnails[strings.Index(thumbnails, "url(\"")+5:]
	thumbnails = thumbnails[:strings.Index(thumbnails, "url(\"")-4]
	return thumbnails
}
