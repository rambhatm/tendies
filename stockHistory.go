package main

import (
	"time"
)

type StockHistory struct {
	Timestamp time.Time `json:"timestamp" bson: "timestamp`
	Price     float64   `json:"price" bson:"price"`
}
