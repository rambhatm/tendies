package main

type userStockData struct {
	Symbol   string
	Quantity int32
	Price    float64
}

//User per-user object that stores the list of stocks/free cash
type User struct {
	Username string
	FreeCash float64
	Stocks   []userStockData
}

//NewUser Creates new user and returns a User object
func NewUser(name string) *User {
	u := new(User)
	u.Username = name
	u.FreeCash = 100000.00
	return u
}

//BuyStock buys stock with symbol for the user. returns false if user cannot buy stock
func (u *User) BuyStock(symbol string, qty int32) bool {
	stock := db.getStock(symbol)
	if (stock.Price * float64(qty)) > u.FreeCash {
		return false
	}
	u.Stocks = append(u.Stocks, userStockData{symbol, qty, stock.Price})
	u.FreeCash -= (stock.Price * float64(qty))
	return true
}
