package mongogogo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"log"
	"mongodbStruct"
	"time"
)

const (
	MongoServerURI = "mongodb://120.126.16.88:27017"
)

var (
	DB *mongo.Database
)

// Insert a Bson data to MongoDB. This function does three things.
// First, connect to MongoDB server.
// Second, find if this data has existed already, and then, inserts it or updates it.
// Finally, cut the connection between local and server.

func MongoInsertOne(client *mongo.Client, DBName string, CollectionName string, mongoDB mongodbStruct.MongoDB) {
	var BsonMongoDB = JsonToBson(mongoDB)
	collection := MongoCollection(DBName, CollectionName, client)
	result := MongoFindOne(collection, BsonMongoDB)
	if result != nil {
		//log.Println("Find the older one and update it .")
		err := MongoUpdateOne(collection, BsonMongoDB)
		handleError(err, "MongoDB Update One")
	} else {
		_, err := collection.InsertOne(context.Background(), BsonMongoDB)
		handleError(err, "MongoDB Insert one")
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
	return collection
}

// Find if there is a data has existed already.
// If has, return Reault, otherwise, return nil.

func MongoFindOne(collection *mongo.Collection, Name mongodbStruct.MongoDB) *mongodbStruct.MongoDB {
	var Result *mongodbStruct.MongoDB
	filter := bson.D{{"host", Name.Host}}
	err := collection.FindOne(context.TODO(), filter).Decode(&Result)
	if err != nil {
		//log.Println("Insert a new one")
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

func PopulateIndex(database string, collection string, client *mongo.Client, ttl int32) {
	c := client.Database(database).Collection(collection)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	index := yieldIndexModel(ttl)
	c.Indexes().CreateOne(context.Background(), index, opts)
	log.Println("Successfully create the index")
}

func yieldIndexModel(ttl int32) mongo.IndexModel {
	var TTL *int32
	TTL = new(int32)
	*TTL = int32(ttl)

	keys := bsonx.Doc{{Key: "createdAt", Value: bsonx.Int32(int32(1))}}
	index := mongo.IndexModel{}
	index.Keys = keys
	var Options *options.IndexOptions
	Options = new(options.IndexOptions)
	Options.ExpireAfterSeconds = TTL
	index.Options = Options
	return index
}

// Encode data from Json to Bson
// In this package, the initial input of MongoInsertOne is Json Type, so we need to decode that to Bson Type
// Return a Bson MongoDB structure type.
func JsonToBson(mongoDB mongodbStruct.MongoDB) mongodbStruct.MongoDB {
	var Decode mongodbStruct.MongoDB
	streams_encode, _ := bson.Marshal(mongoDB)
	bson.Unmarshal(streams_encode, &Decode)
	return Decode
}
