package mongogo

import (
	"encoding/json"
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

// DBConfig ...
type DBConfig struct {
	ServerIP string
	DBPort   string
}

// MongoDB ...
type MongoDB struct {
	Title            string
	Description      string
	Platform         string
	VideoId          string
	Host             string
	Status           string
	Thumbnails       string
	Published        string
	Tags             string
	GeneralTag       string
	Timestamp        string
	Language         string
	ViewCount        int
	Viewers          int
	VideoURL         string
	VideoEmbedded    string
	ChatRoomEmbedded string
	Channel          string
}

type Web struct {
	Host string
}

// ConfigFilePath ...
const ConfigFilePath = "./config.json"

// GetService ...
func GetService(dbName string) *mgo.Database {
	var conf = getDBConfig()
	var db = getDB(*conf, dbName)
	return db
}

// MongogoInitial ...
func MongogoInitial(db *mgo.Database, collection string, data MongoDB) {
	insert(db, collection, data)
}

// func MongogoInitial(mongoDB MongoDB, dbName string, collection string) {
// 	var conf = getDBConfig()
// 	var db = getDB(*conf, dbName)
// 	insert(db, collection, mongoDB)
// }

func MongogoWebInitial(web Web, collection string) {
	var conf = getDBConfig()
	var db = getDB(*conf, "Web")

	c := db.C(collection)
	err := c.Insert(web)
	handleError(err, "Web Insert Error")

}

func insert(db *mgo.Database, collection string, mongoDB MongoDB) {
	c := db.C(collection)
	err := c.Insert(mongoDB)
	handleError(err, "Insert Error")
}

func getDB(conf DBConfig, dbName string) *mgo.Database {
	session, err := mgo.Dial(conf.ServerIP + ":" + conf.DBPort)
	handleError(err, "Dial Error")

	session.SetMode(mgo.Monotonic, true)
	db := session.DB(dbName)
	return db
}

func getDBConfig() *DBConfig {
	file, _ := os.Open(ConfigFilePath)
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := new(DBConfig)
	err := decoder.Decode(&conf)
	handleError(err, "decode")
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
