package main

import (
//"fmt"
)

func main() {
	var db StockDB
	db.init()
	db.updateDB()
	InitHTTPServer(&db, 4000)
	//fmt.Println(db.getStock("AMZN"))
}
