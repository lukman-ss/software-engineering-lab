package inventory

import (
	"context"
	"errors"
	"fmt"
)

type appService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &appService{repo: repo}
}

func (s *appService) CheckAvailability(ctx context.Context, productID string, quantity int) (bool, error) {
	item, err := s.repo.GetByProductID(ctx, productID)
	if err != nil {
		return false, err
	}
	return item.Quantity >= quantity, nil
}

func (s *appService) ReserveItems(ctx context.Context, items map[string]int) error {
	// Reserve each item
	// In production, this would handle partial reservations and failures
	for productID, quantity := range items {
		if err := s.repo.Reserve(ctx, productID, quantity); err != nil {
			// Rollback previous reservations (would need transaction boundary in real impl)
			return fmt.Errorf("reserve product %s: %w", productID, err)
		}
	}
	return nil
}

func (s *appService) ReleaseItems(ctx context.Context, items map[string]int) error {
	for productID, quantity := range items {
		if err := s.repo.Restock(ctx, productID, quantity); err != nil {
			return fmt.Errorf("restock product %s: %w", productID, err)
		}
	}
	return nil
}