package crawler

import (
	"bytes"
	"log"

	//"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
	"gopkg.in/mgo.v2"

	"crawler_livecat/mongogo"
)

// CrawlerHTMl is
// type CrawlerHTML struct {
// 	Title            string
// 	Description      string
// 	Platform         string
// 	VideoId          string
// 	Host             string
// 	Status           string
// 	Thumbnails       string
// 	Published        string
// 	Tags             string
// 	GeneralTag       string
// 	Timestamp        string
// 	Language         string
// 	ViewCount        string
// 	Viewers          string
// 	VideoURL         string
// 	VideoEmbedded    string
// 	ChatRoomEmbedded string
// 	Channel          string
// }

// PikoInitial is intial to crawler Piko
func PikoInitial() {
	log.Println("START [Piko]")
	crawlerPiko("https://tw.pikolive.com/list/youtube/page/1.html")
	log.Println("END [Piko]")

}

func crawlerPiko(url string) {
	service := getService(portPiko)
	defer service.Stop()
	wd := getRemote(portPiko)
	defer wd.Quit()

	doc := getDoc(wd, url)
	pageList := getAllPage(doc)
	for _, url := range pageList {
		doc = getDoc(wd, url)
		getStreamerData(wd, doc)
	}

}

func getDoc(wd selenium.WebDriver, url string) (doc *goquery.Document) {

	wd.Get(url)
	time.Sleep(1 * time.Second)

	var frameHTML string
	frameHTML, err := wd.PageSource()
	handleError(err, "")

	doc, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
	handleError(err, "")
	return
}

func getAllPage(doc *goquery.Document) (pageList []string) {
	doc.Find("div.channel-list nav.paginate").Find("li").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			page, isPageExist := s.Find("a").Attr("href")
			if isPageExist {
				pageList = append(pageList, "https://tw.pikolive.com"+page)
			}
		}
	})
	return pageList[:len(pageList)-1]
}

func getStreamerData(wd selenium.WebDriver, doc *goquery.Document) {
	var db *mgo.Database
	doc.Find("div.channel-list").Find("div.channel-item-wrapper").Each(func(i int, s *goquery.Selection) {

		title := s.Find("div.item-title").Text()
		platform, isPlatformExist := s.Find("div.item-info-left").Attr("title")
		if isPlatformExist {
			if strings.Contains(platform, "twitch") {
				platform = "Twitch"
			} else if strings.Contains(platform, "youtube") {
				platform = "YouTube"
			} else {
				platform = ""
			}
		}
		username := s.Find("span.username").Text()
		channelname := s.Find("div.item-info-real-name").Text()
		thumbnails, _ := s.Find("div.item-image").Find("img").Attr("src")
		thumbnails = strings.Replace(thumbnails, "320x180", "480x360", -1)
		videourl, _ := s.Find("div.channel-item").Find("a").Attr("href")

		userDoc := getDoc(wd, "https://tw.pikolive.com"+videourl)
		publishedInPiko := userDoc.Find("div#live-channel").Find("div#frame-footer").Find("span#startedTime").Text()
		localTime := getLocalTime()
		published := getPublishedTime(localTime, publishedInPiko)

		viewers, _ := strconv.Atoi(userDoc.Find("span#viewers").Text())
		videoembedded, _ := userDoc.Find("div#live-channel").Find("div#frame-view").Find("iframe").Attr("src")
		chatroomembedded, _ := userDoc.Find("div#chat-area").Find("div.chat-room").Find("iframe").Attr("src")
		if platform == "YouTube" {
			channelname = username
			chatroomembedded = chatroomembedded[:strings.Index(chatroomembedded, "&embed_domain")] + "&embed_domain=dlivecat.com"
		}
		mongoDBData := &mongogo.MongoDB{
			Title:            title,
			Description:      "",
			Platform:         platform,
			VideoID:          "",
			Host:             username,
			Status:           "live",
			Thumbnails:       thumbnails,
			Published:        published.String(),
			Tags:             "",
			GeneralTag:       "",
			Timestamp:        localTime.String(),
			Language:         "",
			ViewCount:        viewers,
			Viewers:          viewers,
			VideoURL:         "https://tw.pikolive.com" + videourl,
			VideoEmbedded:    videoembedded,
			ChatRoomEmbedded: chatroomembedded,
			Channel:          channelname,
		}
		if i%100 == 0 {
			db = mongogo.GetService("Crawler")
		}
		mongogo.MongogoInitial(db, "Livestreams", *mongoDBData)
		// log.Println("i is", i)
		// log.Printf("%+v", *mongoDBData)
	})
}

func getPublishedTime(t time.Time, published string) time.Time {
	tmp := strings.Split(published, ":")

	if len(tmp) != 3 {
		return t
	}

	hr, _ := strconv.Atoi(tmp[0])
	min, _ := strconv.Atoi(tmp[1])
	sec, _ := strconv.Atoi(tmp[2])
	t = t.Add(time.Hour * time.Duration(hr) * -1)
	t = t.Add(time.Minute * time.Duration(min) * -1)
	t = t.Add(time.Second * time.Duration(sec) * -1)
	return t
}
