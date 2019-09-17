//This is the only file that directly deals with a database
//currently using  level DB https://github.com/syndtr/goleveldb

package main

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	STOCKDB = "db/stock.db"
	USERDB  = "db/user.db"
)

//InsertStockDB encodes and inserts stock into stock DB
func InsertStockDB(symbol string, stock StockData) bool {
	//Encode t0o gob
	var gobstock bytes.Buffer
	enc := gob.NewEncoder(&gobstock)
	_ = enc.Encode(stock)

	stockdb, err := leveldb.OpenFile(STOCKDB, nil)
	if err != nil {
		log.Fatal("Stockdb open error")
		return false
	}
	defer stockdb.Close()

	err = stockdb.Put([]byte(symbol), gobstock.Bytes(), nil)
	if err != nil {
		log.Fatal("StockDB insert error")
		return false
	}
	return true
}

//GetStockDB decodes and returns stock
func GetStockDB(symbol string) (stock StockData) {
	stockdb, err := leveldb.OpenFile(STOCKDB, nil)
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
