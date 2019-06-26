package repository

import (
	"bytes"
	"crawlers/controller/resource"
	"crawlers/controller/seventeen"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/tebeka/selenium"
)

var (
	seventeenURL = "https://17.live/suggested"
)

type seleniumSeventeenRepository struct{}

// NewSeleniumSeventeenRepository is new a seleniumSeventeenRepository struct
func NewSeleniumSeventeenRepository() seventeen.Repository {
	return &seleniumSeventeenRepository{}
}

func (seventeenRepo *seleniumSeventeenRepository) Init(config resource.Config) {
	log.Println("[Seventeen] START")
	defer log.Println("[Seventeen] END")

	port, gerFreePortErr := resource.GetFreePort()
	if gerFreePortErr != nil {
		resource.HandleError(gerFreePortErr, "[Seventeen] gerFreePortErr")
		return
	}

	service, getServiceErr := resource.GetService(port)
	if getServiceErr != nil {
		resource.HandleError(getServiceErr, "[Seventeen] getServiceErr")
		return
	}
	defer service.Stop()

	wd, getRemoteErr := resource.GetRemote(port)
	if getRemoteErr != nil {
		resource.HandleError(getRemoteErr, "[Seventeen] getRemoteErr")
		return
	}
	defer wd.Quit()

	frameHTML, getDataFromStreamListErr := seventeenRepo.getStreamListHTML(wd, seventeenURL)
	if getDataFromStreamListErr != nil {
		resource.HandleError(getDataFromStreamListErr, "[Seventeen] getDataFromStreamListErr")
		return
	}

	datas, parseStreamListErr := seventeenRepo.parseStreamListHTML(frameHTML)
	if parseStreamListErr != nil {
		resource.HandleError(parseStreamListErr, "[Seventeen] parseStreamListErr")
		return
	}

	for _, data := range datas {
		resource.SendToElasticByHTTPPost(data, config)
	}
	log.Println("Total streams are ", len(datas))
}

func (seventeenRepo *seleniumSeventeenRepository) getStreamListHTML(wd selenium.WebDriver, url string) (string, error) {
	if wd == nil {
		err := errors.New("[Seventeen] wd is nil")
		return "", err
	}

	wd.Get(url)
	wd.MaximizeWindow("")
	time.Sleep(3000 * time.Millisecond)
	wd.ExecuteScript("window.scrollTo(0,5000)", nil)
	time.Sleep(15 * time.Second)

	frameHTML, pageSourceErr := wd.PageSource()
	return frameHTML, pageSourceErr
}

func (seventeenRepo *seleniumSeventeenRepository) parseStreamListHTML(frameHTML string) (datas []resource.Data, err error) {
	doc, newDocumentFromReaderErr := goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
	if newDocumentFromReaderErr != nil {
		resource.HandleError(newDocumentFromReaderErr, "[Seventeen] newDocumentFromReaderErr")
		return nil, newDocumentFromReaderErr
	}

	doc.Find("div.Grid__Row-eThVWD").Find("div.LiveStreamBlock__Block-bdJheI").Each(func(i int, s *goquery.Selection) {
		title := seventeenRepo.getTitle(s)
		thumbnails := seventeenRepo.getThumbnails(s)
		localTime := resource.GetLocalTime()
		viewers, getViewersErr := seventeenRepo.getViewers(s)
		if getViewersErr != nil {
			resource.HandleError(getViewersErr, "[Seventeen] getViewersErr")
			err = getViewersErr
			return
		}
		videoURL, isVideoURLExist := seventeenRepo.getVideoURL(s)
		if !isVideoURLExist {
			err = errors.New("VideoURL not Exist")
			return
		}
		data := resource.Data{
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

func (seventeenRepo *seleniumSeventeenRepository) getTitle(s *goquery.Selection) string {
	title, _ := s.Find("a").Attr("title")
	return title
}

func (seventeenRepo *seleniumSeventeenRepository) getThumbnails(s *goquery.Selection) string {
	thumbnails, _ := s.Find("a").Find("div.LiveStreamBlock__AvatarWrapper-CSDSj").Attr("style")
	thumbnails = thumbnails[strings.Index(thumbnails, "url(\"")+5:]
	thumbnails = thumbnails[:strings.Index(thumbnails, "url(\"")-4]
	return thumbnails
}

func (seventeenRepo *seleniumSeventeenRepository) getViewers(s *goquery.Selection) (int, error) {
	viewersStr := s.Find("a").Find("span.Msg-eAZFzz").Text()
	viewersStr = strings.Replace(viewersStr, ",", "", -1)
	viewers, atoiErr := strconv.Atoi(viewersStr)
	return viewers, atoiErr
}

func (seventeenRepo *seleniumSeventeenRepository) getVideoURL(s *goquery.Selection) (string, bool) {
	url, isVideoURLExist := s.Find("a").Attr("href")
	return "https://17.live" + url, isVideoURLExist
}

func (seventeenRepo *seleniumSeventeenRepository) ParseStreamURLHTML(wd selenium.WebDriver, videoURLs []string) (datas []resource.Data, err error) {
	for _, url := range videoURLs {
		frameHTML, getStreamURLHTMLErr := seventeenRepo.getStreamURLHTML(wd, url)
		if getStreamURLHTMLErr != nil {
			resource.HandleError(getStreamURLHTMLErr, "[Seventeen] getStreamURLHTMLErr")
			break
		}
		doc, newDocumentFromReaderErr := goquery.NewDocumentFromReader(bytes.NewReader([]byte(frameHTML)))
		if newDocumentFromReaderErr != nil {
			resource.HandleError(newDocumentFromReaderErr, "[Seventeen] newDocumentFromReaderErr")
			break
		}
		title := seventeenRepo.getTitleFromHost(doc)
		thumbnails := seventeenRepo.getThumbnailsFromHost(doc)
		viewers := seventeenRepo.getViewersFromHost(doc)
		localTime := resource.GetLocalTime()

		data := resource.Data{
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
			Language:         "zh",
			ViewCount:        viewers,
			Viewers:          viewers,
			VideoURL:         url,
			VideoEmbedded:    "",
			ChatRoomEmbedded: "",
			Channel:          "",
		}
		log.Printf("%+v", data)
		datas = append(datas, data)
	}

	return
}

func (seventeenRepo *seleniumSeventeenRepository) getStreamURLHTML(wd selenium.WebDriver, url string) (string, error) {
	wd.Get(url)
	time.Sleep(3000 * time.Millisecond)
	frameHTML, pageSourceErr := wd.PageSource()
	return frameHTML, pageSourceErr
}

func (seventeenRepo *seleniumSeventeenRepository) getTitleFromHost(s *goquery.Document) string {
	return s.Find("a.StreamerInfo__StreamerName-hcwkwO").Text()
}

func (seventeenRepo *seleniumSeventeenRepository) getThumbnailsFromHost(s *goquery.Document) string {
	thumbnails, _ := s.Find("div.StreamerInfo__StreamerInfoWrapper-bWkiCI").Find("div.withFadeIn__Com-dcgmuw").Attr("style")
	thumbnails = thumbnails[strings.Index(thumbnails, "url(\"")+5:]
	thumbnails = thumbnails[:strings.Index(thumbnails, "url(\"")-4]
	return thumbnails
}

func (seventeenRepo *seleniumSeventeenRepository) getViewersFromHost(s *goquery.Document) int {
	return 0
}
