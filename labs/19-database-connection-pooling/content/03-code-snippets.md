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

Explanation: Pola yang benar adalah selalu melepaskan panggilan I/O jaringan di luar area di mana `db.Conn` sedang dipinjam atau belum di-*close*. `ProcessOrderSafe` menjalankan `externalCall` lebih dulu, baru meminjam koneksi untuk operasi DB singkat. `ProcessOrderUnsafeLeak` meminjam koneksi lalu memanggil `externalCall` sementara koneksi masih terbuka — menahan slot koneksi selama durasi I/O eksternal.

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

	if d.maxConnections > 0 && atomic.LoadInt32(&d.activeConns) >= d.maxConnections {
		return nil, ErrServerOverloaded
	}

	atomic.AddInt32(&d.activeConns, 1)
	atomic.AddInt32(&d.totalCreated, 1)

	return &mockConn{driver: d}, nil
}
```

Explanation: Driver memastikan jika total koneksi aktif di atas ambang `maxConnections`, inisiasi koneksi baru langsung dilempar sebagai error `ErrServerOverloaded`, yang mencerminkan batas ketat `max_connections` di dunia nyata. `connectDelay` mensimulasikan penalti handshake TCP/TLS sebelum batas dicek.

## Snippet 3 — Single-Close Guarantee on Connection

Source File: `internal/pool/mockdb.go`
Purpose: Memastikan pengurangan hitungan koneksi aktif hanya terjadi sekali walau `Close` dipanggil berulang.

```go
func (c *mockConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.closed = true
		atomic.AddInt32(&c.driver.activeConns, -1)
	}
	return nil
}
```

Explanation: Flag `closed` dijaga mutex mencegah double-decrement pada `activeConns`. Tanpa guard ini, penutupan ganda akan membuat hitungan koneksi aktif menjadi negatif dan memalsukan sisa kapasitas server. Divalidasi oleh `TestMockConnDoubleClose`.

## Snippet 4 — Connection Reuse on the Demo

Source File: `cmd/demo/main.go`
Purpose: Membandingkan eksekusi tanpa reuse (`SetMaxIdleConns(0)`) terhadap koneksi yang di-*pool*.

```go
func demoDirectOverhead() {
	driver := pool.NewMockDriver(100, 10*time.Millisecond) // 10ms handshake penalty

	unpooledDB := pool.OpenDB(driver)
	unpooledDB.SetMaxIdleConns(0)

	start := time.Now()
	for i := 0; i < 5; i++ {
		unpooledDB.Exec("SELECT 1")
	}
	fmt.Printf("Unpooled (5 requests): %v\n", time.Since(start))

	pooledDB := pool.OpenDB(driver)
	pooledDB.SetMaxIdleConns(5)
	pooledDB.SetMaxOpenConns(5)

	// Warm up pool
	pooledDB.Exec("SELECT 1")

	start = time.Now()
	for i := 0; i < 5; i++ {
		pooledDB.Exec("SELECT 1")
	}
	fmt.Printf("Pooled (5 requests): %v\n", time.Since(start))
}
```

Explanation: `SetMaxIdleConns(0)` memaksa `sql.DB` menutup koneksi setelah tiap query sehingga tiap request menanggung `connectDelay` 10ms penuh. Konfigurasi pooled menghangatkan pool lebih dulu lalu memakai ulang koneksi yang sudah terbuka. Angka absolut hasil demo bervariasi antar runtime; yang diverifikasi test adalah urutan (*pooled lebih cepat dari unpooled*) dan jumlah koneksi yang dibuat.

## Snippet 5 — Oversized Pool Against Server Limit

Source File: `cmd/demo/main.go`
Purpose: Menunjukkan penolakan koneksi saat pool klien melebihi limit server.

```go
func demoOversizedPool() {
	driver := pool.NewMockDriver(15, 0) // Database server max_connections = 15

	db := pool.OpenDB(driver)
	db.SetMaxOpenConns(50) // Client oversized pool = 50

	var wg sync.WaitGroup
	var successCount int32
	var failCount int32

	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := db.Conn(context.Background())
			if err != nil {
				atomic.AddInt32(&failCount, 1)
				return
			}
			time.Sleep(10 * time.Millisecond)
			conn.Close()
			atomic.AddInt32(&successCount, 1)
		}()
	}
	wg.Wait()

	fmt.Printf("Client attempted: 30, Succeeded: %d, Server Rejected: %d\n", successCount, failCount)
}
```

Explanation: Pool klien diizinkan 50 koneksi padahal server hanya menerima 15. Dari 30 percobaan konkuren, server menerima 15 dan menolak 15 — meniru perilaku `max_connections` yang menolak koneksi baru tanpa peduli apakah CPU/s memory server sehat.

## Snippet 6 — Starvation Test: Leak Holds the Only Connection

Source File: `tests/pool_test.go`
Purpose: Membuktikan bahwa koneksi yang ditahan selama I/O eksternal memblokir request lain hingga timeout.

```go
func TestConnectionStarvationDueToLeak(t *testing.T) {
	mockDriver := pool.NewMockDriver(10, 0)
	db := pool.OpenDB(mockDriver)
	defer db.Close()

	// Strict pool size of 1
	db.SetMaxOpenConns(1)

	svc := pool.NewOrderService(db)

	acquired := make(chan struct{})

	// Start unsafe request that holds connection for 100ms
	go func() {
		_ = svc.ProcessOrderUnsafeLeak(context.Background(), 1, func() error {
			close(acquired)
			time.Sleep(100 * time.Millisecond)
			return nil
		})
	}()

	// Wait until goroutine has acquired the connection and is in externalCall
	select {
	case <-acquired:
	case <-time.After(2 * time.Second):
		t.Fatal("leak goroutine did not acquire connection within 2s")
	}

	// Subsequent request with 20ms timeout should fail due to pool starvation
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	err := svc.ProcessOrderSafe(ctx, 2, nil)
	if err == nil {
		t.Fatalf("expected context deadline exceeded due to pool starvation, got nil")
	}
}
```

Explanation: Dengan `SetMaxOpenConns(1)`, hanya satu koneksi tersedia. Goroutine pertama memakai `ProcessOrderUnsafeLeak` dan menahan koneksi selama 100ms di dalam `externalCall`. Request kedua — yang memakai pola aman sekalipun — tidak mendapat koneksi dan berakhir dengan error karena context 20ms-nya habis. Kebocoran di satu kode merusak request yang kode-nya benar.

## Snippet 7 — Pool Locking Deadlock

Source File: `tests/pool_test.go`
Purpose: Membuktikan bahwa mengakuisisi koneksi kedua saat koneksi pertama masih dipegang pada pool berukuran 1 akan timeout.

```go
func TestPoolLockingDeadlock(t *testing.T) {
	// Pool of size 1: acquiring a second connection while holding the first results in deadlock/timeout
	mockDriver := pool.NewMockDriver(10, 0)
	db := pool.OpenDB(mockDriver)
	defer db.Close()
	db.SetMaxOpenConns(1)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	conn1, err := db.Conn(ctx)
	if err != nil {
		t.Fatalf("failed to acquire first connection: %v", err)
	}
	defer conn1.Close()

	_, err = db.Conn(ctx)
	if err == nil {
		t.Fatalf("expected context deadline exceeded when acquiring 2nd connection on pool of size 1, got nil")
	}
}
```

Explanation: Pola "pegang satu koneksi sambil minta koneksi lain" pada pool kecil berujung timeout. Ini adalah sederhana dari kasus yang dirumuskan HikariCP `pool_size = Tn × (Cm - 1) + 1` — jumlah minimum koneksi agar sekelompok thread yang sama-sama membutuhkan beberapa koneksi tidak saling menunggu selamanya.
