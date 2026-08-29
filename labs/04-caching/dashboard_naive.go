package caching

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"
)

// DashboardNaiveService menghitung statistik secara langsung dari DB tiap request.
// Ini adalah baseline: tidak ada cache, setiap request = query DB + aggregation.
type DashboardNaiveService struct {
	db           *sql.DB
	queryCounter atomic.Int64 // untuk monitoring, bukan production-ready
}

func NewDashboardNaiveService(db *sql.DB) *DashboardNaiveService {
	return &DashboardNaiveService{db: db}
}

// GetDashboard mengembalikan statistik dashboard tanpa caching.
// Setiap pemanggilan menghitung ulang dari database.
func (s *DashboardNaiveService) GetDashboard(ctx context.Context, branchID int64) (Dashboard, error) {
	// Track query count untuk demonstration
	s.queryCounter.Add(1)

	// Combine semua query menjadi satu transaction untuk kompleksitas yang realistis
	var d Dashboard
	d.BranchID = branchID
	d.Date = Today()

	// Query 1: Invoice count hari ini
	row := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM invoices
		WHERE branch_id = $1 AND DATE(created_at) = DATE('now')
	`, branchID)
	if err := row.Scan(&d.InvoiceCountToday); err != nil {
		return Dashboard{}, fmt.Errorf("count invoices: %w", err)
	}

	// Query 2: Total revenue hari ini (tanpa aggregate di DB, manual untuk demo)
	// Dalam real world: SELECT COALESCE(SUM(amount), 0) FROM invoices WHERE...
	row = s.db.QueryRowContext(ctx, `
		SELECT SUM(amount) FROM payments p
		JOIN invoices i ON p.invoice_id = i.id
		WHERE i.branch_id = $1 AND DATE(i.created_at) = DATE('now')
	`, branchID)
	if err := row.Scan(&d.TotalRevenueToday); err != nil {
		return Dashboard{}, fmt.Errorf("sum revenue: %w", err)
	}

	// Query 3: Top mechanic (mechanic dengan paling banyak service record)
	row = s.db.QueryRowContext(ctx, `
		SELECT m.id, m.name, COUNT(s.id) as cnt
		FROM mechanics m
		JOIN service_records s ON s.mechanic_id = m.id
		JOIN invoices i ON s.invoice_id = i.id
		WHERE i.branch_id = $1 AND DATE(i.created_at) = DATE('now')
		GROUP BY m.id, m.name
		ORDER BY cnt DESC LIMIT 1
	`, branchID)
	if err := row.Scan(&d.TopMechanic.MechanicID, &d.TopMechanic.Name, &d.TopMechanic.Count); err != nil {
		return Dashboard{}, fmt.Errorf("top mechanic: %w", err)
	}

	// Query 4: Top sparepart
	row = s.db.QueryRowContext(ctx, `
		SELECT p.id, p.name, COUNT(pi.id) as cnt, COALESCE(SUM(pi.price), 0) as rev
		FROM parts p
		JOIN parts_invoices pi ON pi.part_id = p.id
		JOIN invoices i ON pi.invoice_id = i.id
		WHERE i.branch_id = $1 AND DATE(i.created_at) = DATE('now')
		GROUP BY p.id, p.name
		ORDER BY cnt DESC LIMIT 1
	`, branchID)
	if err := row.Scan(&d.TopSparepart.PartID, &d.TopSparepart.Name,
		&d.TopSparepart.Count, &d.TopSparepart.Revenue); err != nil {
		return Dashboard{}, fmt.Errorf("top sparepart: %w", err)
	}

	// Query 5: Vehicle count baru hari ini
	row = s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM vehicles
		WHERE branch_id = $1 AND DATE(created_at) = DATE('now')
	`, branchID)
	if err := row.Scan(&d.VehicleCountToday); err != nil {
		return Dashboard{}, fmt.Errorf("count vehicles: %w", err)
	}

	// Query 6: Active customer (purchased dalam 30 hari)
	row = s.db.QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT c.id) FROM customers c
		JOIN invoices i ON i.customer_id = c.id
		WHERE i.branch_id = $1 AND i.created_at >= DATE('now', '-30 days')
	`, branchID)
	if err := row.Scan(&d.ActiveCustomer); err != nil {
		return Dashboard{}, fmt.Errorf("count customers: %w", err)
	}

	// Simulate processing time (not a real benchmark)
	select {
	case <-ctx.Done():
		return Dashboard{}, ctx.Err()
	case <-time.After(50 * time.Millisecond): // simulate aggregation
	}

	return d, nil
}

// QueryCount mengembalikan jumlah query yang pernah dieksekusi.
func (s *DashboardNaiveService) QueryCount() int64 {
	return s.queryCounter.Load()
}