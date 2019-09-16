package main

/*
Routes for the server

/ - main app html
/stock/ - return list of symbols to search
/stock/{symbol} - return stockdata for the symbol

*/
import (
	"encoding/json"
	//"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getAllStockSymbols(resp http.ResponseWriter, req *http.Request) {
	log.Printf("endpoint: /stocks/")
}

func getStock(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	symbol := vars["symbol"]
	log.Printf("endpoint: /stock/%s", symbol)
	json.NewEncoder(resp).Encode(db.getStock(symbol))
}

func app(w http.ResponseWriter, r *http.Request) {

}

//InitHTTPServer Creates the HTTP server for the given port
func InitHTTPServer(port int) {
	log.Printf("Initing server on port %d", port)
	router := mux.NewRouter().StrictSlash(true)
	_ = router
	//var templates = template.Must(template.ParseFiles("app.html"))

	http.HandleFunc("/", app)

	//router.HandleFunc("/stocks/", getAllStockSymbols)
	//router.HandleFunc("/stock/{symbol}/", getStock)

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), router))
}
