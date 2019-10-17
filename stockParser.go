package main

/*
This file handles the job of genererating/updating the stock DB
*/
import (
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

//StockData is the per-stock type and holds stock info, in the database it is the value for the key(Symbol)
type StockData struct {
	Symbol    string  `json: "symbol" bson: "symbol"`
	Name      string  `json: "name" bson: "name"`
	Weight    float64 `json: "weight" bson: "weight"`
	Price     float64 `json: "price" bson: "price"`
	Change    float64 `json:"-" bson:"-"`
	ChangePct string  `json: "-" bson: "-"`
}

//Remove the , from prices
func normalizeAmerican(old string) string {
	return strings.Replace(old, ",", "", -1)
}

//ParseAndUpdateStockDB gets stock data from website and updates the databse
func ParseAndUpdateStockDB() {

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
	)

	c.OnHTML("table", func(e *colly.HTMLElement) {
		e.ForEach("table tr", func(x int, ele *colly.HTMLElement) {
			name := ele.ChildText("td:nth-of-type(2)")
			symbol := ele.ChildText("td:nth-of-type(3)")
			weight, _ := strconv.ParseFloat(normalizeAmerican(ele.ChildText("td:nth-of-type(4)")), 64)
			price, _ := strconv.ParseFloat(normalizeAmerican(ele.ChildText("td:nth-of-type(5)")), 64)
			change, _ := strconv.ParseFloat(normalizeAmerican(ele.ChildText("td:nth-of-type(6)")), 64)
			changePct := ele.ChildText("td:nth-of-type(7)")

			_ = UpdateStockToDB(StockData{strings.ToLower(symbol), name, weight, price, change, changePct})

		})
	})

	err := c.Visit("https://www.slickcharts.com/sp500")
	if err != nil {
		log.Fatal("Unable to visit site")
	}
}
