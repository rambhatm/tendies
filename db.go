//This is the only file that directly deals with a database
//currently using https://github.com/syndtr/goleveldb

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

func InsertStockDB(symbol string, stock StockData) bool {
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
