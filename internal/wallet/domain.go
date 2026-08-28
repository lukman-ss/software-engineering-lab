package wallet

import (
	"context"
	"errors"
	"time"
)

var (
	ErrWalletNotFound = errors.New("wallet not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type Wallet struct {
	ID        string
	UserID    string
	Balance   int64
	Currency  string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionType string

const (
	TxTypeDeposit  TransactionType = "deposit"
	TxTypeWithdraw TransactionType = "withdraw"
	TxTypePayment  TransactionType = "payment"
	TxTypeRefund   TransactionType = "refund"
)

type WalletTransaction struct {
	ID        string
	WalletID  string
	Amount    int64 // positive for deposit, negative for withdraw
	Type      TransactionType
	Reference string // e.g. OrderID or PaymentID
	CreatedAt time.Time
}

type Repository interface {
	GetByUserID(ctx context.Context, userID string) (*Wallet, error)
	// UpdateBalance should handle concurrent updates and append a transaction record atomically
	UpdateBalance(ctx context.Context, walletID string, amount int64, txType TransactionType, reference string) error
}

type Service interface {
	GetBalance(ctx context.Context, userID string) (int64, error)
	Debit(ctx context.Context, userID string, amount int64, reference string) error
	Credit(ctx context.Context, userID string, amount int64, reference string) error
}