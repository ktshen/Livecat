package mongogo

import (
	"crawlers/controller/resource"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongogo ...
type Mongogo struct {
	Session *mgo.Session
	DB      *mgo.Database
}

// NotificationData ...
type NotificationData struct {
	Platform string `bson:"platform"`
	Pagelink string `bson:"pagelink"`
	Account  string `bson:"account"`
	Chatroom string `bson:"chatroom"`
	Email    string `bson:"email"`
}

// Init ...
func (mongogo *Mongogo) Init(serverIP string, port string, dBName string) {
	mongogo.getSession(serverIP, port)
	mongogo.Session.SetMode(mgo.Monotonic, true)
	mongogo.getDB(dBName)
}

func (mongogo *Mongogo) getSession(serverIP string, port string) error {
	var err error
	mongogo.Session, err = mgo.Dial(serverIP + ":" + port)
	return err
}

func (mongogo *Mongogo) getDB(dbName string) {
	mongogo.DB = mongogo.Session.DB(dbName)
}

// Insert ...
func (mongogo *Mongogo) Insert(collection string, data resource.Data) {
	c := mongogo.DB.C(collection)
	insertErr := c.Insert(data)
	if insertErr != nil {
		resource.HandleError(insertErr, "[Mongogo] InsertErr")
	}
}

// Find ...
func (mongogo *Mongogo) Find(collection string, key string, value string) []NotificationData {
	c := mongogo.DB.C(collection)
	datas := []NotificationData{}
	err := c.Find(bson.M{key: value}).All(&datas)
	if err != nil {
		log.Println(err)
		return nil
	}
	return datas
}

// Remove ...
func (mongogo *Mongogo) Remove(collection string, key string, value string) error {
	c := mongogo.DB.C(collection)
	removeErr := c.Remove(bson.M{key: value})
	return removeErr
}
