package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	go ParseAndUpdateStockDB()

	templates := template.Must(template.ParseGlob("templates/*.htm"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Execute(w, nil)
	})

	http.HandleFunc("/stocks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			templates.Execute(w, nil)
		}
		symbol := r.FormValue("stockSearch")
		log.Printf("GET request on /stocks : %s", symbol)
		stock := GetStockDB(symbol)

		templates.Execute(w, struct {
			Success bool
			Stock   StockData
		}{true, stock})
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
