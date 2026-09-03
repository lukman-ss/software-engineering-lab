package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lukman-ss/software-engineering-lab/internal/inventory"
	"github.com/lukman-ss/software-engineering-lab/internal/notification"
	"github.com/lukman-ss/software-engineering-lab/internal/payment"
	"github.com/lukman-ss/software-engineering-lab/pkg/util"
)

type appService struct {
	repo      Repository
	inventory inventory.Service
	payment   payment.Service
	notify    notification.Service
}

func NewService(repo Repository, inv inventory.Service, pay payment.Service, notif notification.Service) Service {
	return &appService{
		repo:      repo,
		inventory: inv,
		payment:   pay,
		notify:    notif,
	}
}

func (s *appService) CreateOrder(ctx context.Context, userID string, items []OrderItem) (*Order, error) {
	// Validate items
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// Calculate total
	var total int64
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("quantity must be positive")
		}
		item.Subtotal = item.UnitPrice * int64(item.Quantity)
		total += item.Subtotal
	}

	orderID := util.NewOrderID()
	now := time.Now()

	order := &Order{
		ID:          orderID,
		UserID:      userID,
		Status:      StatusPending,
		TotalAmount: total,
		Items:       items,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Set order ID on items
	for i := range items {
		items[i].ID = util.NewID()
		items[i].OrderID = orderID
		items[i].CreatedAt = now
	}

	// TODO: Persist order and items in a transaction
	// Reserve inventory
	invMap := make(map[string]int)
	for _, item := range items {
		invMap[item.ProductID] = item.Quantity
	}
	if err := s.inventory.ReserveItems(ctx, invMap); err != nil {
		return nil, fmt.Errorf("reserve inventory: %w", err)
	}

	// Emit notification (fire and forget, best effort)
	_ = s.notify.NotifyOrderCreated(ctx, userID, orderID)

	return order, nil
}

func (s *appService) GetByID(ctx context.Context, id string) (*Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *appService) MarkAsPaid(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != StatusPending {
		return ErrInvalidStatus
	}

	return s.repo.UpdateStatus(ctx, id, StatusPaid)
}

func (s *appService) CancelOrder(ctx context.Context, id string) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.Status != StatusPending {
		return ErrInvalidStatus
	}

	// Release inventory
	invMap := make(map[string]int)
	for _, item := range order.Items {
		invMap[item.ProductID] = item.Quantity
	}
	if err := s.inventory.ReleaseItems(ctx, invMap); err != nil {
		return err
	}

	return s.repo.UpdateStatus(ctx, id, StatusCancelled)
}
