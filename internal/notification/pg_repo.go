package notification

import (
	"context"
	"database/sql"
)

type pgRepo struct {
	db *sql.DB
}

func NewPGRepo(db *sql.DB) Repository {
	return &pgRepo{db: db}
}

func (r *pgRepo) Create(ctx context.Context, n *Notification) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO notifications (id, user_id, type, payload, created_at) VALUES ($1,$2,$3,$4,NOW())`,
		n.ID, n.UserID, n.Type, n.Payload)
	return err
}

func (r *pgRepo) MarkAsSent(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE notifications SET sent_at=NOW() WHERE id=$1`, id)
	return err
}