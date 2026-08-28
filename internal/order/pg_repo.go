package order

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

func (r *pgRepo) Create(ctx context.Context, order *Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert order
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO orders (id, user_id, status, total_amount, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		order.ID, order.UserID, order.Status, order.TotalAmount, order.CreatedAt, order.UpdatedAt); err != nil {
		return err
	}

	// Insert items
	for _, item := range order.Items {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO order_items (id, order_id, product_id, quantity, unit_price, subtotal, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
			item.ID, item.OrderID, item.ProductID, item.Quantity, item.UnitPrice, item.Subtotal, item.CreatedAt); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *pgRepo) GetByID(ctx context.Context, id string) (*Order, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, status, total_amount, created_at, updated_at FROM orders WHERE id=$1`, id)
	var o Order
	if err := row.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalAmount, &o.CreatedAt, &o.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	// Load items
	rows, err := r.db.QueryContext(ctx, `SELECT id, order_id, product_id, quantity, unit_price, subtotal, created_at FROM order_items WHERE order_id=$1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.Subtotal, &item.CreatedAt); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	return &o, nil
}

func (r *pgRepo) UpdateStatus(ctx context.Context, id string, status Status) error {
	result, err := r.db.ExecContext(ctx, `UPDATE orders SET status=$1, updated_at=NOW() WHERE id=$2`, status, id)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrOrderNotFound
	}
	return nil
}
