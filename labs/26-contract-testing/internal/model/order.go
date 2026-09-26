package model

// Order represents domain entity.
type Order struct {
	ID           string
	Status       string
	CustomerID   string
	CustomerName string
	Total        int64
	Notes        string
}

// OrderResponseV1 matches original V1 consumer contract expectations.
type OrderResponseV1 struct {
	ID       string             `json:"id"`
	Status   string             `json:"status"` // "IN_PROGRESS", "COMPLETED"
	Customer CustomerResponseV1 `json:"customer"`
	Total    int64              `json:"total"` // Integer amount in IDR
	Notes    string             `json:"notes,omitempty"`
}

type CustomerResponseV1 struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OrderResponseBreaking represents provider modifications that violate contract.
type OrderResponseBreaking struct {
	ID       string                   `json:"id"`
	Status   string                   `json:"status"` // "in_progress" (BREAKING: enum casing)
	Customer CustomerResponseBreaking `json:"customer"`
	Total    string                   `json:"total"` // "150000" (BREAKING: primitive type int->str)
	Notes    string                   `json:"notes,omitempty"`
}

type CustomerResponseBreaking struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"` // (BREAKING: renamed name -> full_name)
}

// OrderResponseV2 represents safe evolutionary schema.
type OrderResponseV2 struct {
	ID       string             `json:"id"`
	Status   string             `json:"status"` // "in_progress"
	Customer CustomerResponseV2 `json:"customer"`
	Total    string             `json:"total"` // "150000"
	Currency string             `json:"currency"`
}

type CustomerResponseV2 struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}
