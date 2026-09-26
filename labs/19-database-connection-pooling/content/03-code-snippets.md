## Snippet 1 — Safe Processing vs Unsafe Leak

Source File: `internal/pool/service.go`
Purpose: Menganalogikan pemisahan pemanggilan eksternal dari masa penahanan koneksi versus melakukan panggilan saat koneksi masih ditahan.

```go
// ProcessOrderSafe performs external I/O outside of the DB transaction/connection.
func (s *OrderService) ProcessOrderSafe(ctx context.Context, orderID int, externalCall func() error) error {
	// 1. External Call first (or after), never while holding the connection
	if externalCall != nil {
		if err := externalCall(); err != nil {
			return err
		}
	}

	// 2. Short, bounded DB operation
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx, "UPDATE orders SET status = 'completed' WHERE id = ?", orderID)
	return err
}

// ProcessOrderUnsafeLeak holds the DB connection open during an external network call.
func (s *OrderService) ProcessOrderUnsafeLeak(ctx context.Context, orderID int, externalCall func() error) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.ExecContext(ctx, "UPDATE orders SET status = 'processing' WHERE id = ?", orderID)
	if err != nil {
		return err
	}

	// External call while holding connection open!
	if externalCall != nil {
		if err := externalCall(); err != nil {
			return err
		}
	}

	return nil
}
```

Explanation: Pola yang benar adalah selalu melepaskan panggilan I/O jaringan di luar area di mana `db.Conn` sedang dipinjam atau belum di-*close*.

## Snippet 2 — Enforcing Server Connection Limits

Source File: `internal/pool/mockdb.go`
Purpose: Mensimulasikan server database menolak koneksi ketika batasan perangkat keras terlewati.

```go
func (d *MockDriver) Open(name string) (driver.Conn, error) {
	if d.connectDelay > 0 {
		time.Sleep(d.connectDelay)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.maxConnections > 0 && d.activeConns >= d.maxConnections {
		return nil, ErrServerOverloaded
	}

	atomic.AddInt32(&d.activeConns, 1)
	atomic.AddInt32(&d.totalCreated, 1)

	return &mockConn{driver: d}, nil
}
```

Explanation: Driver memastikan jika total koneksi aktif di atas ambang `maxConnections`, inisiasi koneksi baru langsung dilempar sebagai error, yang mencerminkan batas ketat yang terjadi di dunia nyata.