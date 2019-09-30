package main

type userStockData struct {
	Symbol   string
	Quantity int32
	Price    float64
}

//User per-user object that stores the list of stocks/free cash
type User struct {
	//ID   string   `json:"id" bson:"_id"`
	Auth AuthData `json: "auth" bson: "auth"`
	Cash float64  `json: "cash" bson: "cash"`
	//Tradekeys    []uint64 //array of pointers to trades done by user
}

//NewUser Creates new user and returns a User object
func NewUser(name string, plaintextPassword string) User {
	u := User{
		Auth: NewAuthData(name, plaintextPassword),
		Cash: 100000.00,
	}
	return u
}

//BuyStock buys stock with symbol for the user. returns false if user cannot buy stock
func (u *User) BuyStock(symbol string, qty int32) bool {
	stock := GetStockDB(symbol)
	if (stock.Price * float64(qty)) > u.Cash {
		return false
	}
	//u.Stocks = append(u.Stocks, userStockData{symbol, qty, stock.Price})
	u.Cash -= (stock.Price * float64(qty))
	return true

}
