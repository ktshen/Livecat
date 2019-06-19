package repository

import (
	"bytes"
	"crawlers/controller/resource"
	"crawlers/controller/youtube"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

type seleniumYouTubeRepository struct{}

// NewSeleniumYouTubeRepository is new a seleniumYouTubeRepository struct
func NewSeleniumYouTubeRepository() youtube.Repository {
	return &seleniumYouTubeRepository{}
}

func (youtubeRepo *seleniumYouTubeRepository) Init(config resource.Config, youtubeURL string) {
	log.Println("[YouTube] START")
	defer log.Println("[YouTube] END")

	port, gerFreePortErr := resource.GetFreePort()
	if gerFreePortErr != nil {
		resource.HandleError(gerFreePortErr, "[YouTube] gerFreePortErr")
		return
	}

	service, getServiceErr := resource.GetService(port)
	if getServiceErr != nil {
		resource.HandleError(getServiceErr, "[YouTube] getServiceErr")
		return
	}
	defer service.Stop()

	wd, getRemoteErr := resource.GetRemote(port)
	if getRemoteErr != nil {
		resource.HandleError(getRemoteErr, "[YouTube] getRemoteErr")
		return
	}
	defer wd.Quit()

	frameHTML, getDataFromStreamListErr := youtubeRepo.getStreamListHTML(wd, youtubeURL)
	if getDataFromStreamListErr != nil {
		resource.HandleError(getDataFromStreamListErr, "[YouTube] getDataFromStreamListErr")
		return
	}
	videoURLs, parseStreamListHTMLErr := youtubeRepo.parseStreamListHTML(frameHTML)
	if parseStreamListHTMLErr != nil {
		resource.HandleError(parseStreamListHTMLErr, "[YouTube] parseStreamListHTMLErr")
		return
	}
	datas, parseStreamURLHTMLErr := youtubeRepo.ParseStreamURLHTML(wd, videoURLs)
	if parseStreamURLHTMLErr != nil {
		resource.HandleError(parseStreamURLHTMLErr, "[YouTube] parseStreamURLHTMLErr")
		return
	}

	for _, data := range datas {
		resource.SendToElasticByHTTPPost(data, config)
	}
	log.Println("Total streams are ", len(datas))
}

func (youtubeRepo *seleniumYouTubeRepository) getStreamListHTML(wd selenium.WebDriver, url string) (string, error) {
	if wd == nil {
		return "", errors.New("[YouTube] wd is nil")
	}

	wd.Get(url)
	time.Sleep(1000 * time.Millisecond)

	elem, findElementErr := wd.FindElement(selenium.ByCSSSelector, "div#continuations")
	if findElementErr != nil {
		return "", findElementErr
	}
	for countt := 0; countt < 10; countt++ {
		elem.MoveTo(0, 0)
		time.Sleep(1 * time.Second)
	}

	frameHTML, pageSourceErr := wd.PageSource()
	return frameHTML, pageSourceErr
}

func (youtubeRepo *seleniumYouTubeRepository) parseStreamListHTML(frameHTML string) (videoURLs []string, err error) {
	doc, newDocumentFromReaderErr := goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
	if newDocumentFromReaderErr != nil {
		resource.HandleError(newDocumentFromReaderErr, "[YouTube] newDocumentFromReaderErr")
		return nil, newDocumentFromReaderErr
	}

	doc.Find("ytd-item-section-renderer.style-scope").Find("ytd-playlist-video-list-renderer.style-scope").Find("ytd-playlist-video-renderer.style-scope").Each(func(i int, s *goquery.Selection) {
		videoURL := youtubeRepo.getYouTubeVideoURL(s)
		videoURLs = append(videoURLs, videoURL)
	})

	return
}

func (youtubeRepo *seleniumYouTubeRepository) getYouTubeVideoURL(s *goquery.Selection) (videoURL string) {
	videoURL, _ = s.Find("div#content").Find("a").Attr("href")
	end := strings.Index(videoURL, "&list")
	if end != -1 {
		videoURL = videoURL[:end]
	}
	return "https://www.youtube.com" + videoURL
}

func (youtubeRepo *seleniumYouTubeRepository) ParseStreamURLHTML(wd selenium.WebDriver, videoURLs []string) (datas []resource.Data, err error) {
	for idx, url := range videoURLs {
		frameHTML, getStreamURLHTMLErr := youtubeRepo.getStreamURLHTML(wd, url)
		if getStreamURLHTMLErr != nil {
			resource.HandleError(getStreamURLHTMLErr, "[YouTube] getStreamURLHTMLErr")
			break
		}
		doc, newDocumentFromReaderErr := goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
		if newDocumentFromReaderErr != nil {
			resource.HandleError(newDocumentFromReaderErr, "[YouTube] newDocumentFromReaderErr")
			break
		}

		title := youtubeRepo.getYouTubeTitle(doc)
		host := youtubeRepo.getYouTubeHost(doc)
		videoID, videoEmbedded, chatEmbedded, thumbnails := getYouTubeEmbedded(url)
		localTime := resource.GetLocalTime()
		tags := youtubeRepo.getYouTubeTag(doc)
		viewers := youtubeRepo.getYouTubeViewers(doc)
		data := resource.Data{
			Title:            title,
			Description:      "",
			Platform:         "YouTube",
			VideoID:          videoID,
			Host:             host,
			Status:           "live",
			Thumbnails:       thumbnails,
			Published:        localTime,
			Tags:             tags,
			GeneralTag:       "",
			Timestamp:        localTime,
			Language:         "",
			ViewCount:        viewers,
			Viewers:          viewers,
			VideoURL:         url,
			VideoEmbedded:    videoEmbedded,
			ChatRoomEmbedded: chatEmbedded,
			Channel:          videoID,
		}
		datas = append(datas, data)
		if idx%100 == 0 {
			log.Println(idx)
		}

	}

	return
}

func (youtubeRepo *seleniumYouTubeRepository) getStreamURLHTML(wd selenium.WebDriver, url string) (string, error) {
	wd.Get(url)
	time.Sleep(3000 * time.Millisecond)
	frameHTML, pageSourceErr := wd.PageSource()
	return frameHTML, pageSourceErr
}

func (youtubeRepo *seleniumYouTubeRepository) getYouTubeTitle(s *goquery.Document) (title string) {
	return s.Find("ytd-page-manager#page-manager").Find("div#columns").Find("div#primary").Find("div#primary-inner").Find("div#info").Find("div#info-contents").Find("ytd-video-primary-info-renderer.style-scope").Find("h1.ytd-video-primary-info-renderer").Text()
}

func (youtubeRepo *seleniumYouTubeRepository) getYouTubeHost(s *goquery.Document) (title string) {
	return s.Find("ytd-page-manager#page-manager").Find("div#columns").Find("div#primary").Find("div#primary-inner").Find("div#meta").Find("div#meta-contents").Find("div#top-row").Find("div#owner-container").Find("a.yt-simple-endpoint").Text()
}

func (youtubeRepo *seleniumYouTubeRepository) getYouTubeViewers(s *goquery.Document) (viewers int) {
	str := s.Find("ytd-page-manager#page-manager").Find("div#columns").Find("div#primary").Find("div#primary-inner").Find("div#info").Find("div#info-contents").Find("ytd-video-primary-info-renderer.style-scope").Find("div#info").Find("yt-view-count-renderer").Find("span.view-count").Text()

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
	return
}

func (youtubeRepo *seleniumYouTubeRepository) getYouTubeTag(s *goquery.Document) (tags string) {
	var taglist []string
	s.Find("ytd-page-manager#page-manager").Find("div#columns").Find("div#primary").Find("div#primary-inner").Find("div#info").Find("div#info-contents").Find("ytd-video-primary-info-renderer.style-scope").Find("yt-formatted-string.super-title").Find("a").Each(func(i int, s *goquery.Selection) {
		tag := strings.Replace(s.Text(), "#", "", -1)
		taglist = append(taglist, tag)
	})
	tags = strings.Join(taglist, ", ")
	return
}

func getYouTubeEmbedded(url string) (videoID string, videoEmbedded string, chatEmbedded string, thumbnails string) {
	videoID = url[strings.Index(url, "/watch?v=")+9:]
	videoEmbedded = "https://www.youtube.com/embed/" + videoID + "?rel=0&autoplay=1"
	chatEmbedded = "https://www.youtube.com/live_chat?v=" + videoID + "&is_popout=1&embed_domain=www.ilivenet.com"
	thumbnails = "https://i.ytimg.com/vi/" + videoID + "/hqdefault_live.jpg?sqp=CJy9p-gF-oaymwEZCNACELwBSFXyq4qpAwsIARUAAIhCGAFwAQ==&rs=AOn4CLAvPzycs4nNURfIDs9_ptV5REh8nw"
	return
}
