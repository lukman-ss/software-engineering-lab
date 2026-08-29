package caching

import "time"

// Dashboard menampilkan statistik workshop/bengkel.
type Dashboard struct {
	BranchID          int64       `json:"branch_id"`
	InvoiceCountToday int         `json:"invoice_count_today"`
	TotalRevenueToday float64     `json:"total_revenue_today"`
	TopMechanic       MechanicRank `json:"top_mechanic"`
	TopSparepart      SparepartRank `json:"top_sparepart"`
	VehicleCountToday int         `json:"vehicle_count_today"`
	ActiveCustomer    int         `json:"active_customer"`
	Date              string      `json:"date"` // YYYY-MM-DD
}

// MechanicRank adalah ranking mekanik teratas
type MechanicRank struct {
	MechanicID int64  `json:"mechanic_id"`
	Name       string `json:"name"`
	Count      int    `json:"count"`
}

// SparepartRank adalah ranking sparepart teratas
type SparepartRank struct {
	PartID  int64  `json:"part_id"`
	Name    string `json:"name"`
	Count   int    `json:"count"`
	Revenue float64 `json:"revenue"`
}

// Today returns date in YYYY-MM-DD format
func Today() string {
	return time.Now().Format("2006-01-02")
}