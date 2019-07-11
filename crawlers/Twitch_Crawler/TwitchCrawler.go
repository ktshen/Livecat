package main

import (
	"fmt"
	"os"

	//"mongogogo"
	//"douyustreams_crawler"
	"Twitch_Crawler/gameInformationUpdate"
	"Twitch_Crawler/twitchstreams_crawler"
	//"time"
	//"mongogo"
)

func main() {

	argsWithProg := os.Args
	//argsWithoutProg := os.Args[1:] //You can get individual args with normal indexing.
	//arg := os.Args[3]
	// fmt.Println(argsWithoutProg)
	// fmt.Println(arg)
	length := len(os.Args)
	fmt.Println(argsWithProg)
	if length < 2 {
		fmt.Printf("Your insert is ")
		fmt.Printf("%s\n", argsWithProg)
		fmt.Println("You must insert at least one argumant !")
		fmt.Println("For example : ./Twitchcrawler [Update/Start]")
		//fmt.Println("Update : Update the GameInformation Json File.")
		//fmt.Println("Start  : Start getting Twitch API.")
		return
	}

	fmt.Println("Start process ! ")

	switch os.Args[1] {
	case "Start":
		for {
			twitchstreams_crawler.TwitchCrawlerInitial()
		}
	case "Update":
		gameInformationUpdate.Update()
	}
}
