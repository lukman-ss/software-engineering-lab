package services

import (
	"errors"
	"sync"
)

type OrderState string

const (
	OrderPending   OrderState = "PENDING"
	OrderApproved  OrderState = "APPROVED"
	OrderCancelled OrderState = "CANCELLED"
)

type OrderService struct {
	mu     sync.Mutex
	orders map[string]OrderState
	locks  map[string]bool // semantic lock
}

func NewOrderService() *OrderService {
	return &OrderService{
		orders: make(map[string]OrderState),
		locks:  make(map[string]bool),
	}
}

func (s *OrderService) CreateOrder(orderID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.locks[orderID] {
		return errors.New("order is semantically locked by another operation")
	}

	s.orders[orderID] = OrderPending
	s.locks[orderID] = true
	return nil
}

func (s *OrderService) ApproveOrder(orderID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.orders[orderID] == OrderCancelled {
		return errors.New("cannot approve cancelled order")
	}
	s.orders[orderID] = OrderApproved
	delete(s.locks, orderID)
	return nil
}

func (s *OrderService) CancelOrder(orderID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[orderID] = OrderCancelled
	delete(s.locks, orderID)
	return nil
}

func (s *OrderService) GetOrder(orderID string) OrderState {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.orders[orderID]
}

type PaymentService struct {
	mu          sync.Mutex
	payments    map[string]float64
	processedID map[string]bool // idempotency key
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		payments:    make(map[string]float64),
		processedID: make(map[string]bool),
	}
}

func (s *PaymentService) ProcessPayment(paymentID string, amount float64, shouldFail bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.processedID[paymentID] {
		return nil // Idempotent success
	}

	if shouldFail {
		return errors.New("insufficient funds")
	}

	s.payments[paymentID] = amount
	s.processedID[paymentID] = true
	return nil
}

func (s *PaymentService) RefundPayment(paymentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.payments, paymentID)
	return nil
}

func (s *PaymentService) HasPayment(paymentID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.payments[paymentID]
	return exists
}

type InventoryService struct {
	mu        sync.Mutex
	stock     map[string]int
	reserved  map[string]int
}

func NewInventoryService(initialStock map[string]int) *InventoryService {
	return &InventoryService{
		stock:    initialStock,
		reserved: make(map[string]int),
	}
}

func (s *InventoryService) Reserve(item string, qty int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stock[item] < qty {
		return errors.New("out of stock")
	}

	s.stock[item] -= qty
	s.reserved[item] += qty
	return nil
}

func (s *InventoryService) Release(item string, qty int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.reserved[item] < qty {
		return errors.New("cannot release more than reserved")
	}

	s.reserved[item] -= qty
	s.stock[item] += qty
	return nil
}

func (s *InventoryService) GetStock(item string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.stock[item]
}
