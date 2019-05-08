# *_Twitch Crawler_*

This repository contains four packages to call the APIs.
You need Go version at least `1.12.x`

* [Set GOOS](#goos-setting)
* [Build the project](#build)
* [Update the game infromation Json file](#update-the-json-file)
* [Start the main program](#run-the-main-program)

## GOOS setting
----

If your os is `mac`, please try `$ export GOOS=darwin`
If your os is `linux`, please try `$ export GOOS=linux`

## Build
----
Make sure you download the whole file completely.

    $ go build TwitchCrawler.go

## Update the Json file
----
Please update the game information at first, which will spend you some time.
Please wait for it patiently. Then, it will create a Json file which contains about 10000
pieces of game information data.

    $ sudo ./TwitchCrawler Update

## Run the main program
----
    $ sudo ./TwitchCrawler Start

After runnuing this program, you will see the output.
It's like
```
Start process !
Convert Json file into this crawler.
Twitch Crawler is running ...
   Initial      1   422
   .
   .
   .
```
