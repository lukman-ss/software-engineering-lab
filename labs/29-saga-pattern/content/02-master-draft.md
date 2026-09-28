# Saga Pattern: Mengelola Transaksi Terdistribusi Tanpa 2PC

## Problem

Dalam arsitektur microservices dengan pola *Database-per-Service*, setiap layanan mengelola datanya sendiri secara independen. Ketika sebuah proses bisnis melibatkan beberapa layanan—misalnya pembuatan pesanan yang memerlukan pembuatan order, pemrosesan pembayaran, dan reservasi inventaris—developer tidak dapat lagi mengandalkan transaksi ACID tradisional atau protokol *Two-Phase Commit (2PC)*. 

2PC mengunci sumber daya di beberapa database secara sinkron hingga seluruh peserta memberikan suara setuju (*commit*). Hal ini mengakibatkan penurunan ketersediaan (*availability*), latensi tinggi, dan rentan terhadap kegagalan jaringan (*network partition*).

## Why This Matters

Meskipun 2PC memberikan konsistensi ketat (*strict ACID consistency*), protokol ini terbukti tidak cocok untuk sistem terdistribusi skala besar karena sifatnya yang *blocking*. Tanpa alternatif seperti Saga Pattern, sistem microservices tidak dapat menyelesaikan transaksi lintas layanan secara andal tanpa membuat seluruh database tergabung dalam satu kunci global yang rapuh.

## Mental Model

Bayangkan sebuah proses estafet atau barisan transaksi mandiri di mana setiap pelari menyelesaikan tugasnya secara independen dan langsung mencatatnya ke database lokal. 
- Jika semua pelari berhasil menyelesaikan bagiannya, transaksi selesai dengan sukses (*eventual consistency*).
- Jika ada pelari di tengah jalan yang gagal, maka pelari-pelari sebelumnya harus berlari mundur (*compensating transactions*) untuk membatalkan efek dari tindakan yang sudah terlanjur dilakukan.

## Core Concept

Saga didefinisikan sebagai urutan transaksi lokal (*sequence of local transactions*). Setiap transaksi lokal memperbarui database pada satu layanan dan mempublikasikan event atau pesan untuk memicu langkah berikutnya.

Terdapat tiga jenis transaksi utama dalam saga:
1. **Compensable Transactions:** Transaksi yang dapat dibatalkan atau dikompensasi dengan aksi sebaliknya (*opposite effect*).
2. **Pivot Transactions:** Titik balik (*point of no return*) dalam saga. Jika transaksi pivot berhasil, maka saga dijamin akan mencapai status akhir (tidak dapat dibatalkan lagi).
3. **Retryable Transactions:** Transaksi yang mengikuti pivot transaction dan dijamin sukses (bersifat idempotent).

## Failure Scenario

Katakanlah seorang pelanggan memesan produk, pembayaran berhasil diproses, namun saat reservasi inventaris dilakukan, stok ternyata kosong. 
Jika menggunakan ACID monolitik, rollback otomatis terjadi. Namun dalam arsitektur terdistribusi, pembayaran yang sudah terlanjur masuk ke payment gateway tidak bisa di-rollback secara otomatis oleh database. Saga mengatasi hal ini dengan menjalankan aksi kompensasi secara eksplisit: melakukan refund pembayaran dan membatalkan pesanan (*cancel order*).

## How It Works

Dalam implementasi lab ini, terdapat dua pendekatan utama koordinasi saga:
1. **Orchestration:** Sebuah koordinator terpusat (*Orchestrator*) mengatur urutan eksekusi langkah-langkah saga, memantau log status, dan memicu kompensasi berbasis tumpukan (*LIFO stack*) jika terjadi kegagalan.
2. **Choreography:** Layanan-layanan saling berkoordinasi secara terdesentralisasi dengan mempublikasikan dan mendengarkan event melalui *EventBus*.

## Architecture

Arsitektur lab ini terdiri dari:
- **`internal/saga`**: Berisi mesin orchestrator (`orchestrator.go`) yang mengelola log langkah dan eksekusi kompensasi LIFO, serta event bus (`choreography.go`).
- **`internal/services`**: Berisi mock service untuk Order, Payment, dan Inventory yang dilengkapi dengan kunci semantik (*semantic locks*) dan kunci idempotensi (*idempotency keys*).

## Implementation

Orchestrator mengelola daftar langkah (`Step`) dan log eksekusi (`StepLog`). Setiap langkah memiliki fungsi `Execute` dan `Compensate`.

```go
type Step struct {
	Name       string
	Execute    func(ctx context.Context) error
	Compensate func(ctx context.Context) error
}
```

## Code Walkthrough

Berikut adalah cuplikan inti dari eksekusi orchestrator di `internal/saga/orchestrator.go` yang menangani kegagalan dan pemanggilan kompensasi LIFO:

```go
	for _, step := range steps {
		err := step.Execute(ctx)
		o.mu.Lock()
		if err != nil {
			o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusFailed})
			o.mu.Unlock()

			compErr := o.compensate(context.Background(), executed)
			if compErr != nil {
				return fmt.Errorf("step %s failed: %w; compensation errors: %v", step.Name, err, compErr)
			}
			return fmt.Errorf("step %s failed: %w", step.Name, err)
		}
		o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusExecuted})
		o.mu.Unlock()
		executed = append(executed, step)
	}
```

Fungsi kompensasi berjalan secara LIFO (dari indeks terakhir ke pertama):

```go
func (o *Orchestrator) compensate(ctx context.Context, executed []Step) error {
	var compErrors []error
	for i := len(executed) - 1; i >= 0; i-- {
		step := executed[i]
		if step.Compensate != nil {
			err := step.Compensate(ctx)
			o.mu.Lock()
			if err != nil {
				o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusCompensateFailed})
				compErrors = append(compErrors, fmt.Errorf("compensation %s: %w", step.Name, err))
			} else {
				o.logs = append(o.logs, StepLog{Name: step.Name, Status: StatusCompensated})
			}
			o.mu.Unlock()
		}
	}
	if len(compErrors) > 0 {
		return fmt.Errorf("%v", compErrors)
	}
	return nil
}
```

## What the Tests Prove

Pengujian dalam `tests/saga_test.go` membuktikan:
1. **Happy Path:** Seluruh langkah eksekusi (CreateOrder -> ProcessPayment -> ReserveInventory -> ApproveOrder) berhasil dan mengubah status pesanan menjadi `APPROVED`.
2. **Failure LIFO Rollback:** Kegagalan pada reservasi inventaris memicu kompensasi pembayaran (refund) dan pembatalan pesanan secara berurutan dari belakang ke depan (LIFO).
3. **Idempotency:** Pemanggilan pembayaran berulang dengan ID yang sama berhasil tanpa menduplikasi pemotongan dana.
4. **Semantic Locking:** Pembuatan pesanan dengan ID yang sama saat status masih pending ditolak oleh kunci semantik.
5. **Concurrency Safety:** Eksekusi paralel 10 saga secara bersamaan lulus uji `-race`.

## Recovery / Rollback

Jika transaksi kompensasi mengalami kegagalan permanen (misalnya gangguan pada payment gateway saat refund), sistem tidak dapat mengandalkan rollback otomatis. Penanganan operasional yang direkomendasikan meliputi:
- Penggunaan *Dead-Letter Queues (DLQ)* untuk mengkarantina event kompensasi yang gagal.
- Peringatan (*alerting*) otomatis ke tim *on-call*.
- Rekonsiliasi manual melalui konsol admin atau penyesuaian di luar jalur (*out-of-band adjustments*).

## Production Considerations

- **Idempotensi Wajib:** Setiap langkah dalam saga (terutama *retryable* dan *compensable*) harus dirancang agar aman dipanggil berkali-kali.
- **Kurangnya Isolasi (Isolation):** Sagas tidak menyediakan tingkat isolasi "I" dalam ACID. Anomali seperti *dirty reads*, *fuzzy reads*, dan *lost updates* dapat terjadi dan memerlukan countermeasures seperti *semantic locks* atau *reread values*.

## Common Mistakes

- Mengandalkan rollback otomatis seperti pada transaksi ACID monolitik.
- Mengabaikan idempotensi pada saat retry, mengakibatkan duplikasi transaksi finansial.
- Tidak menangani kegagalan kompensasi (*compensating failure*) secara operasional.

## Case Study

Dalam skenario e-commerce checkout yang diuji pada lab ini:
- **Langkah 1:** `CreateOrder` (menciptakan pesanan dengan status `PENDING` dan menerapkan *semantic lock*).
- **Langkah 2:** `ProcessPayment` (memproses pembayaran dengan kunci idempotensi).
- **Langkah 3:** `ReserveInventory` (mengurangi stok barang).
- **Langkah 4:** `ApproveOrder` (mengubah status pesanan menjadi `APPROVED` dan menghapus *semantic lock*).
Ketika stok habis pada langkah ke-3, orchestrator secara otomatis memicu refund pembayaran dan pembatalan pesanan.

## Checklist

- [x] Pastikan setiap langkah memiliki fungsi kompensasi yang sesuai.
- [x] Terapkan kunci idempotensi pada operasi pembayaran atau mutasi data eksternal.
- [x] Gunakan *semantic locks* untuk melindungi data selama saga berjalan.
- [x] Jalankan uji konkurensi (`go test -race`).

## Key Takeaways

1. Saga menggantikan 2PC untuk arsitektur *database-per-service*.
2. Atomilitas dicapai di level saga, bukan pada setiap transaksi lokal individual.
3. Kompensasi berjalan secara LIFO (Last-In, First-Out).
4. Sagas mengorbankan isolasi demi ketersediaan (*availability*).
5. Idempotensi dan *semantic locks* adalah kunci untuk menghindari anomali data.

## Sources

- Microsoft Azure Architecture Center: *Saga Design Pattern* (https://learn.microsoft.com/en-us/azure/architecture/patterns/saga)
- Chris Richardson / Microservices.io: *Pattern: Saga* (https://microservices.io/patterns/data/saga.html)
- Chris Richardson: *Microservices Patterns* (Manning Publications, 2018)
