package controller

import (
	"crawlers/controller/mongogo"
	"crawlers/controller/resource"
	"crawlers/controller/seventeen"
	"crawlers/controller/watermelon"
	"crawlers/controller/youtube"
	"log"
	"time"

	"github.com/spf13/viper"
)

// Controller represent a struct to controll crawler interface
type Controller struct {
	YouTube        youtube.Repository
	YouTubeFreq    time.Duration
	Watermelon     watermelon.Repository
	WatermelonFreq time.Duration
	Seventeen      seventeen.Repository
	SeventeenFreq  time.Duration
	ConfigName     string
	ConfigType     string
	ConfigPath     string
}

var (
	youtubeURLs = []string{
		"https://www.youtube.com/playlist?list=PLU12uITxBEPFJGVb2zSgCaWvMBe7vHonB", // 直播
		"https://www.youtube.com/playlist?list=PLiCvVJzBupKlQ50jZqLas7SAztTMEYv1f", // 遊戲
		"https://www.youtube.com/playlist?list=PL8fVUTBmJhHJrxHg_uNTMyRmsWbFltuQV", // 運動
		"https://www.youtube.com/playlist?list=PL57quI9usf_sQMlMeQrUr5O_pCncZilx3", // 科技
		"https://www.youtube.com/playlist?list=PLIFqWCuxNyoiKKthaTBqjIH6m6A9INomt", // 動物
		"https://www.youtube.com/playlist?list=PLU12uITxBEPFnoOrc_w0oJL6CEeKRhpcb", // 行動應用程式
		"https://www.youtube.com/playlist?list=PLU12uITxBEPHvBRHoUt-fzKr8Iz1HpfUC", // 網路攝影機直播
	}
)

// Init ...
func (controller *Controller) Init() {
	config, _ := controller.ReadConfigFile(controller.ConfigName, controller.ConfigType, controller.ConfigPath)
	log.Println("config is", config)
	youtubeTicker := time.NewTicker(controller.YouTubeFreq)
	watermelonTicker := time.NewTicker(controller.WatermelonFreq)
	seventeenTicker := time.NewTicker(controller.SeventeenFreq)

	for _, youtubeURL := range youtubeURLs {
		go controller.YouTube.Init(config, youtubeURL)
	}
	// go controller.Watermelon.Init(config)
	// go controller.Seventeen.Init(config)

	for {
		select {
		case <-youtubeTicker.C:
			for _, youtubeURL := range youtubeURLs {
				go controller.YouTube.Init(config, youtubeURL)
			}
		case <-watermelonTicker.C:
			// go controller.Watermelon.Init(config)
		case <-seventeenTicker.C:
			// go controller.Seventeen.Init(config)
		}
	}
}

// ReadConfigFile is to read config file
func (controller *Controller) ReadConfigFile(configName string, configType string, configPath string) (resource.Config, error) {
	viper, readConfigFileError := controller.SetViper(configName, configType, configPath)
	if readConfigFileError != nil {
		log.Fatalf("Error read config file, %s", readConfigFileError)
	}
	var configurations resource.Config
	unmarshalErr := viper.Unmarshal(&configurations)
	return configurations, unmarshalErr
}

// SetViper is setting viper about config name, type, path
func (controller *Controller) SetViper(configName string, configType string, path string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	return v, err
}

// StreamerNotification ...
func (controller *Controller) StreamerNotification() {
	config, _ := controller.ReadConfigFile(controller.ConfigName, controller.ConfigType, controller.ConfigPath)
	mongogo := mongogo.Mongogo{}
	mongogo.Init(config.ElasticServerIP, config.DBPort, "Host")
	// err := mongogo.Remove("testCollection", "title", "test")
	datas := mongogo.Find("Host", "platform", "17直播")
	var urlS []string
	for _, data := range datas {
		urlS = append(urlS, data.Pagelink)
	}
	port, gerFreePortErr := resource.GetFreePort()
	if gerFreePortErr != nil {
		resource.HandleError(gerFreePortErr, "[Notification] gerFreePortErr")
		return
	}
	service, getServiceErr := resource.GetService(port)
	if getServiceErr != nil {
		resource.HandleError(getServiceErr, "[Notification] getServiceErr")
		return
	}
	defer service.Stop()
	wd, getRemoteErr := resource.GetRemote(port)
	if getRemoteErr != nil {
		resource.HandleError(getRemoteErr, "[Notification] getRemoteErr")
		return
	}
	defer wd.Quit()
	controller.Seventeen.ParseStreamURLHTML(wd, urlS)
}
