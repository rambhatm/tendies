package main

import (
	"testing"
)

func TestUser(t *testing.T) {
	u := NewUser("ram")
	if u == nil {
		t.Errorf("Unable to allocate user")
	}
	stock := StockData{"AAPL", "apple", 55.0, 111.1, 11.1, "2"}
	err := InsertStockDB("AAPL", stock)
	if err != true {
		t.Errorf("Unable to insert into stock DB")
	}
	bought := u.BuyStock("AAPL", 10)
	if bought == false {
		t.Errorf("Could not buy stock for user")
	}
}
