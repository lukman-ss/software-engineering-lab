package compat

import "time"

type PhoneEntry struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Number    string `json:"number"`
	IsPrimary bool   `json:"is_primary"`
}

type User struct {
	ID        int
	Name      string
	Phone     *string // Legacy field; nil once contracted
	CreatedAt time.Time
}

// UserResponse is an enriched additive response supporting both legacy and modern clients
type UserResponse struct {
	ID     int          `json:"id"`
	Name   string       `json:"name"`
	Phone  string       `json:"phone,omitempty"` // Legacy field maintained for v1 consumers
	Phones []PhoneEntry `json:"phones"`          // New field for v2 consumers
}

// LegacyConsumerDTO models an un-upgraded client expecting only string phone
type LegacyConsumerDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// ModernConsumerDTO models an upgraded client expecting phones array
type ModernConsumerDTO struct {
	ID     int          `json:"id"`
	Name   string       `json:"name"`
	Phones []PhoneEntry `json:"phones"`
}
