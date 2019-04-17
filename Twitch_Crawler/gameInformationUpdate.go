package main

import (
	"fmt"
	"twitchstreams_crawler"
)

func main() {

	fmt.Println("Initialize the game list. ...")
	twitchstreams_crawler.GetAllGameInformation()
	fmt.Println("Initialize the game list completed.")
}
