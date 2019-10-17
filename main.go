package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/robfig/cron"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	cronJob := cron.New()
	cronJob.AddFunc("@hourly", ParseAndUpdateStockDB)
	cronJob.Start()

	templates := template.Must(template.ParseGlob("templates/*.htm"))

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Execute(w, nil)
	})

	http.HandleFunc("/stocks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			templates.Execute(w, nil)
		}
		symbol := r.FormValue("stockSearch")
		log.Printf("GET request on /stocks : %s", symbol)
		found, stock := GetStockDB(symbol)

		if found {
			templates.Execute(w, struct {
				Success bool
				Stock   StockData
			}{true, stock})
		}
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
