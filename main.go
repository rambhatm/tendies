package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type stockData struct {
	Name      string
	Weight    float64
	Price     float64
	Change    float64
	ChangePct string
}

type stockDB struct {
	sp500 map[string]stockData
}

func normalizeAmerican(old string) string {
	return strings.Replace(old, ",", "", -1)
}
func (db *stockDB) init() {
	db.sp500 = make(map[string]stockData)
	return
}

func (db *stockDB) updateDB() {

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
			db.sp500[symbol] = stockData{name, weight, price, change, changePct}
			//fmt.Println(db.sp500[symbol])
		})
	})
	//fmt.Printf("Sp500 database extracted %d", len(db.sp500))

	err := c.Visit("https://www.slickcharts.com/sp500")
	if err != nil {
		fmt.Println(err)
	}
}

func (db stockDB) getStock(stk string) stockData {
	return db.sp500[stk]
}

func main() {
	var db stockDB
	db.init()
	db.updateDB()
	fmt.Println(db.getStock("AMZN"))
}
