Twitch Crawler
==============

This repository contains four packages to call the APIs.
You need Go version at least `1.12.x`

GOOS
----

If your os is mac, please try `$ export GOOS=darwin`.
If your os is linux, please try `$ export GOOS=linux`.

# Get started
`$ go build TwitchCrawler.go`
## The first time

Please update the game information at first, which will spend you some time,
please wait for it patiently. Then, it will create a Json file which contains about 10000
pieces of game information data.
`./TwitchCrawler Update`

## Get TwitchAPI
`./TwitchCrawler Start`
