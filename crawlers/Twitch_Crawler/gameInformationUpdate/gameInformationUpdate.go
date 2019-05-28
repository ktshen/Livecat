package gameInformationUpdate

import (
	"Twitch_Crawler/twitchstreams_crawler"
	"fmt"
)

func Update() {
	fmt.Println("Initialize the game list. ...")
	twitchstreams_crawler.GetAllGameInformation()
	fmt.Println("Initialize the game list completed.")
}
