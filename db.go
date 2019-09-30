//This is the only file that directly deals with a database
//currently using  level DB bindings for go https://github.com/syndtr/goleveldb

package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	stockdbfile     = "db/stock.db"
	herokuDB        = "heroku_kmkcwhkx"
	registeredUsers = "users"
	tradedb         = "db/trade.db"
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

///MONGODB

//inserts key value pair into dbfile
func insertDB(dbfile string, key string, val bytes.Buffer) bool {
	db, err := leveldb.OpenFile(dbfile, nil)
	if err != nil {
		log.Fatal("DB open error")
		return false
	}
	defer db.Close()

	err = db.Put([]byte(key), val.Bytes(), nil)
	if err != nil {
		log.Fatal("DB open error")
		return false
	}
	return true
}

//InsertStockDB encodes and inserts stock into stock DB
func InsertStockDB(symbol string, stock StockData) bool {
	//Encode to gob,needed for structs
	var gobstock bytes.Buffer
	enc := gob.NewEncoder(&gobstock)
	_ = enc.Encode(stock)

	return insertDB(stockdbfile, symbol, gobstock)
}

//GetStockDB decodes and returns stock
func GetStockDB(symbol string) (stock StockData) {
	stockdb, err := leveldb.OpenFile(stockdbfile, nil)
	if err != nil {
		log.Fatal("Stockdb open error ", err)
		return
	}
	defer stockdb.Close()

	data, err := stockdb.Get([]byte(symbol), nil)
	//Decode the value from gob
	gobstock := bytes.NewBuffer(data)
	dec := gob.NewDecoder(gobstock)
	dec.Decode(&stock)
	return
}
