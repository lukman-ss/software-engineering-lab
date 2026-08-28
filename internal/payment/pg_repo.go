package payment

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

func (r *pgRepo) Create(ctx context.Context, p *Payment) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO payments (id, order_id, amount, status, idempotency_key, payment_method, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())`,
		p.ID, p.OrderID, p.Amount, p.Status, p.IdempotencyKey, p.PaymentMethod)
	return err
}

func (r *pgRepo) GetByID(ctx context.Context, id string) (*Payment, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, order_id, amount, status, idempotency_key, payment_method, external_id, paid_at, created_at, updated_at FROM payments WHERE id=$1`, id)
	var p Payment
	var ext sql.NullString
	var paidAt sql.NullTime
	if err := row.Scan(&p.ID, &p.OrderID, &p.Amount, &p.Status, &p.IdempotencyKey, &p.PaymentMethod, &ext, &paidAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPaymentNotFound
		}
		return nil, err
	}
	if ext.Valid {
		p.ExternalID = ext.String
	}
	if paidAt.Valid {
		t := paidAt.Time
		p.PaidAt = &t
	}
	return &p, nil
}

func (r *pgRepo) GetByIdempotencyKey(ctx context.Context, key string) (*Payment, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, order_id, amount, status, idempotency_key, payment_method, external_id, paid_at, created_at, updated_at FROM payments WHERE idempotency_key=$1`, key)
	var p Payment
	var ext sql.NullString
	var paidAt sql.NullTime
	if err := row.Scan(&p.ID, &p.OrderID, &p.Amount, &p.Status, &p.IdempotencyKey, &p.PaymentMethod, &ext, &paidAt, &p.CreatedAt, &p.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if ext.Valid {
		p.ExternalID = ext.String
	}
	if paidAt.Valid {
		t := paidAt.Time
		p.PaidAt = &t
	}
	return &p, nil
}

func (r *pgRepo) Update(ctx context.Context, p *Payment) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE payments SET status=$1, paid_at=$2, updated_at=NOW() WHERE id=$3`,
		p.Status, p.PaidAt, p.ID)
	return err
}