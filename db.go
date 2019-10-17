//This is the only file that directly deals with a database
//currently using  level DB bindings for go https://github.com/syndtr/goleveldb

package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (

	//Heroku mongo DB name
	herokuDB = "heroku_kmkcwhkx"

	//Collections
	registeredUsers = "users"
	stockInfo       = "stockInfo" //Stores current price and stock details
	stockHist       = "stockHist"
	trades          = "trades"
)

// Set client options
//TODO use env var
var clientOptions = options.Client().ApplyURI("")

func InsertUserToDB(u User) (err error) {
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("cannot connect to mongodb")
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(herokuDB).Collection(registeredUsers)

	insertResult, err := collection.InsertOne(context.TODO(), u)
	if err != nil {
		log.Printf("error: unable to insert to user DB %s", err)
		return
	}

	//_, err = collection.Indexes().CreateOne(
	//	context.Background(),
	//	mongo.IndexModel{
	//		Keys:    bsonx.Doc{{"auth.username", bsonx.Int32(1)}},
	//		Options: options.Index().SetUnique(true),
	//	},
	//)

	log.Printf("Inserted new user: %s id: %s ", u.Auth.Username, insertResult.InsertedID)
	return
}

func GetUserFromDB(username string) (found bool, u User) {
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("cannot connect to mongodb")
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(herokuDB).Collection(registeredUsers)

	filter := bson.D{{"username", username}}

	err = collection.FindOne(context.TODO(), filter).Decode(&u)
	if err != nil {
		log.Printf("Unable to find username: %s err: %s", username, err)
		found = false
		return
	}
	found = true
	return
}

//UpdateStockToDB updates the current stock value in stockinfo, and updates history if price is different
func UpdateStockToDB(s StockData) (err error) {
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("cannot connect to mongodb")
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(herokuDB).Collection(stockInfo)
	filter := bson.D{{"symbol", s.Symbol}}

	var curr StockData
	err = collection.FindOne(context.TODO(), filter).Decode(&curr)
	if err != nil {
		log.Printf("<DB error> finding stock: %s in stockinfo : %s", s.Symbol, err)
		return
	}
	if curr.Symbol == s.Symbol {
		if curr.Price != s.Price {
			//Update the price and set history
			filter := bson.D{{"name", "Ash"}}

			update := bson.D{
				{"$set", bson.D{
					{"price", s.Price},
				}},
			}

			updateResult, err = collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Printf("<DB error>Updating stock: %s in stockinfo : %s", s.Symbol, err)
				return
			}

			fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
		}
	}
}
