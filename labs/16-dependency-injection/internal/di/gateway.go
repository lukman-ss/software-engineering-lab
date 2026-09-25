package di

import "fmt"

// Money is a value object instantiated directly. (Finding 5)
type Money struct {
	Amount   int
	Currency string
}

// PaymentGateway abstracts the concrete implementation. (Finding 1)
type PaymentGateway interface {
	Charge(m Money) error
}

// RealGateway is the production implementation simulating a network call.
type RealGateway struct{}

func (g *RealGateway) Charge(m Money) error {
	fmt.Printf("RealGateway charging %d %s\n", m.Amount, m.Currency)
	return nil
}
