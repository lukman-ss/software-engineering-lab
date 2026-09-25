# Diagrams

## Diagram 1 — Siklus Tunggu Deadlock (Circular Wait)

Diagram ini mengilustrasikan kondisi *deadlock* pada pemanggilan `TransferNaive` yang terjadi ketika dua goroutine saling menunggu lock yang dipegang pihak lain.

```text
Goroutine 1: Transfer(A -> B)                 Goroutine 2: Transfer(B -> A)
==============================                 ==============================
  1. Lock(Account A) [BERHASIL]                  1. Lock(Account B) [BERHASIL]
           |                                              |
           v                                              v
  2. Tunggu delay transaksi                      2. Tunggu delay transaksi
           |                                              |
           v                                              v
  3. Lock(Account B)                             3. Lock(Account A)
     [TERBLOKIR: B dipegang G2]                     [TERBLOKIR: A dipegang G1]
           |                                              |
           +-------------------> [DEADLOCK] <-------------+
                                       |
                                       v
                     Context Timeout / Deadlock Detection
                                       |
                                       v
                     Salah satu dipilih menjadi VICTIM
                          (Return: ErrDeadlock)
```

---

## Diagram 2 — Pencegahan Deadlock dengan Urutan Penguncian (Lock Ordering)

Diagram ini menggambarkan alur kerja `TransferOrdered` di mana goroutine selalu meminta lock dalam urutan alfabetis `Account.ID`.

```text
Goroutine 1: TransferOrdered(A, B)            Goroutine 2: TransferOrdered(B, A)
==================================            ==================================
Urutkan: ID "A" < ID "B"                      Urutkan: ID "A" < ID "B"
Urutan lock: Akun A -> Akun B                 Urutan lock: Akun A -> Akun B
              |                                             |
              v                                             v
  1. Bersaing memperebutkan Lock(Akun A)        1. Bersaing memperebutkan Lock(Akun A)
              |                                             |
     [Mendapatkan Lock A]                           [Menunggu Lock A]
              |                                             |
  2. Kerjakan transfer & Lock(Akun B)                       |
              |                                             |
  3. Lepas Lock(B) & Lock(A)                                |
              |                                             |
              +-------------------------------------------->+ (Lock A bebas)
                                                            |
                                                2. Goroutine 2 memperoleh Lock A
                                                3. Goroutine 2 mengunci Lock B
                                                4. Selesai tanpa Deadlock!
```

---

## Diagram 3 — Siklus Penanganan Deadlock dengan Application Retry

Diagram alur penanganan error deadlock pada fungsi `TransferWithRetry`.

```text
               +---------------------------+
               |   Mulai Eksekusi Transfer |
               +---------------------------+
                             |
                             v
               +---------------------------+
               |   attemptCtx dibatasi     |
               |       (10ms timeout)      |
               +---------------------------+
                             |
                             v
               +---------------------------+
               |  Eksekusi TransferNaive   |
               +---------------------------+
                             |
                   +---------+---------+
                   | Apakah error nil? |
                   +---------+---------+
                   |                   |
               YA  |                   | TIDAK
                   v                   v
      +------------------+    +------------------------------+
      | Transfer Sukses  |    | Apakah ErrDeadlock / Timeout |
      | (Kembalikan nil) |    +--------------+---------------+
      +------------------+                   |
                                   YA        |       TIDAK
                                             v         |
                             +-------------------+     v
                             | Iterasi < Max?    |  +--------------------+
                             +---------+---------+  | Kembalikan Error   |
                             |         |            | Asli (Bukan victim)|
                         YA  |         | TIDAK      +--------------------+
                             v         v
                   +-----------+   +-------------------+
                   | Backoff   |   | Kembalikan        |
                   | (Sleep)   |   | ErrDeadlock       |
                   +-----+-----+   +-------------------+
                         |
                         +---> (Kembali ke Mulai Eksekusi)
```

---

## Diagram 4 — Arsitektur Modul Simulasi Lab 18

Struktur paket internal yang merepresentasikan interaksi lapisan aplikasi dan entitas data.

```text
+--------------------------------------------------------------+
|                     cmd/demo/main.go                         |
|   Eksekusi demonstrasi skenario konkurensi end-to-end        |
+--------------------------------------------------------------+
                               |
                               | Memanggil strategi
                               v
+--------------------------------------------------------------+
|                internal/transfer/transfer.go                 |
|  - TransferNaive(ctx, from, to, amount, delay)               |
|  - TransferOrdered(ctx, acc1, acc2, amount, delay)           |
|  - TransferWithRetry(ctx, from, to, amount, delay, maxRetry) |
+--------------------------------------------------------------+
                               |
                               | Sinkronisasi channel & timeout
                               v
+--------------------------------------------------------------+
|                  internal/bank/account.go                    |
|  - Account { ID, Balance, ch chan struct{} }                 |
|  - Lock(ctx) -> nil | ErrDeadlock                            |
|  - Unlock()                                                  |
+--------------------------------------------------------------+
```
