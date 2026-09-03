package isolation

import (
	"errors"
)

var (
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrAccountNotFound        = errors.New("account not found")
	ErrSerializationFailure   = errors.New("serialization failure (40001)")
	ErrMaxRetryExceeded       = errors.New("max retry attempts exceeded")
	ErrNegativeTransferAmount = errors.New("transfer amount must be positive")
	ErrDeadlockDetected       = errors.New("deadlock detected (40P01)")
	ErrInvalidMaxAttempts     = errors.New("max attempts must be positive")
	ErrSameAccountTransfer    = errors.New("source and destination account must be different")
)

type Account struct {
	ID      int    `json:"id"`
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
}

type Invoice struct {
	ID     int    `json:"id"`
	Amount int64  `json:"amount"`
	Status string `json:"status"`
}

type IsolationLevel string

const (
	LevelReadUncommitted IsolationLevel = "READ UNCOMMITTED"
	LevelReadCommitted   IsolationLevel = "READ COMMITTED"
	LevelRepeatableRead  IsolationLevel = "REPEATABLE READ"
	LevelSerializable    IsolationLevel = "SERIALIZABLE"
)

// DeterministicLockOrder sorts two IDs in ascending order to prevent deadlocks.
func DeterministicLockOrder(id1, id2 int) (firstID, secondID int) {
	if id1 < id2 {
		return id1, id2
	}
	return id2, id1
}
