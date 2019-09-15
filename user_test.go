package main

import (
	"testing"
)

func TestUser(t *testing.T) {
	u := NewUser("ram")
	if u == nil {
		t.Errorf("Unable to allocate user")
	}
	var db StockDB
	db.init()
	db.updateDB()

	success := u.BuyStock("AMZN", 50)
	if success != true {
		t.Errorf("Unable to buy stock")
	}
}
