package isolation

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrAccountNotFound        = errors.New("account not found")
	ErrSerializationFailure   = errors.New("serialization failure (40001)")
	ErrMaxRetryExceeded       = errors.New("max retry attempts exceeded")
	ErrNegativeTransferAmount = errors.New("transfer amount must be positive")
)

type Account struct {
	ID      int    `json:"id"`
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
}

type IsolationLevel string

const (
	LevelReadUncommitted IsolationLevel = "READ UNCOMMITTED"
	LevelReadCommitted   IsolationLevel = "READ COMMITTED"
	LevelRepeatableRead  IsolationLevel = "REPEATABLE READ"
	LevelSerializable    IsolationLevel = "SERIALIZABLE"
)

type WalletRepository interface {
	GetBalance(ctx context.Context, tx *sql.Tx, accountID int) (int64, error)
	TransferNaive(ctx context.Context, fromID, toID int, amount int64) error
	TransferWithLock(ctx context.Context, fromID, toID int, amount int64) error
	TransferRepeatableRead(ctx context.Context, fromID, toID int, amount int64) error
	TransferSerializable(ctx context.Context, fromID, toID int, amount int64) error
	TransferSerializableWithRetry(ctx context.Context, fromID, toID int, amount int64, maxRetries int) error
}
