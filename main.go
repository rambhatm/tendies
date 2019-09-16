package main

//"fmt"
import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	var db StockDB
	db.init()
	db.updateDB()
	//InitHTTPServer(4000)
	//fmt.Println(db.getStock("AMZN"))

	templates := template.Must(template.ParseFiles("templates/app.htm"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.Execute(w, nil)
	})

	http.HandleFunc("/stocks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			templates.Execute(w, nil)
		}
		symbol := r.FormValue("stockSearch")
		log.Printf("POST request on /stocks : %s", symbol)

		templates.Execute(w, struct {
			success bool
			message string
		}{true, "yay!!" + symbol})
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
