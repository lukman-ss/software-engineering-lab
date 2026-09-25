package di

import "errors"

// Container acts as a Service Locator.
type Container interface {
	GetPaymentGateway() PaymentGateway
}

// BadProcessor injects the container directly (Anti-Pattern) (Finding 4)
type BadProcessor struct {
	container Container
}

func NewBadProcessor(c Container) *BadProcessor {
	return &BadProcessor{container: c}
}

func (p *BadProcessor) ProcessPayment(amount int) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	m := Money{Amount: amount, Currency: "USD"}
	return p.container.GetPaymentGateway().Charge(m)
}
