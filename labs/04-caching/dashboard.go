package caching

import "time"

// Clock interface for testability - allows injecting deterministic time sources.
type Clock interface {
	Now() time.Time
}

// SystemClock uses the real system clock.
type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now() }

// defaultClock is the global clock instance.
var defaultClock Clock = SystemClock{}

// SetClock allows replacing the clock (for testing).
func SetClock(clock Clock) {
	defaultClock = clock
}

// RestoreClock restores default system clock.
func RestoreClock() {
	defaultClock = SystemClock{}
}

// Today returns current date in YYYY-MM-DD format (using local time).
func Today() string {
	return TodayInLocation(defaultClock.Now(), time.Local)
}

// TodayInLocation returns current date in YYYY-MM-DD format for the given timezone.
func TodayInLocation(now time.Time, loc *time.Location) string {
	return now.In(loc).Format("2006-01-02")
}

// Dashboard menampilkan statistik workshop/bengkel.
type Dashboard struct {
	BranchID          int64         `json:"branch_id"`
	InvoiceCountToday int           `json:"invoice_count_today"`
	TotalRevenueToday float64       `json:"total_revenue_today"`
	TopMechanic       MechanicRank  `json:"top_mechanic"`
	TopSparepart      SparepartRank `json:"top_sparepart"`
	VehicleCountToday int           `json:"vehicle_count_today"`
	ActiveCustomer    int           `json:"active_customer"`
	Date              string        `json:"date"` // YYYY-MM-DD
}

// MechanicRank adalah ranking mekanik teratas
type MechanicRank struct {
	MechanicID int64  `json:"mechanic_id"`
	Name       string `json:"name"`
	Count      int    `json:"count"`
}

// SparepartRank adalah ranking sparepart teratas
type SparepartRank struct {
	PartID  int64   `json:"part_id"`
	Name    string  `json:"name"`
	Count   int     `json:"count"`
	Revenue float64 `json:"revenue"`
}
