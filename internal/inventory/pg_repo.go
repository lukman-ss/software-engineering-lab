package inventory

import (
	"context"
	"database/sql"
	"errors"
)

type pgRepo struct {
	db *sql.DB
}

func NewPGRepo(db *sql.DB) Repository {
	return &pgRepo{db: db}
}

func (r *pgRepo) GetByProductID(ctx context.Context, productID string) (*InventoryItem, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, product_id, quantity, version, created_at, updated_at FROM inventory WHERE product_id=$1`, productID)
	var item InventoryItem
	if err := row.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Version, &item.CreatedAt, &item.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInventoryNotFound
		}
		return nil, err
	}
	return &item, nil
}

// Reserve optimistically locks by version
func (r *pgRepo) Reserve(ctx context.Context, productID string, quantity int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // safe to call even after Commit

	// Select for update to lock the row
	row := tx.QueryRowContext(ctx, `SELECT id, quantity, version FROM inventory WHERE product_id=$1 FOR UPDATE`, productID)
	var id string
	var qty, version int
	if err := row.Scan(&id, &qty, &version); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInventoryNotFound
		}
		return err
	}

	if qty < quantity {
		return ErrInsufficientStock
	}

	newQty := qty - quantity
	newVersion := version + 1

	_, err = tx.ExecContext(ctx, `UPDATE inventory SET quantity=$1, version=$2, updated_at=NOW() WHERE id=$3`, newQty, newVersion, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *pgRepo) Restock(ctx context.Context, productID string, quantity int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE inventory SET quantity = quantity + $1, updated_at=NOW() WHERE product_id=$2`, quantity, productID)
	return err
}