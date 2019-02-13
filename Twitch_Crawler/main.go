package main

import (
	"fmt"
	"streams_crawler"
	"time"
)

func main() {
	START := time.Now()
	Time := START.Format("Mon Jan _2 15:04:05 2006")
	fmt.Println("[START]  Time is ", Time)
	streams_crawler.TwitchAPI()
	fmt.Println("Your motion is compeleted")
}
