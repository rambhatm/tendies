package main

import (
	"testing"
)

func TestDB(t *testing.T) {
	stock := StockData{"AAPL", "apple", 55.0, 111.1, 11.1, "2"}
	err := InsertStockDB("AAPL", stock)
	if err != true {
		t.Errorf("Unable to insert into stock DB")
	}
	stock2 := GetStockDB("AAPL")
	if stock2.Name != "apple" {
		t.Errorf("Error with stock get from DB")
	}
}
