package repository

import (
	"bytes"
	"crawlers/controller/resource"
	"crawlers/controller/watermelon"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

var (
	watermelonURL = "https://live.ixigua.com/category/1/109"
)

type seleniumWatermelonRepository struct{}

// NewSeleniumWatermelonRepository is new a seleniumWatermelonRepository struct
func NewSeleniumWatermelonRepository() watermelon.Repository {
	return &seleniumWatermelonRepository{}
}

func (watermelonRepo *seleniumWatermelonRepository) Init(config resource.Config) {
	log.Println("[Watermelon] START")
	defer log.Println("[Watermelon] END")

	port, gerFreePortErr := resource.GetFreePort()
	if gerFreePortErr != nil {
		resource.HandleError(gerFreePortErr, "[Watermelon] gerFreePortErr")
		return
	}

	service, getServiceErr := resource.GetService(port)
	if getServiceErr != nil {
		resource.HandleError(getServiceErr, "[Watermelon] getServiceErr")
		return
	}
	defer service.Stop()

	wd, getRemoteErr := resource.GetRemote(port)
	if getRemoteErr != nil {
		resource.HandleError(getRemoteErr, "[Watermelon] getRemoteErr")
		return
	}
	defer wd.Quit()

	frameHTML, getDataFromStreamListErr := watermelonRepo.getStreamListHTML(wd, watermelonURL)
	if getDataFromStreamListErr != nil {
		resource.HandleError(getDataFromStreamListErr, "[Watermelon] getDataFromStreamListErr")
		return
	}

	datas, parseStreamListErr := watermelonRepo.parseStreamListHTML(frameHTML)
	if parseStreamListErr != nil {
		resource.HandleError(parseStreamListErr, "[Watermelon] parseStreamListErr")
		return
	}

	for _, data := range datas {
		resource.SendToElasticByHTTPPost(data, config)
	}
	log.Println("Total streams are ", len(datas))
}

func (watermelonRepo *seleniumWatermelonRepository) getStreamListHTML(wd selenium.WebDriver, url string) (string, error) {
	if wd == nil {
		err := errors.New("[Watermelon] wd is nil")
		return "", err
	}

	wd.Get(url)
	time.Sleep(100 * time.Millisecond)

	elem, findElementErr := wd.FindElement(selenium.ByCSSSelector, "div.Feed-footer")
	if findElementErr != nil {
		resource.HandleError(findElementErr, "[Watermelon] findElementErr")
		return "", findElementErr
	}

	for count := 0; count < 300; count++ {
		elem.MoveTo(0, 0)
		time.Sleep(100 * time.Millisecond)
	}

	frameHTML, pageSourceErr := wd.PageSource()
	return frameHTML, pageSourceErr
}

func (watermelonRepo *seleniumWatermelonRepository) parseStreamListHTML(frameHTML string) (datas []resource.Data, err error) {
	doc, newDocumentFromReaderErr := goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
	if newDocumentFromReaderErr != nil {
		resource.HandleError(newDocumentFromReaderErr, "[Watermelon] newDocumentFromReaderErr")
		return nil, newDocumentFromReaderErr
	}
	doc.Find("div.feed-item__list").Find("div.feed-card").Each(func(i int, s *goquery.Selection) {
		title := watermelonRepo.getTitle(s)
		host := watermelonRepo.getHost(s)
		thumbnails := watermelonRepo.getThumbnails(s)
		localTime := resource.GetLocalTime()
		tags := watermelonRepo.getTags(s)
		viewers, getViewersErr := watermelonRepo.getViewers(s)
		if getViewersErr != nil {
			resource.HandleError(getViewersErr, "[Watermelon] getViewersErr")
			err = getViewersErr
			return
		}
		videoURL, isVideoURLExist := watermelonRepo.getVideoURL(s)
		if !isVideoURLExist {
			err = errors.New("VideoURL not Exist")
			return
		}
		data := resource.Data{
			Title:            title,
			Description:      "",
			Platform:         "西瓜直播",
			VideoID:          "",
			Host:             host,
			Status:           "live",
			Thumbnails:       thumbnails,
			Published:        localTime,
			Tags:             tags,
			GeneralTag:       "",
			Timestamp:        localTime,
			Language:         "zh",
			ViewCount:        viewers,
			Viewers:          viewers,
			VideoURL:         videoURL,
			VideoEmbedded:    "",
			ChatRoomEmbedded: "",
			Channel:          "",
		}
		datas = append(datas, data)
	})
	return
}

func (watermelonRepo *seleniumWatermelonRepository) getTitle(s *goquery.Selection) string {
	return s.Find("a.feed-card__title").Text()
}

func (watermelonRepo *seleniumWatermelonRepository) getHost(s *goquery.Selection) string {
	return s.Find("span.feed-card__user__name").Text()
}

func (watermelonRepo *seleniumWatermelonRepository) getThumbnails(s *goquery.Selection) string {
	thumbnails, _ := s.Find("img.feed-card__cover__opacity-img").Attr("src")
	thumbnails = strings.Replace("http:"+thumbnails, "580x327", "480x360", -1)
	return thumbnails
}

func (watermelonRepo *seleniumWatermelonRepository) getTags(s *goquery.Selection) string {
	return s.Find("div.feed-card__activity-tag").Text()
}

func (watermelonRepo *seleniumWatermelonRepository) getViewers(s *goquery.Selection) (int, error) {
	viewersStr := s.Find("span.feed-card__user__count").Text()
	if strings.Contains(viewersStr, "w") {
		viewersStr = strings.Replace(viewersStr, "w", "", -1)
		viewers, parseFloatErr := strconv.ParseFloat(viewersStr, 32)
		return int(viewers * 10000), parseFloatErr
	}
	viewers, parseFloatErr := strconv.ParseFloat(viewersStr, 32)
	return int(viewers), parseFloatErr
}

func (watermelonRepo *seleniumWatermelonRepository) getVideoURL(s *goquery.Selection) (string, bool) {
	url, isVideoURLExist := s.Find("a").Attr("href")
	return "https://live.ixigua.com" + url, isVideoURLExist
}

func (watermelonRepo *seleniumWatermelonRepository) ParseStreamURLHTML(wd selenium.WebDriver, videoURLs []string) (datas []resource.Data, err error) {
	for _, url := range videoURLs {
		frameHTML, getStreamURLHTMLErr := watermelonRepo.getStreamURLHTML(wd, url)
		if getStreamURLHTMLErr != nil {
			resource.HandleError(getStreamURLHTMLErr, "[Watermelon] getStreamURLHTMLErr")
			break
		}
		doc, newDocumentFromReaderErr := goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
		if newDocumentFromReaderErr != nil {
			resource.HandleError(newDocumentFromReaderErr, "[Watermelon] newDocumentFromReaderErr")
			break
		}

		title := watermelonRepo.getTitleFromHost(doc)
		host := watermelonRepo.getHostFromHost(doc)
		thumbnails := watermelonRepo.getThumbnailsFromHost(doc)
		tags := watermelonRepo.getTagsFromHost(doc)
		viewers := watermelonRepo.getViewersFromHost(doc)
		localTime := resource.GetLocalTime()

		data := resource.Data{
			Title:            title,
			Description:      "",
			Platform:         "西瓜直播",
			VideoID:          "",
			Host:             host,
			Status:           "live",
			Thumbnails:       thumbnails,
			Published:        localTime,
			Tags:             tags,
			GeneralTag:       "",
			Timestamp:        localTime,
			Language:         "zh",
			ViewCount:        viewers,
			Viewers:          viewers,
			VideoURL:         url,
			VideoEmbedded:    "",
			ChatRoomEmbedded: "",
			Channel:          "",
		}
		datas = append(datas, data)
	}
	return
}

func (watermelonRepo *seleniumWatermelonRepository) getStreamURLHTML(wd selenium.WebDriver, url string) (string, error) {
	wd.Get(url)
	time.Sleep(3000 * time.Millisecond)
	frameHTML, pageSourceErr := wd.PageSource()
	return frameHTML, pageSourceErr
}

func (watermelonRepo *seleniumWatermelonRepository) getTitleFromHost(s *goquery.Document) string {
	return ""
}

func (watermelonRepo *seleniumWatermelonRepository) getHostFromHost(s *goquery.Document) string {
	return ""
}

func (watermelonRepo *seleniumWatermelonRepository) getTagsFromHost(s *goquery.Document) string {
	return ""
}

func (watermelonRepo *seleniumWatermelonRepository) getThumbnailsFromHost(s *goquery.Document) string {
	return ""
}

func (watermelonRepo *seleniumWatermelonRepository) getViewersFromHost(s *goquery.Document) int {
	return 2
}
