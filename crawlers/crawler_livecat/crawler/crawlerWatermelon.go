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

// WatermelonInitial ...
func WatermelonInitial() {
	log.Println("START [西瓜直播]")
	crawlerWatermelon("https://live.ixigua.com/category/1/109")
	log.Println("END [西瓜直播]")
}

func crawlerWatermelon(url string) {
	service := getService(portWatermelon)
	defer service.Stop()
	wd := getRemote(portWatermelon)
	defer wd.Quit()
	doc := getWatermelonDoc(wd, url)
	getWatermelonDATA(doc)
}

func getWatermelonDoc(wd selenium.WebDriver, url string) (doc *goquery.Document) {
	wd.Get(url)
	time.Sleep(1 * time.Second)
	// wd.MaximizeWindow("")
	elem, err := wd.FindElement(selenium.ByCSSSelector, "div.Feed-footer")
	if err != nil {
		panic(err)
	}
	for countt := 0; countt < 10; countt++ {
		elem.MoveTo(0, 0)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(3 * time.Second)
	var frameHTML string
	frameHTML, err = wd.PageSource()
	handleError(err, "")

	doc, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
	handleError(err, "")

	return
}

func getWatermelonDATA(doc *goquery.Document) {
	var db *mgo.Database
	doc.Find("div.feed-item__list").Find("div.feed-card").Each(func(i int, s *goquery.Selection) {

		link, _ := s.Find("a").Attr("href")
		localTime := getLocalTime()
		viewersStr := strings.Replace(s.Find("span.feed-card__user__count").Text(), "w", "", -1)
		viewers, _ := strconv.ParseFloat(viewersStr, 32)
		thumbnails, _ := s.Find("div.feed-card__cover").Attr("style")
		thumbnails = "http:" + thumbnails[strings.Index(thumbnails, "//"):len(thumbnails)-1]
		thumbnails = strings.Replace(thumbnails, "580x327", "480x360", -1)
		mongoDBData := &mongogo.MongoDB{
			Title:            s.Find("div.feed-card__title").Text(),
			Description:      "",
			Platform:         "西瓜直播",
			VideoId:          "",
			Host:             s.Find("span.feed-card__user__name").Text(),
			Status:           "live",
			Thumbnails:       thumbnails,
			Published:        "",
			Tags:             s.Find("feed-card__activity-tag").Text(),
			GeneralTag:       "",
			Timestamp:        localTime.String(),
			Language:         "Simplified Chinese",
			ViewCount:        int(viewers * 10000),
			Viewers:          int(viewers * 10000),
			VideoURL:         "https://live.ixigua.com" + link,
			VideoEmbedded:    "",
			ChatRoomEmbedded: "",
			Channel:          "",
		}
		if i%100 == 0 {
			db = mongogo.GetService("Crawler")
		}
		mongogo.MongogoInitial(db, "Livestreams", *mongoDBData)
		//log.Println("i is", i)
		//log.Printf("%+v", *mongoDBData)
	})
}
