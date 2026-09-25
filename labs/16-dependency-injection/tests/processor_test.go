package tests

import (
	"errors"
	"testing"

	"lab16/internal/di"
)

type MockGateway struct {
	ChargedMoney di.Money
	ShouldFail   bool
}

func (m *MockGateway) Charge(money di.Money) error {
	if m.ShouldFail {
		return errors.New("gateway unavailable")
	}
	m.ChargedMoney = money
	return nil
}

func TestProcessor_Success(t *testing.T) {
	mock := &MockGateway{}
	proc := di.NewProcessor(mock)

	err := proc.ProcessPayment(50)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if mock.ChargedMoney.Amount != 50 || mock.ChargedMoney.Currency != "USD" {
		t.Errorf("expected 50 USD, got %+v", mock.ChargedMoney)
	}
}

func TestProcessor_GatewayError(t *testing.T) {
	mock := &MockGateway{ShouldFail: true}
	proc := di.NewProcessor(mock)

	err := proc.ProcessPayment(50)
	if err == nil {
		t.Fatalf("expected error from gateway, got nil")
	}
}

func TestProcessor_InvalidAmount(t *testing.T) {
	mock := &MockGateway{}
	proc := di.NewProcessor(mock)

	err := proc.ProcessPayment(-10)
	if err == nil {
		t.Fatalf("expected error for negative amount, got nil")
	}

	if mock.ChargedMoney.Amount != 0 {
		t.Errorf("gateway should not have been called, but got amount %d", mock.ChargedMoney.Amount)
	}
}

type MockContainer struct {
	gateway di.PaymentGateway
}

func (c *MockContainer) GetPaymentGateway() di.PaymentGateway {
	return c.gateway
}

func TestBadProcessor_Success(t *testing.T) {
	mock := &MockGateway{}
	container := &MockContainer{gateway: mock}
	badProc := di.NewBadProcessor(container)

	err := badProc.ProcessPayment(75)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if mock.ChargedMoney.Amount != 75 {
		t.Errorf("expected 75, got %d", mock.ChargedMoney.Amount)
	}
}
