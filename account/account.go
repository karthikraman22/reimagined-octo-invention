package account

import "context"

// Account is ...
type Account struct {
	Identifier    string
	Number        string
	Balance       float64
	BalanceUpdate chan<- balanceUpdate
}

// Balance update is modeled as a request
type balanceUpdate struct {
	amount   float64
	response chan<- float64
	ctx      context.Context
}

// Credit does ...
func (a *Account) Credit(amount float64, ctx context.Context) <-chan float64 {
	res := make(chan float64, 1)
	go func(credit balanceUpdate) {
		a.BalanceUpdate <- credit
	}(balanceUpdate{amount: amount, response: res, ctx: ctx})
	return res
}
