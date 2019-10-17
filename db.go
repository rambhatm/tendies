//This is the only file that directly deals with a database
//currently using  level DB bindings for go https://github.com/syndtr/goleveldb

package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

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
var mongodbURI = os.Getenv("MONGODB_URI")
var clientOptions = options.Client().ApplyURI(mongodbURI)

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

func insertHistoryDB(client *mongo.Client, s StockData) (err error) {
	//Set the price in history databse for the symbol collection
	histCollection := client.Database(herokuDB).Collection(s.Symbol)
	var sHist StockHistory
	sHist.Timestamp = time.Now()
	sHist.Price = s.Price

	_, err = histCollection.InsertOne(context.TODO(), sHist)
	if err != nil {
		log.Printf("<DB error> unable to add to historyDB", err)
		return
	}
	return
}

func getStock(client *mongo.Client, sym string) (found bool, s StockData) {
	infoCollection := client.Database(herokuDB).Collection(stockInfo)
	filter := bson.D{{"symbol", sym}}

	err := infoCollection.FindOne(context.TODO(), filter).Decode(&s)
	if err != nil {
		log.Printf("<DB error> finding stock: %s in stockinfo : %s", sym, err)
		return
	}
	if sym == s.Symbol {
		found = true
	}
	return
}

//UpdateStockToDB updates the current stock value in stockinfo, and updates history if price is different
func UpdateStockToDB(s StockData) (err error) {
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("cannot connect to mongodb")
	}
	defer client.Disconnect(context.TODO())

	found, curr := getStock(client, s.Symbol)

	if found {
		if curr.Price != s.Price {
			//Update the price and set history
			filter := bson.D{{"name", "Ash"}}

			update := bson.D{
				{"$set", bson.D{
					{"price", s.Price},
				}},
			}

			infoCollection := client.Database(herokuDB).Collection(stockInfo)
			_, err = infoCollection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				log.Printf("<DB error>Updating stock: %s in stockinfo : %s", s.Symbol, err)
				return
			}

			insertHistoryDB(client, s)

		}
		//noting else to do, we have the latest price
	} else {
		//insert into stockinfo for the first time
		infoCollection := client.Database(herokuDB).Collection(stockInfo)
		_, err = infoCollection.InsertOne(context.TODO(), s)
		if err != nil {
			log.Printf("<DB error>: unable to insert stock info %s", err)
			return
		}
		insertHistoryDB(client, s)
	}
	return
}

//GetStockDB gets the current price and stock details
func GetStockDB(symbol string) (stockFound bool, s StockData) {
	symbol = strings.ToLower(symbol)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("cannot connect to mongodb")
	}
	defer client.Disconnect(context.TODO())

	stockFound, s = getStock(client, symbol)
	return
}
