package wallet

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lukman-ss/software-engineering-lab/pkg/util"
)

type pgRepo struct {
	db *sql.DB
}

func NewPGRepo(db *sql.DB) Repository {
	return &pgRepo{db: db}
}

func (r *pgRepo) GetByUserID(ctx context.Context, userID string) (*Wallet, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, balance, currency, version, created_at, updated_at FROM wallets WHERE user_id=$1`, userID)
	var w Wallet
	if err := row.Scan(&w.ID, &w.UserID, &w.Balance, &w.Currency, &w.Version, &w.CreatedAt, &w.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, err
	}
	return &w, nil
}

// UpdateBalance adjusts balance and records a transaction
func (r *pgRepo) UpdateBalance(ctx context.Context, walletID string, amount int64, txType TransactionType, reference string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lock the wallet row
	row := tx.QueryRowContext(ctx, `SELECT balance, version FROM wallets WHERE id=$1 FOR UPDATE`, walletID)
	var balance, version int
	if err := row.Scan(&balance, &version); err != nil {
		return err
	}

	newBalance := balance + amount
	if newBalance < 0 {
		return ErrInsufficientBalance
	}

	newVersion := version + 1
	if _, err := tx.ExecContext(ctx, `UPDATE wallets SET balance=$1, version=$2, updated_at=NOW() WHERE id=$3`, newBalance, newVersion, walletID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO wallet_transactions (id, wallet_id, amount, type, reference, created_at) VALUES ($1,$2,$3,$4,$5,NOW())`,
		util.NewTxID(), walletID, amount, string(txType), reference); err != nil {
		return err
	}

	return tx.Commit()
}
