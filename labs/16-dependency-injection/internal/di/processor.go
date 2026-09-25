package di

import "errors"

// Processor uses Constructor Injection (Finding 3)
type Processor struct {
	gateway PaymentGateway
}

// NewProcessor ensures valid state and immutability.
func NewProcessor(g PaymentGateway) *Processor {
	return &Processor{gateway: g}
}

func (p *Processor) ProcessPayment(amount int) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	m := Money{Amount: amount, Currency: "USD"} // Value object direct instantiation (Finding 5)
	return p.gateway.Charge(m)
}
