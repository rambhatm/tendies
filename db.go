//This is the only file that directly deals with a database
//currently using  level DB bindings for go https://github.com/syndtr/goleveldb

package main

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	stockdbfile = "db/stock.db"
	userdbfile  = "db/user.db"
	tradedb     = "db/trade.db"
)

//inserts key value pair into dbfile
func insertDB(dbfile string, key string, val bytes.Buffer) bool {
	db, err := leveldb.OpenFile(dbfile, nil)
	if err != nil {
		log.Fatal("%s DB open error", dbfile)
		return false
	}
	defer db.Close()

	err = db.Put([]byte(key), val.Bytes(), nil)
	if err != nil {
		log.Fatal("%s DB open error", dbfile)
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
		log.Fatal("Stockdb open error")
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

//InsertUserDB encodes and inserts user into the DB
func InsertUserDB(username string, u User) bool {
	//Encode to gob,needed for structs
	var gobuser bytes.Buffer
	enc := gob.NewEncoder(&gobuser)
	_ = enc.Encode(u)

	return insertDB(userdbfile, username, gobuser)
}
