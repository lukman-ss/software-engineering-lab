# Pola Transactional Outbox: Atomicitas Penulisan Database dan Event

## Problem

Pada sistem terdistribusi berbasis event, layanan sering perlu **membuat perubahan data** (misalnya, mencatat order baru) **dan** **memublikasikan event** (misalnya, `OrderCreated`) secara bersamaan. Pendekatan sembarangan melakukan dua operasi ini secara terpisah:

1. Commit perubahan ke database.
2. Publish event ke message broker.

Jika langkah (1) berhasil tetapi (2) gagal — broker down, jaringan putus, atau error lain — maka **database sudah menyimpan order tetapi event tidak pernah sampai ke broker**. Sistem berada dalam keadaan inkonsisten: data ada tapi tidak ada notifikasinya.

Masalah ini disebut **dual-write problem**.

## Why This Matters

Konsistensi data adalah fondasi sistem terdistribusi. Jika event tidak terpublish, layanan lain yang bergantung pada event tersebut — misalnya layanan notifikasi email, laporan keuangan, atau warehouse data — tidak akan pernah tahu bahwa order telah dibuat. Akibatnya:

- Pengguna tidak menerima konfirmasi order.
- Laporan keuangan tidak lengkap.
- Sistem pelacakan order gagal.
- Data di downstream service tidak pernah sinkron.

Transactional Outbox memastikan **tidak ada titik kegagalan tunggal** di antara penulisan database dan publikasi event.

## Mental Model

Bayangkan setiap perubahan data sebagai sebuah **transaksi lokal**. Di dalam transaksi yang sama, kita menuliskan:

1. **Entitas bisnis** (misalnya, `Order`) ke tabel `orders`.
2. **Record event** ke tabel `outbox` yang berada **di database yang sama**.

Jika transaksi berhasil (commit), kedua record tersebut tersedia. Jika gagal (rollback), tidak satu pun tersedia.

Setelah commit, sebuah **message relay** — proses asinkron terpisah — memantaut tabel `outbox`, mengambil record dengan status `PENDING`, dan mempublikasikannya ke message broker. Setelah berhasil dipublikasikan, relay menandai record tersebut sebagai `PROCESSED`.

Jika relay gagal setelah publish tetapi sebelum menandai status, message akan **dipublikasikan kembali** pada siklus berikutnya. Oleh karena itu, konsumen harus **idempoten** — menerima dua kali message yang sama namun hanya memproses satu kali.

## Core Concept

### Atomicitas dalam Satu Transaksi Lokal

Poli ini tidak mengandalkan transaksi terdistribusi (2PC). Ia **memindahkan publikasi event ke dalam boundary transaksi database lokal** dengan menyimpan event sebagai record di dalam tabel outbox yang sama.

```go
// service.go:42-52
func (s *OrderService) CreateOrderWithOutbox(orderID string, customerID string, amount float64) error {
	tx := s.db.BeginTx()
	order := Order{
		ID:         orderID,
		CustomerID: customerID,
		Amount:     amount,
		Status:     OrderStatusCreated,
	}

	payload, err := json.Marshal(order)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("failed to marshal order payload: %w", err)
	}

	msg := OutboxMessage{
		ID:        fmt.Sprintf("evt-%s", orderID),
		EventType: "OrderCreated",
		Payload:   string(payload),
		Status:    MessageStatusPending,
		CreatedAt: time.Now(),
	}

	if err := tx.SaveOrder(order); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.SaveOutbox(msg); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
```

Source File: `internal/outbox/service.go:18-53`

Di dalam `Commit()`, staged mutation untuk `orders` dan `outbox` dituliskan ke database yang sama secara atom:

```go
// db.go:99-116
func (tx *Tx) Commit() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.closed {
		return ErrTxClosed
	}
	tx.closed = true

	tx.db.mu.Lock()
	defer tx.db.mu.Unlock()
	for k, v := range tx.stagedOrders {
		tx.db.orders[k] = v
	}
	for k, v := range tx.stagedOutbox {
		tx.db.outbox[k] = v
	}
	return nil
}
```

Source File: `internal/outbox/db.go:99-116`

Jika `Rollback()` dipanggil, perubahan yang distagingkan tidak pernah dituliskan:

```go
// db.go:118-127
func (tx *Tx) Rollback() error {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	if tx.closed {
		return ErrTxClosed
	}
	tx.closed = true
	// discard staged mutations
	return nil
}
```

Source File: `internal/outbox/db.go:118-127`

### Message Relay (Polling Publisher)

Relay adalah worker asinkron yang memantaut tabel outbox secara berkala:

```go
// relay.go:24-37
func (r *Relay) Start() {
	go func() {
		ticker := time.NewTicker(r.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.PollAndDispatch()
			case <-r.stopChan:
				return
			}
		}
	}()
}
```

Source File: `internal/outbox/relay.go:24-37`

Pada setiap tick, relay mengambil record dengan status `PENDING`, memublikasikannya ke broker, dan menandai status menjadi `PROCESSED`:

```go
// relay.go:43-59
func (r *Relay) PollAndDispatch() int {
	pending := r.db.GetPendingOutbox()
	dispatched := 0
	for _, msg := range pending {
		err := r.broker.Publish(msg)
		if err == nil {
			err = r.db.MarkOutboxProcessed(msg.ID)
			if err != nil {
				log.Printf("failed to mark outbox msg %s as processed: %v\n", msg.ID, err)
			} else {
				dispatched++
			}
		} else {
			log.Printf("failed to publish outbox msg %s: %v\n", msg.ID, err)
		}
	}
	return dispatched
}
```

Source File: `internal/outbox/relay.go:43-59`

Relay menggunakan **polling publisher pattern** — bukan transaction log tailing (CDC). Ini berarti ada latency yang ditentukan oleh `pollInterval`.

### Idempotent Consumer

Karena pola ini memberikan **at-least-once delivery**, konsumen harus menangani duplikat. Setiap message memiliki ID unik (`evt-<orderID>`), dan konsumen melacak ID yang sudah diproses:

```go
// consumer.go:19-31
func (c *Consumer) Handle(msg OutboxMessage) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.processedIDs[msg.ID] {
		// Duplicate detected, ignore
		return false
	}

	c.processedIDs[msg.ID] = true
	c.received = append(c.received, msg)
	return true
}
```

Source File: `internal/outbox/consumer.go:19-31`

Jika message dengan ID yang sama diterima lagi, `Handle` mengembalikan `false` dan tidak menambahkan ke `received`.

## Architecture

Lab ini mengimplementasikan pola ini sebagai **simulasi in-memory** berbasis Go (pure standard library, tidak ada dependency eksternal):

| Komponen | File | Peran |
|---|---|---|
| `Order` / `OutboxMessage` | `internal/outbox/model.go` | Domain models dengan status `CREATED`/`CANCELLED` dan `PENDING`/`PROCESSED` |
| `DB` / `Tx` | `internal/outbox/db.go` | Mock database thread-safe dengan transaksi lokal (`BeginTx`, `Commit`, `Rollback`) |
| `Broker` / `MockBroker` | `internal/outbox/broker.go` | Mock message broker dengan simulasi kegagalan (`failNext`) |
| `OrderService` | `internal/outbox/service.go` | Logika bisnis: `CreateOrderWithOutbox` (atomic) vs `CreateOrderDualWriteNaive` (vulnerable) |
| `Relay` | `internal/outbox/relay.go` | Worker asinkron yang mempoling tabel outbox dan mempublikasikan ke broker |
| `Consumer` | `internal/outbox/consumer.go` | Konsumen idempotent dengan deduplication berdasarkan event ID |

Model data:

```go
// model.go:12-17
type Order struct {
	ID         string
	CustomerID string
	Amount     float64
	Status     OrderStatus
}
```

Source File: `internal/outbox/model.go:12-32`

```go
// model.go:19-32
type OutboxMessage struct {
	ID        string
	EventType string
	Payload   string
	Status    MessageStatus
	CreatedAt time.Time
}
```

Source File: `internal/outbox/model.go:26-32`

## Implementation

### Struktur Tabel Outbox

Lab menggunakan representasi sederhana dari tabel outbox yang ditemukan pada literatur (Debezium, microservices.io):

- `ID` — identifier unik event (format: `evt-<orderID>`)
- `EventType` — jenis event (misalnya, `OrderCreated`)
- `Payload` — representasi JSON dari entitas
- `Status` — `PENDING` atau `PROCESSED`
- `CreatedAt` — timestamp pembuatan

### Simulasi Database Transaksional

Database direpresentasikan sebagai struct `DB` dengan dua map: `orders` dan `outbox`. Transaksi (`Tx`) menyimpan perubahan ke dalam **buffer staged** (`stagedOrders`, `stagedOutbox`). Pada `Commit()`, staged changes dituliskan atomik ke kedua map secara bersamaan. Pada `Rollback()`, buffer dibuang.

Mutex `sync.RWMutex` pada `DB` dan `sync.Mutex` pada `Tx` memastikan thread-safety.

### Dual-Write vs Outbox

Service mengekspor dua fungsi:

1. `CreateOrderWithOutbox` — menulis order dan outbox message dalam satu transaksi.
2. `CreateOrderDualWriteNaive` — menuliskan order ke DB, commit, lalu mencoba publish ke broker di luar transaksi.

Kedua-duanya ditulis eksplisit untuk membandingkan kegagalan.

## Code Walkthrough

Demo (`cmd/demo/main.go`) berjalan tiga skenario secara berurutan:

### Skenario 1: Dual-Write Problem

```go
// demo/main.go:21-27
broker.SetFailNext(true) // broker failure simulation
err := service.CreateOrderDualWriteNaive(broker, "order-dual-fail", "cust-1", 150.0)
if err != nil {
	fmt.Printf("Direct write failed: %v\n", err)
}
_, dbFound := db.GetOrder("order-dual-fail")
fmt.Printf("State Inconsistency: Order in DB = %v, Broker Message Count = %d\n", dbFound, len(broker.GetPublished()))
```

Source File: `cmd/demo/main.go:21-27`

Broker diset agar gagal (`failNext = true`). Service menuliskan order ke DB dan sukses commit. Kemudian broker gagal. Hasilnya: order ada di DB (true), tetapi tidak ada message yang dipublikasikan (0). **Sistem inkonsisten.**

### Skenario 2: Transactional Outbox Solution

```go
// demo/main.go:35-40
err = service.CreateOrderWithOutbox("order-outbox-success", "cust-2", 300.0)
if err != nil {
	fmt.Printf("Failed to create order: %v\n", err)
	return
}
fmt.Println("Order and Outbox record atomically saved to DB.")
```

Source File: `cmd/demo/main.go:35-40`

Order dan outbox message dituliskan atomik. Relay yang sudah distart memantau tabel outbox dan mempublikasikan message ke broker.

### Skenario 3: At-Least-Once & Idempotent Consumer

```go
// demo/main.go:56-60
duplicateMsg := published[0]
acceptedAgain := consumer.Handle(duplicateMsg)
fmt.Printf("Consumer processing duplicate delivery: accepted=%v (Duplicate safely skipped!)\n", acceptedAgain)
fmt.Printf("Total events processed by consumer: %d\n", consumer.GetReceivedCount())
```

Source File: `cmd/demo/main.go:56-60`

Pesan yang sama dikirim lagi ke konsumen. Karena ID sudah tercatat, konsumen menolaknya (`accepted=false`). Total event yang diproses tetap 1.

## What the Tests Prove

Ada 5 tes yang semuanya lolos, termasuk di bawah Go race detector.

```text
=== RUN   TestTransactionalOutbox_HappyPath
--- PASS: TestTransactionalOutbox_HappyPath (0.05s)
=== RUN   TestTransactionalOutbox_Rollback
--- PASS: TestTransactionalOutbox_Rollback (0.03s)
=== RUN   TestTransactionalOutbox_Idempotency_DuplicateDelivery
--- PASS: TestTransactionalOutbox_Idempotency_DuplicateDelivery (0.00s)
=== RUN   TestDualWriteProblem_Failure
--- PASS: TestDualWriteProblem_Failure (0.00s)
=== RUN   TestTransactionalOutbox_ConcurrentWrites
--- PASS: TestTransactionalOutbox_ConcurrentWrites (0.05s)
PASS
ok  	github.com/software-engineering-lab/labs/21-outbox-pattern/tests	1.257s
```

Source: `engineering-audit/03-test-audit.md:34-47`

### Happy Path (`TestTransactionalOutbox_HappyPath`)

```go
// tests/outbox_test.go:12-59
func TestTransactionalOutbox_HappyPath(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	relay := outbox.NewRelay(db, broker, 10*time.Millisecond)
	service := outbox.NewOrderService(db)
	consumer := outbox.NewConsumer()

	relay.Start()
	defer relay.Stop()

	err := service.CreateOrderWithOutbox("o-1", "c-1", 100.50)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	order, ok := db.GetOrder("o-1")
	if !ok || order.Amount != 100.50 {
		t.Fatalf("expected order to be persisted in db")
	}

	time.Sleep(50 * time.Millisecond)

	published := broker.GetPublished()
	if len(published) != 1 {
		t.Fatalf("expected 1 message published to broker, got %d", len(published))
	}
	if published[0].EventType != "OrderCreated" {
		t.Fatalf("expected OrderCreated event")
	}

	msg, ok := db.GetOutbox("evt-o-1")
	if !ok || msg.Status != outbox.MessageStatusProcessed {
		t.Fatalf("expected outbox message to be marked processed")
	}

	processed := consumer.Handle(published[0])
	if !processed {
		t.Fatalf("expected consumer to process new message")
	}
	if consumer.GetReceivedCount() != 1 {
		t.Fatalf("expected consumer to have 1 message")
	}
}
```

Source File: `tests/outbox_test.go:11-59`

Membuktikan seluruh alur happy path: atomic save → relay dispatch → broker menerima → status berubah `PROCESSED` → consumer memproses.

### Rollback Safety (`TestTransactionalOutbox_Rollback`)

```go
// tests/outbox_test.go:61-90
func TestTransactionalOutbox_Rollback(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	relay := outbox.NewRelay(db, broker, 10*time.Millisecond)

	relay.Start()
	defer relay.Stop()

	tx := db.BeginTx()
	_ = tx.SaveOrder(outbox.Order{ID: "o-2", Amount: 50})
	_ = tx.SaveOutbox(outbox.OutboxMessage{ID: "evt-o-2", Status: outbox.MessageStatusPending})

	_ = tx.Rollback()

	_, ok := db.GetOrder("o-2")
	if ok {
		t.Fatalf("expected order to not be saved")
	}
	_, ok = db.GetOutbox("evt-o-2")
	if ok {
		t.Fatalf("expected outbox to not be saved")
	}

	time.Sleep(30 * time.Millisecond)
	if len(broker.GetPublished()) > 0 {
		t.Fatalf("expected no messages sent to broker")
	}
}
```

Source File: `tests/outbox_test.go:61-90`

Membuktikan bahwa rollback benar-benar **membuang kedua record**. Tidak ada order di DB, tidak ada outbox, dan relay tidak mempublikasikan message apapun.

### Idempotent Consumer (`TestTransactionalOutbox_Idempotency_DuplicateDelivery`)

```go
// tests/outbox_test.go:92-115
func TestTransactionalOutbox_Idempotency_DuplicateDelivery(t *testing.T) {
	consumer := outbox.NewConsumer()

	msg := outbox.OutboxMessage{
		ID:        "evt-duplicate",
		EventType: "OrderCreated",
		Payload:   "data",
	}

	p1 := consumer.Handle(msg)
	if !p1 {
		t.Fatalf("expected first delivery to be processed")
	}

	p2 := consumer.Handle(msg)
	if p2 {
		t.Fatalf("expected duplicate delivery to be rejected by consumer")
	}

	if consumer.GetReceivedCount() != 1 {
		t.Fatalf("expected exactly 1 message processed despite duplicate delivery")
	}
}
```

Source File: `tests/outbox_test.go:92-115`

Membuktikan konsumen memproses message pertama (`true`) dan menolak duplikat (`false`), dengan counter tetap 1.

### Dual-Write Failure (`TestDualWriteProblem_Failure`)

```go
// tests/outbox_test.go:117-139
func TestDualWriteProblem_Failure(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	service := outbox.NewOrderService(db)

	broker.SetFailNext(true)

	err := service.CreateOrderDualWriteNaive(broker, "o-bug", "c-2", 200)
	if err == nil {
		t.Fatalf("expected error from dual write naive approach")
	}

	_, ok := db.GetOrder("o-bug")
	if !ok {
		t.Fatalf("expected order to be saved in DB")
	}

	if len(broker.GetPublished()) != 0 {
		t.Fatalf("expected no message published to broker")
	}
}
```

Source File: `tests/outbox_test.go:117-139`

Membuktikan skenario kegagalan dual-write: DB punya order, broker tidak menerima message.

### Concurrency Safety (`TestTransactionalOutbox_ConcurrentWrites`)

```go
// tests/outbox_test.go:141-170
func TestTransactionalOutbox_ConcurrentWrites(t *testing.T) {
	db := outbox.NewDB()
	broker := outbox.NewMockBroker()
	relay := outbox.NewRelay(db, broker, 5*time.Millisecond)
	service := outbox.NewOrderService(db)

	relay.Start()
	defer relay.Stop()

	var wg sync.WaitGroup
	workers := 10
	ordersPerWorker := 10

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(wID int) {
			defer wg.Done()
			for j := 0; j < ordersPerWorker; j++ {
				_ = service.CreateOrderWithOutbox(
					"o-concurrent-1",
					"c-multi",
					float64(wID*100+j),
				)
			}
		}(i)
	}

	wg.Wait()
	time.Sleep(50 * time.Millisecond)
}
```

Source File: `tests/outbox_test.go:141-170`

10 goroutine menulis 10 order masing-masing sambil relay aktif memantau. Diuji dengan `go test -race` — **tidak ditemukan race condition**.

## Recovery / Rollback

### Recovery pada Relay

Jika relay gagal publish ke broker (misalnya, broker down), message tetap di tabel outbox dengan status `PENDING`. Pada siklus polling berikutnya, relay akan **mencoba lagi**. Karena tidak ada retry policy dengan exponential backoff di implementasi ini, relay akan mencoba setiap `pollInterval`.

Jika relay berhasil publish tapi gagal menandai status `PROCESSED` (race condition pada `MarkOutboxProcessed`), message akan tetap `PENDING` dan dipublikasikan kembali pada siklus berikutnya. Karena konsumen bersifat idempotent, publikasi duplikat tidak menyebabkan efek samping.

### Rollback pada Transaksi

Jika terjadi error sebelum `Commit()` — seperti kegagalan marshal JSON — `Rollback()` dipanggil dan tidak ada perubahan yang dituliskan ke database. Ini diuji secara eksplisit pada `TestTransactionalOutbox_Rollback`.

## Production Considerations

> **Penting**: Implementasi ini adalah **simulasi in-memory**. Berikut adalah pertimbangan untuk produksi berdasarkan batasan dan rekomendasi dari research dan engineering notes.

### Apa yang Tidak Diimplementasikan

1. **Persisten pada disk** — In-memory map tidak bertahan restart proses.
2. **Exponential backoff** pada relay — Polling dilakukan dengan interval tetap.
3. **Outbox cleanup/retention** — Record dengan status `PROCESSED` tidak pernah dihapus.
4. **CDC / Transaction Log Tailing** — Hanya polling publisher yang diimplementasikan.

### Rekomendasi untuk Skala Produksi

Berdasarkan research (Finding 6 dan Open Questions), sistem produksi perlu:

- **Monitoring metrik**: jumlah event belum diproses, usia event tertua, failure rate, retry count, throughput.
- **Cleanup**: arsip atau hapus record `PROCESSED` setelah periode retensi.
- **Indexing**: indeks pada kolom `Status` dan `CreatedAt` untuk performa polling.
- **Retry policy**: exponential backoff dan DLQ untuk event yang gagal berulang kali.
- **Schema evolution**: strategi backward compatibility untuk payload JSON/Avro.

Nilai-nilai ambang batas SLA seperti "2 detik normal" dan "47+ menit error" adalah **contoh ilustratif**, bukan standar universal. Mereka harus disesuaikan dengan objektif layanan.

## Common Mistakes

1. **Publish event di luar transaksi** — ini justru yang menyebabkan dual-write problem.
2. **Konsumen tidak idempotent** — at-least-once delivery berarti duplikat pasti terjadi.
3. **Lupa menandai event sebagai `PROCESSED`** — relay akan mempublikasikan ulang berulang-ulang.
4. **Tidak mengatur cleanup** — tabel outbox tumbuh tanpa batas.
5. **Polling interval terlalu agresif** — menyebabkan tekanan pada database.
6. **Menggunakan 2PC sebagai solusi** — tidak layak karena coupling dan kompleksitas.

## Case Study: Demo End-to-End

Demo (`go run ./cmd/demo`) menampilkan tiga skenario berturut-turut:

```text
=== Lab 21: Transactional Outbox Pattern Demo ===

[Scenario 1: The Dual-Write Problem]
Direct write failed: failed to publish to broker after DB commit: broker unavailable
State Inconsistency: Order in DB = true, Broker Message Count = 0

[Scenario 2: Transactional Outbox Solution]
Creating order with transactional outbox...
Order and Outbox record atomically saved to DB.
Broker received messages: 1
 - Event ID: evt-order-outbox-success, Type: OrderCreated, Payload: {"ID":"order-outbox-success","CustomerID":"cust-2","Amount":300,"Status":"CREATED"}
 - Consumer processing initial message: accepted=true

[Scenario 3: At-Least-Once Delivery & Idempotent Consumer]
Simulating duplicate delivery to consumer...
Consumer processing duplicate delivery: accepted=false (Duplicate safely skipped!)
Total events processed by consumer: 1

=== Demo Complete ===
```

Source: `engineering/03-execution-result.md:44-62`

**Skenario 1** membuktikan kegagalan dual-write: order tersimpan di DB, tetapi broker tidak menerima event.

**Skenario 2** membuktikan solusi outbox: order dan outbox message disimpan atomik, relay mempublikasikan ke broker, konsumen menerima message.

**Skenario 3** membuktikan idempotensi: message duplikat ditolak oleh konsumen.

## Checklist

- [x] Order dan outbox message dituliskan dalam satu transaksi
- [x] Rollback benar-benar membuang kedua record
- [x] Relay mempolling tabel outbox dan mempublikasikan ke broker
- [x] Status outbox berubah dari `PENDING` ke `PROCESSED` setelah publish
- [x] Konsumen idempotent — menolak message duplikat
- [x] Dual-write problem ditunjukkan: DB konsisten, broker tidak dapat message
- [x] Concurrency safety teruji dengan race detector
- [x] Demo dapat dijalankan dengan `go run ./cmd/demo`

## Key Takeaways

1. **Atomicitas bukan berarti sinkron**: Dua record (entitas + event) dituliskan atomik, tetapi publikasi ke broker tetap asinkron.
2. **Outbox bukan solusi exactly-once**: Ia menyediakan at-least-once delivery yang mengharuskan konsumen bersifat idempotent.
3. **Polling vs CDC**: Polling publisher lebih sederhana tetapi menambah latency. CDC (log tailing) lebih responsif tetapi database-spesifik.
4. **Konsumsi duplicate adalah fitur, bukan bug**: Relay yang crash sebelum menandai `PROCESSED` akan mempublikasikan ulang — konsumen harus siap menanganinya.
5. **Simulasi vs produksi**: Lab ini membuktikan konsep dengan simulasi in-memory. Produksi perlu cleanup, monitoring, index, dan retry policy.
6. **SLA thresholds bersifat ilustratif**: Nilai ambang seperti 2s atau 47m adalah contoh, harus disesuaikan dengan SLO.
7. **Dual-write problem nyata**: `CreateOrderDualWriteNaive` secara eksplisit memperlihatkan betapa mudahnya sistem menjadi inkonsisten.

## Sources

### Research
- `research/05-report.md` — Research Report (6 Findings, APPROVED)
- `research/06-open-questions.md` — Open Questions
- `research-audit/07-verdict.md` — Research Audit Verdict (APPROVED)
- `research-audit/06-gaps.md` — Research Gap Analysis

### Implementation
- `internal/outbox/model.go` — Domain models
- `internal/outbox/db.go` — Transactional database
- `internal/outbox/broker.go` — Mock message broker
- `internal/outbox/service.go` — Order service
- `internal/outbox/relay.go` — Polling relay
- `internal/outbox/consumer.go` — Idempotent consumer

### Tests
- `tests/outbox_test.go` — 5 test functions

### Engineering
- `engineering/01-design.md` — Engineering Design
- `engineering/02-implementation-notes.md` — Implementation Notes
- `engineering/03-execution-result.md` — Execution Results
- `engineering-audit/06-verdict.md` — Engineering Audit Verdict (APPROVED)
- `engineering-audit/03-test-audit.md` — Test Audit
