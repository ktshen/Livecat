package mongogogo

import (
	"context"
	"crawler_livecat/mongodbStruct"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	//"mongogo"
)

//const ConfigFilePath = "./config.json"
const (
	//MongoDBName    = "Crawler"
	MongoServerURI = "mongodb://120.126.16.88:27017"
)

var (
	DB *mongo.Database
)

// Insert a Bson data to MongoDB. This function does three things.
// First, connect with MongoDB server.
// Second, find if this data has existed already and then, inserts it or updates it.
// Finally, cuts the connection between local and server.

func MongoInsertOne(client *mongo.Client, DBName string, CollectionName string, mongoDB mongodbStruct.MongoDB) {
	var BsonMongoDB = JsonToBson(mongoDB)
	collection := MongoCollection(DBName, CollectionName, client)
	result := MongoFindOne(collection, BsonMongoDB)
	if result != nil {
		log.Println("Find the older one and update it .")
		err := MongoUpdateOne(collection, BsonMongoDB)
		handleError(err, "")
	} else {
		_, err := collection.InsertOne(context.Background(), BsonMongoDB)
		handleError(err, "")
	}
}

func MongoConnect() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoServerURI))
	handleError(err, "")
	ctx, _ := context.WithCancel(context.Background())
	err = client.Connect(ctx)
	handleError(err, "")
	err = client.Ping(context.TODO(), nil)
	handleError(err, "")
	log.Println("Connection to MongoDB.")
	return client
}

func MongoDisconnect(client *mongo.Client) error {
	err := client.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	log.Println("Connection to MongoDB closed.")
	return nil
}

func MongoCollection(DBName string, CollectionName string, client *mongo.Client) *mongo.Collection {
	collection := client.Database(DBName).Collection(CollectionName)
	//log.Println("COLLECTIONS", collection)
	return collection
}

// Find if there is a data has existed already.
// If has, return Reault, otherwise, return error.

func MongoFindOne(collection *mongo.Collection, Name mongodbStruct.MongoDB) *mongodbStruct.MongoDB {
	var Result *mongodbStruct.MongoDB
	filter := bson.D{{"host", Name.Host}}
	err := collection.FindOne(context.TODO(), filter).Decode(&Result)
	if err != nil {
		log.Println("Insert a new one")
		return nil
	}
	return Result
}

// If there is the older one in document, update it
// Set filter that you want to find the key out.
// Set update by using bson.D{,} encoding. The first parameter must start with $.
// 		$currentDate	Sets the value of a field to current date, either as a Date or a Timestamp.
// 		$inc	Increments the value of the field by the specified amount.
// 		$min	Only updates the field if the specified value is less than the existing field value.
// 		$max	Only updates the field if the specified value is greater than the existing field value.
// 		$mul	Multiplies the value of the field by the specified amount.
// 		$rename	Renames a field.
// 		$set	Sets the value of a field in a document.
// 		$setOnInsert	Sets the value of a field if an update results in an insert of a document. Has no effect on update operations that modify existing documents.
// 		$unset	Removes the specified field from a document.
func MongoUpdateOne(collection *mongo.Collection, Name mongodbStruct.MongoDB) error {
	filter := bson.D{{"host", Name.Host}}
	update := bson.D{{"$set", Name}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// Decodine data from Json to Bson
// In this package, the initial input of MongoInsertOne is Json Type, so I need to decode that to Bson Type
// Return a Bson typr MongoDB structure.
func JsonToBson(mongoDB mongodbStruct.MongoDB) mongodbStruct.MongoDB {
	var Decode mongodbStruct.MongoDB
	streams_encode, _ := bson.Marshal(mongoDB)
	bson.Unmarshal(streams_encode, &Decode)
	//Decode = append(Decode, bson.E{"expireAfterSeconds", 10})
	return Decode
}
