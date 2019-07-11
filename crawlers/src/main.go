package main

import (
	"crawlers/controller"
	seventeenRepo "crawlers/controller/seventeen/repository"
	watermelonRepo "crawlers/controller/watermelon/repository"
	youtubeRepo "crawlers/controller/youtube/repository"
	"flag"
	"time"
)

// Vars are the variables about arguments
var (
	configName = flag.String("configName", "config", "The configName represent config file name ")
	configType = flag.String("configType", "yaml", "The configType represent config file type ")
	configPath = flag.String("configPath", "./", "The configPath represent config file configPath ")
)

func main() {
	flag.Parse()

	controller := controller.Controller{
		YouTube:        youtubeRepo.NewSeleniumYouTubeRepository(),
		YouTubeFreq:    10 * time.Minute,
		Watermelon:     watermelonRepo.NewSeleniumWatermelonRepository(),
		WatermelonFreq: 10 * time.Minute,
		Seventeen:      seventeenRepo.NewSeleniumSeventeenRepository(),
		SeventeenFreq:  10 * time.Minute,
		ConfigName:     *configName,
		ConfigType:     *configType,
		ConfigPath:     *configPath,
	}
	// go controller.StreamerNotification()
	controller.Init()
}
