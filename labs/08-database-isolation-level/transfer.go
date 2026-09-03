package isolation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/lib/pq"
)

type PostgresWalletRepo struct {
	db *sql.DB
}

func NewPostgresWalletRepo(db *sql.DB) *PostgresWalletRepo {
	return &PostgresWalletRepo{db: db}
}

// ResetAccounts resets account balances to initial state
func (r *PostgresWalletRepo) ResetAccounts(ctx context.Context, aliceBalance, bobBalance, charlieBalance int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx for reset: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 1", aliceBalance); err != nil {
		return fmt.Errorf("reset alice balance: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 2", bobBalance); err != nil {
		return fmt.Errorf("reset bob balance: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = 3", charlieBalance); err != nil {
		return fmt.Errorf("reset charlie balance: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM isolation_transfer_audit"); err != nil {
		return fmt.Errorf("clear audit: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit reset: %w", err)
	}
	return nil
}

// GetAccount reads single account details
func (r *PostgresWalletRepo) GetAccount(ctx context.Context, id int) (*Account, error) {
	var a Account
	err := r.db.QueryRowContext(ctx, "SELECT id, owner, balance FROM isolation_accounts WHERE id = $1", id).
		Scan(&a.ID, &a.Owner, &a.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, fmt.Errorf("get account %d: %w", id, err)
	}
	return &a, nil
}

// GetBalance reads current balance with ordinary non-locking SELECT
func (r *PostgresWalletRepo) GetBalance(ctx context.Context, tx *sql.Tx, accountID int) (int64, error) {
	var balance int64
	var queryRunner interface {
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	}
	if tx != nil {
		queryRunner = tx
	} else {
		queryRunner = r.db
	}

	err := queryRunner.QueryRowContext(ctx, "SELECT balance FROM isolation_accounts WHERE id = $1", accountID).Scan(&balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrAccountNotFound
		}
		return 0, fmt.Errorf("get balance %d: %w", accountID, err)
	}
	return balance, nil
}

// TransferNaive: READ COMMITTED without FOR UPDATE
// Vulnerable to lost update / race condition when multiple transfers read stale balance
func (r *PostgresWalletRepo) TransferNaive(ctx context.Context, fromID, toID int, amount int64) error {
	if fromID == toID {
		return ErrSameAccountTransfer
	}
	if amount <= 0 {
		return ErrNegativeTransferAmount
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	// 1. Read sender balance (ordinary SELECT - no lock!)
	fromBalance, err := r.GetBalance(ctx, tx, fromID)
	if err != nil {
		return err
	}

	// 2. Validate in application logic
	if fromBalance < amount {
		return ErrInsufficientFunds
	}

	// 3. Read receiver balance
	toBalance, err := r.GetBalance(ctx, tx, toID)
	if err != nil {
		return err
	}

	// 4. Overwrite sender balance with calculated value (Lost Update vulnerability)
	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = $2", fromBalance-amount, fromID)
	if err != nil {
		return fmt.Errorf("update from balance: %w", err)
	}

	// 5. Overwrite receiver balance with calculated value
	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = $1 WHERE id = $2", toBalance+amount, toID)
	if err != nil {
		return fmt.Errorf("update to balance: %w", err)
	}

	// 6. Record audit
	_, err = tx.ExecContext(ctx, "INSERT INTO isolation_transfer_audit (from_account_id, to_account_id, amount) VALUES ($1, $2, $3)", fromID, toID, amount)
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}

	return tx.Commit()
}

// TransferWithLock: READ COMMITTED + SELECT ... FOR UPDATE
// Acquires locks in deterministic order (smaller ID first) to prevent deadlock
func (r *PostgresWalletRepo) TransferWithLock(ctx context.Context, fromID, toID int, amount int64) error {
	if fromID == toID {
		return ErrSameAccountTransfer
	}
	if amount <= 0 {
		return ErrNegativeTransferAmount
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	firstID, secondID := DeterministicLockOrder(fromID, toID)

	// Lock both accounts deterministically
	var b1, b2 int64
	err = tx.QueryRowContext(ctx, "SELECT balance FROM isolation_accounts WHERE id = $1 FOR UPDATE", firstID).Scan(&b1)
	if err != nil {
		return fmt.Errorf("lock account %d: %w", firstID, err)
	}
	err = tx.QueryRowContext(ctx, "SELECT balance FROM isolation_accounts WHERE id = $1 FOR UPDATE", secondID).Scan(&b2)
	if err != nil {
		return fmt.Errorf("lock account %d: %w", secondID, err)
	}

	var senderBalance int64
	if firstID == fromID {
		senderBalance = b1
	} else {
		senderBalance = b2
	}

	if senderBalance < amount {
		return ErrInsufficientFunds
	}

	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
	if err != nil {
		return fmt.Errorf("deduct balance: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		return fmt.Errorf("add balance: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO isolation_transfer_audit (from_account_id, to_account_id, amount) VALUES ($1, $2, $3)", fromID, toID, amount)
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}

	return tx.Commit()
}

// WrapTxError wraps transaction errors with sentinel errors while preserving the original error chain.
// It preserves both errors.Is(err, ErrSerializationFailure) and errors.As(err, &pq.Error) semantics.
func WrapTxError(action string, err error) error {
	if err == nil {
		return nil
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "40001":
			return fmt.Errorf("%s: %w: %w", action, ErrSerializationFailure, err)
		case "40P01":
			return fmt.Errorf("%s: %w: %w", action, ErrDeadlockDetected, err)
		}
	}
	return fmt.Errorf("%s: %w", action, err)
}

// IsRetryableTxError checks whether an error is transient and safe to retry (40001 serialization failure or 40P01 deadlock).
func IsRetryableTxError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrSerializationFailure) || errors.Is(err, ErrDeadlockDetected) {
		return true
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "40001" || pqErr.Code == "40P01"
	}
	return false
}

// IsSerializationError is exposed for tests and callers.
func IsSerializationError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ErrSerializationFailure) {
		return true
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "40001"
	}
	return false
}

// TransferRepeatableRead: REPEATABLE READ transaction
func (r *PostgresWalletRepo) TransferRepeatableRead(ctx context.Context, fromID, toID int, amount int64) error {
	if fromID == toID {
		return ErrSameAccountTransfer
	}
	if amount <= 0 {
		return ErrNegativeTransferAmount
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return WrapTxError("begin tx", err)
	}
	defer tx.Rollback()

	fromBalance, err := r.GetBalance(ctx, tx, fromID)
	if err != nil {
		return err
	}

	if fromBalance < amount {
		return ErrInsufficientFunds
	}

	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
	if err != nil {
		return WrapTxError("deduct", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		return WrapTxError("add", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO isolation_transfer_audit (from_account_id, to_account_id, amount) VALUES ($1, $2, $3)", fromID, toID, amount)
	if err != nil {
		return WrapTxError("audit", err)
	}

	if err := tx.Commit(); err != nil {
		return WrapTxError("commit", err)
	}

	return nil
}

// TransferSerializable: SSI (Serializable Snapshot Isolation)
func (r *PostgresWalletRepo) TransferSerializable(ctx context.Context, fromID, toID int, amount int64) error {
	if fromID == toID {
		return ErrSameAccountTransfer
	}
	if amount <= 0 {
		return ErrNegativeTransferAmount
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return WrapTxError("begin tx", err)
	}
	defer tx.Rollback()

	fromBalance, err := r.GetBalance(ctx, tx, fromID)
	if err != nil {
		return err
	}

	if fromBalance < amount {
		return ErrInsufficientFunds
	}

	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
	if err != nil {
		return WrapTxError("deduct", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		return WrapTxError("add", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO isolation_transfer_audit (from_account_id, to_account_id, amount) VALUES ($1, $2, $3)", fromID, toID, amount)
	if err != nil {
		return WrapTxError("audit", err)
	}

	if err := tx.Commit(); err != nil {
		return WrapTxError("commit", err)
	}

	return nil
}

type TxOperation func(ctx context.Context) error

type RetryPolicy struct {
	BaseDelay time.Duration
	MaxDelay  time.Duration
	RandSrc   func() float64
	Sleep     func(ctx context.Context, d time.Duration) error
}

func defaultSleep(ctx context.Context, d time.Duration) error {
	select {
	case <-time.After(d):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func RetryTransaction(
	ctx context.Context,
	maxAttempts int,
	policy RetryPolicy,
	operation TxOperation,
) error {
	if maxAttempts <= 0 {
		return ErrInvalidMaxAttempts
	}
	if policy.RandSrc == nil {
		policy.RandSrc = rand.Float64
	}
	if policy.BaseDelay <= 0 {
		policy.BaseDelay = 10 * time.Millisecond
	}
	if policy.MaxDelay <= 0 {
		policy.MaxDelay = 1 * time.Second
	}
	if policy.Sleep == nil {
		policy.Sleep = defaultSleep
	}

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := operation(ctx)
		if err == nil {
			return nil
		}
		lastErr = err

		if IsRetryableTxError(err) && attempt < maxAttempts {
			delayUpper := float64(policy.BaseDelay) * float64(int(1) << uint(attempt-1))
			if delayUpper > float64(policy.MaxDelay) {
				delayUpper = float64(policy.MaxDelay)
			}
			delay := time.Duration(delayUpper * policy.RandSrc())
			policy.Sleep(ctx, delay)
			continue
		}
		break
	}

	if IsRetryableTxError(lastErr) {
		return fmt.Errorf("%w: %w", ErrMaxRetryExceeded, lastErr)
	}
	return lastErr
}

func (r *PostgresWalletRepo) TransferSerializableWithRetry(ctx context.Context, fromID, toID int, amount int64, maxAttempts int) error {
	return r.TransferSerializableWithRetryPolicy(ctx, fromID, toID, amount, maxAttempts, RetryPolicy{
		BaseDelay: 10 * time.Millisecond,
		MaxDelay:  1 * time.Second,
		RandSrc:   rand.Float64,
	})
}

func (r *PostgresWalletRepo) TransferSerializableWithRetryPolicy(ctx context.Context, fromID, toID int, amount int64, maxAttempts int, policy RetryPolicy) error {
	repo := r
	return RetryTransaction(ctx, maxAttempts, policy, func(ctx context.Context) error {
		return repo.TransferSerializable(ctx, fromID, toID, amount)
	})
}
