package model

import "time"

type User struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	Age       int        `json:"age"`
	Status    string     `json:"status"` // 'active', 'suspended', 'pending'
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"` // Foreign Key to User(ID)
	TotalCents int64     `json:"total_cents"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
