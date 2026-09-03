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
	_, err := r.db.ExecContext(ctx, `
		UPDATE isolation_accounts SET balance = $1 WHERE id = 1;
		UPDATE isolation_accounts SET balance = $2 WHERE id = 2;
		UPDATE isolation_accounts SET balance = $3 WHERE id = 3;
		DELETE FROM isolation_transfer_audit;
	`, aliceBalance, bobBalance, charlieBalance)
	return err
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

// TransferRepeatableRead: REPEATABLE READ transaction
func (r *PostgresWalletRepo) TransferRepeatableRead(ctx context.Context, fromID, toID int, amount int64) error {
	if amount <= 0 {
		return ErrNegativeTransferAmount
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
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
		if isSerializationError(err) {
			return ErrSerializationFailure
		}
		return fmt.Errorf("deduct: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		if isSerializationError(err) {
			return ErrSerializationFailure
		}
		return fmt.Errorf("add: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO isolation_transfer_audit (from_account_id, to_account_id, amount) VALUES ($1, $2, $3)", fromID, toID, amount)
	if err != nil {
		if isSerializationError(err) {
			return ErrSerializationFailure
		}
		return fmt.Errorf("audit: %w", err)
	}

	if err := tx.Commit(); err != nil {
		if isSerializationError(err) {
			return ErrSerializationFailure
		}
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// TransferSerializable: SSI (Serializable Snapshot Isolation)
func (r *PostgresWalletRepo) TransferSerializable(ctx context.Context, fromID, toID int, amount int64) error {
	if amount <= 0 {
		return ErrNegativeTransferAmount
	}

	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
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
		if isSerializationError(err) {
			return ErrSerializationFailure
		}
		return fmt.Errorf("deduct: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE isolation_accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
	if err != nil {
		if isSerializationError(err) {
			return ErrSerializationFailure
		}
		return fmt.Errorf("add: %w", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO isolation_transfer_audit (from_account_id, to_account_id, amount) VALUES ($1, $2, $3)", fromID, toID, amount)
	if err != nil {
		if isSerializationError(err) {
			return ErrSerializationFailure
		}
		return fmt.Errorf("audit: %w", err)
	}

	if err := tx.Commit(); err != nil {
		if isSerializationError(err) {
			return ErrSerializationFailure
		}
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// TransferSerializableWithRetry retries on 40001 serialization failures with exponential backoff & jitter
func (r *PostgresWalletRepo) TransferSerializableWithRetry(ctx context.Context, fromID, toID int, amount int64, maxRetries int) error {
	baseDelay := 10 * time.Millisecond
	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := r.TransferSerializable(ctx, fromID, toID, amount)
		if err == nil {
			return nil
		}

		if IsSerializationError(err) && attempt < maxRetries {
			// Exponential backoff with full jitter
			sleep := time.Duration(float64(baseDelay) * float64(int(1)<<uint(attempt)) * (0.5 + rand.Float64()*0.5))
			select {
			case <-time.After(sleep):
				continue
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return err
	}
	return ErrMaxRetryExceeded
}

func isSerializationError(err error) bool {
	if err == nil {
		return false
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "40001"
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
