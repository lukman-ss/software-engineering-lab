package wallet

import (
	"context"
	"errors"
	"fmt"

	"github.com/lukman-ss/software-engineering-lab/pkg/util"
)

type appService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &appService{repo: repo}
}

func (s *appService) GetBalance(ctx context.Context, userID string) (int64, error) {
	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}
	return wallet.Balance, nil
}

func (s *appService) Debit(ctx context.Context, userID string, amount int64, reference string) error {
	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	return s.repo.UpdateBalance(ctx, wallet.ID, -amount, TxTypePayment, reference)
}

func (s *appService) Credit(ctx context.Context, userID string, amount int64, reference string) error {
	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	return s.repo.UpdateBalance(ctx, wallet.ID, amount, TxTypeRefund, reference)
}

// Helper to create a wallet for a user (for testing)
func CreateWalletForUser(ctx context.Context, repo Repository, userID string) error {
	// In real code, this would be in a separate setup function
	return nil
}
