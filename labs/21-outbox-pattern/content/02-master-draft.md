# Pola Transactional Outbox: Mengatasi Dual-Write Problem dengan Atomicity Database

## Problem
Dalam sistem mikrolayanan, kasus penggunaan umum melibatkan pembaruan database (mis: menyimpan pesanan) diikuti dengan pengiriman peristiwa ke message broker (mis: `OrderCreated` ke Kafka). Operasi ini tidak dapat dibuat atomik karena transaksi database hanya mengontrol perubahan pada database itu sendiri; ia tidak dapat melakukan `ROLLBACK` pada sistem eksternal seperti RabbitMQ atau endpoint HTTP. Jika commit database berhasil tetapi pengiriman ke broker gagal, terjadi ketidakkonsistenan: data tersimpan di database tetapi pesan terkait tidak pernah terkirim. Kondisi ini dikenal sebagai **dual-write problem**.

## Mengapa Masalah Ini Penting?
Ketidakkonsistenan antar layanan mengakibatkan data yang tidak sinkron antar sistem, memicu:
- Laporan bisnis yang tidak akurat karena event hilang
- Proses downstream (inventory, pengiriman) tidak terpicu
- Mekanisme penyelesaian yang kompleks dan berbasis manual untuk mendeteksi dan memperbaiki ketidaksesuaian
Transaksi terdistribusi dua-fase (2PC) bukan solusi yang praktis karena banyak database dan broker tidak mendukungnya, serta memperkenalkan kompleksitas operasional dan ketergantungan yang tidak diinginkan.

## Model Mental
Pola Transactional Outbox mengalihkan tugas pengiriman pesan dari operasi yang terpisah menjadi bagian dari transaksi database yang sama. Alih-alih menulis ke database lalu ke broker dalam dua langkah, aplikasi:
1. Memperbarui entitas bisnis (mis: tabel `orders`)
2. Menyimpan pesan ke tabel `outbox_events` — **dalam transaksi database yang sama**

Karena kedua operasi berbagi satu batas transaksi kommit/rollback, keduanya sukses bersama atau gagal bersama. Setelah transaksi commit, proses asynchronous terpisah yang disebut **message relay** membaca tabel outbox dan mengirimkan pesan ke broker. Jika relay gagal (mis: jaringan down), pesan tetap aman di tabel outbox untuk diretry kemudian. Ini menjamin bahwa **tidak ada keadaan dimana database berubah tetapi pesan hilang**.

## Konsep Inti
### Tabel Outbox
Tabel outbox biasanya memiliki kolom:
- `id`: UUID unik untuk deteksi duplikasi
- `aggregateid`: ID akar agregat (mis: OrderID) — sering digunakan sebagai kunci Kafka untuk pemartisian
- `aggregatetype`: tipe agregat (mis: "Order")
- `type`: tipe peristiwa (mis: "OrderCreated")
- `payload`: JSON bersifat event state (mis: Order lengkap sebagai JSON)

> **Catatan Implementasi Lab**: Untuk kesederhanaan, `OutboxMessage` di lab ini hanya menyertakan `ID`, `EventType`, `Payload`, `Status`, dan `CreatedAt` (lihat `model.go:26-32`). Kolom `aggregatetype` dan `aggregateid` dihilangkan karena lab ini tidak mengimplementasikan pemartisian Kafka atau pemetaan aggregate — mereka ada di desain production/Debezium, bukan di demonstrasi lab ini.

### Implementasi Relay
Lab ini menerapkan pola **Polling Publisher** untuk message relay:
- Goroutine berjalan periodik (ticker) mengambil semua record dengan status `PENDING`
- Setiap record dipublish ke broker melalui antarmuka `Broker.Publish()`
- Jika publish sukses, record ditandai sebagai `PROCESSED`
- Jika publish gagal, record tetap `PENDING` untuk percobaan berikutnya
- Interval polling dan ukuran batch dapat dikonfigurasi; nilai yang lebih kecil mengurangi latensi tetapi meningkatkan beban database

### Konsumen Idempotent
Karena relay mungkin mengirimkan pesan lebih dari sekali (mis: crash setelah publish tetapi sebelum menandai sebagai processed), konsumen harus idempotent. Dalam implementasi ini:
- Konsumen menjaga `map[string]bool` dari `processedIDs`
- Saat `Handle(msg)` dipanggil:
  - Jika `msg.ID` sudah ada di map, kembalikan `false` (duplikat, diabaikan)
  - Jika belum, tandai sebagai diproses, simpan pesan, kembalikan `true`

## Arsitektur
Berikut diagram aliran komponen dalam lab ini:
```text
[ Klien ] → [ Layanan Pemesanan ]
                     |
                     v (Transaksi Tunggal DB)
       +-----------------------------+
       |  Database Transaksional     |
       |  - tabel orders             |
       |  - tabel outbox_events      |
       +-----------------------------+
                     ^
                     | Poll (SELECT ... WHERE status='PENDING')
             [ Outbox Relay ]
                     |
                     v Publish
             [ Message Broker (Mock) ]
                     |
                     v Deliver
            [ Konsumen Idempotent ] → (Penyimpanan message_log internal)
```

## Implementasi
Kode sebenarnya berlokasi di bawah package `internal/outbox`:
- `model.go`: struktur domain `Order`, `OutboxMessage` dan konstanta status
- `db.go`: database transaksional dalam memori dengan dukungan `BeginTx`, `Commit`, `Rollback`, dan pengelolaan stage untuk operasi bersifat atomic
- `service.go`: logika bisnis yang membandingkan:
  - `CreateOrderWithOutbox`: penulisan atomic (order + outbox dalam satu Tx)
  - `CreateOrderDualWriteNaive`: penulisan naif (DB commit terpisah dari broker publish) — menunjukkan kerentanan
- `relay.go`: pekerja polling asinkron yang mengimplementasikan `PollAndDispatch()`
- `consumer.go`: konsumen yang melacak ID yang telah diproses untuk deduplikasi
- `broker.go`: antarmuka broker mock dengan kemampuan menyimulasikan kegagalan jaringan melalui `SetFailNext()`

## Penjelasan Kode
### Transaksi Atomik (service.go:17-53)
```go
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
**Tujuan**: Menunjukkan bahwa `order` dan `outbox_event` disimpan bersama dalam satu transaksi. Jika ada kesalahan selama marshal JSON atau penyimpanan outbox, seluruh transaksi di-rollback.

### Dual-Write Naif (service.go:55-90)
```go
func (s *OrderService) CreateOrderDualWriteNaive(broker Broker, orderID string, customerID string, amount float64) error {
    tx := s.db.BeginTx()
    order := Order{ID: orderID, CustomerID: customerID, Amount: amount, Status: OrderStatusCreated}

    if err := tx.SaveOrder(order); err != nil {
        _ = tx.Rollback()
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    // Direct write to broker outside DB transaction
    msg := OutboxMessage{
        ID:        fmt.Sprintf("evt-%s", orderID),
        EventType: "OrderCreated",
        Payload:   fmt.Sprintf(`{"ID":"%s","Amount":%f}`, orderID, amount),
        Status:    MessageStatusPending,
        CreatedAt: time.Now(),
    }

    if err := broker.Publish(msg); err != nil {
        return fmt.Errorf("failed to publish to broker after DB commit: %w", err)
    }

    return nil
}
```
**Tujuan**: Meniru praktik berisiko di mana write ke broker terjadi setelah commit database, sehingga rentan terhadap kegagalan broker yang menyisakan DB berubah tanpa pesan terkait.

### Polling Relay (relay.go:50-67)
```go
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
**Tujuan**: Mengambil semua pesan pending, berusaha publish ke broker, dan hanya menandai sebagai processed jika publish sukses. Kegagalan publish tidak mengubah status, sehingga pesan dapat dicoba lagi pada poll berikutnya.

### Konsumen Idempotent (consumer.go:19-31)
```go
func (c *Consumer) Handle(msg OutboxMessage) bool {
    c.mu.Lock()
    defer c.mu.Unlock()

    if c.processedIDs[msg.ID] {
        // Duplikat terdeteksi, abaikan
        return false
    }

    c.processedIDs[msg.ID] = true
    c.received = append(c.received, msg)
    return true
}
```
**Tujuan**: Menjamin bahwa meskipun relay mengirimkan pesan duplikat (at-least-once delivery), logika bisnis hanya dieksekusi sekali per ID unik.

## Apa yang Dibuktikan oleh Tes
Suite `tests/outbox_test.go` menyebarkan 8 kasus uji yang memverifikasi perilaku di atas:
1. `TestTransactionalOutbox_HappyPath`: order + outbox atomik, relay mengirimkan, consumer memproses
2. `TestTransactionalOutbox_Rollback`: rollback menghapus order dan outbox, tidak ada pesan ke broker
3. `TestTransactionalOutbox_Idempotent_DuplicateDelivery`: konsumen mengabaikan pengiriman duplikat
4. `TestDualWriteProblem_Failure`: pendekatan dual-write naif menghasilkan order di DB tanpa pesan ke broker saat broker down
5. `TestTransactionalOutbox_ConcurrentWrites`: 10 goroutine menyimbolayanan pesanan concurrency menghasilkan tepat 100 pesan terpublish, 0 pending, tanpa lomba data
6. `TestTransactionalOutbox_PurgeProcessed`: membersihkan record outbox yang sudah diproses
7. `TestTransactionalOutbox_RelayRetryAfterBrokerFailure`: relay melakukan retry otomatis setelah kegagalan broker sementara
8. `TestTransactionalOutbox_ConcurrentConsumers`: banyak konsumen mengakses map pemrosesan secara konkuren masih menghasilkan deduplikasi yang benar

Semua tes lolos dengan `go test -race ./...`, memastikan tidak ada lomba data.

## Pertimbangan Produksi
1. **Operasional Outbox**: Tabel outbox harus dibersihkan (diarsipkan atau dihapus) setelah event diproses untuk mencegah pertumbuhan tak terbatas. Contoh: 1 juta event/hari → tabel sangat besar setelah setahun jika tidak dibersihkan.
2. **Monitoring**: Metrik kunci yang perlu diamati:
   - Jumlah event outbox yang belum diproses
   - Usia tertua event outbox yang belum diproses (indikator keterlambatan relay atau masalah broker)
   - Tingkat kegagalan publish
   - Jumlah retry
   - Laju pemrosesan (event/detik)
   Ambang batas contoh dari spesifikasi lab: usia tertua event ~2 detik dalam kondisi normal; 47+ menit menunjukkan masalah serius. Nilai ini bersifat ilustratif dan harus di-tuning sesuai SLO layanan masing-masing.
3. **Desain Payload**: Untuk menghindari biaya I/O dan serialisasi yang tinggi, pertimbangkan menyimpan hanya ID penting (thin event) dan biarkan konsumen mengambil keadaan lengkap bila diperlukan, alih-alih mengirimkan seluruh objek besar (fat event).
4. **Batasan Implementasi Lab**:
   - Database adalah implementasi dalam memori untuk keperluan demonstrasi, bukan pengganti PostgreSQL/MySQL produksi
   - Hanya pola Polling Publisher yang diimplementasikan; Transaction Log Tailing (mis: Debezium membaca WAL Postgres) tidak ditunjukkan
   - Tidak ada mekanisme dead-letter queue untuk payload yang tidak dapat diproses secara berulang
5. **Jaminan Pengiriman**: Pola ini memberikan pengiriman **at-least-once** (bukan exactly-once). Karena itu, konsumen **wajib** idempotent untuk menangani pengiriman duplikat dengan aman.

## Kesalahan Umum
1. Menganggap bahwa outbox memberi jaminan exactly-once delivery — padahal hanya at-least-once; duplikat masih mungkin terjadi
2. Menyimpan payload besar (mis: objek dengan 500 field) di tabel outbox, menyebabkan bloat dan biaya I/O tinggi
3. Melupakan membersihkan tabel outbox, mengakibatkan pertumbuhan tak terbatas dan penurunan performa query
4. Mengabaikan monitoring usia tertua event outbox yang belum diproses, yang merupakan indikator awal masalah dalam saluran relay → broker
5. Mengasumsikan bahwa konsumen tidak perlu idempotent karena berasa "relay saya sangat andal"