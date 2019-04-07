package twitchstreams_crawler

/*
import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	//"log"
	"os"
	"strconv"
)

type DBConfig struct {
	ServerIP string
	DBPort   string
}

/*
type Streams struct {
	Id           string
	Title        string
	Descriptiong string
	Platform     string
	Video_Id     string
	Host         string
	Status       string
	Thumbnails   string
	Published    string
	Tags         string
	General_Tag  string
	Timestamp    string
	Language     string
	Viewcount    string
	Viewers      int
	Video_URL    string
}


const ConfigFilePath = "./config.json"

func MongogoInitial(streams Streams) {
	//var conf = getDBConfig()
	var db = getDB()
	insert(db, streams.Title, streams)
	insert(db, streams.Host, streams)
	insert(db, streams.Platform, streams)
	insert(db, strconv.Itoa(streams.Viewers), streams)
}

func insert(db *mgo.Database, collection string, streams Streams) {
	c := db.C(collection)
	err := c.Insert(streams)
	//log.Println("error")
	handleError(err, "")
}

func getDB() *mgo.Database {
	session, err := mgo.Dial("120.126.16.88:27017")
	handleError(err, "")

	session.SetMode(mgo.Monotonic, true)
	db := session.DB("Stream")
	return db
}

func getDBConfig() *DBConfig {
	file, _ := os.Open(ConfigFilePath)
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := new(DBConfig)
	err := decoder.Decode(&conf)
	handleError(err, "")
	return conf
}


func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}
*/
