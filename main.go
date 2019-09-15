package main

//"fmt"

func main() {
	var db StockDB
	db.init()
	db.updateDB()
	InitHTTPServer(4000)
	//fmt.Println(db.getStock("AMZN"))
}
