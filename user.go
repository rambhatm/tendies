package main

type userStockData struct {
	Symbol   string
	Quantity int32
	Price    float64
}

//User per-user object that stores the list of stocks/free cash
type User struct {
	Auth AuthData
	Cash float64
	//Tradekeys    []uint64 //array of pointers to trades done by user
}

//NewUser Creates new user and returns a User object
func NewUser(name string, passHash string) User {
	u := User{
		Auth: AuthData{
			Username:     name,
			PasswordHash: passHash,
		},
		Cash: 100000.00,
	}
	if ok := InsertUserDB(name, u); !ok {
		//TODO XXX
		return u
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
