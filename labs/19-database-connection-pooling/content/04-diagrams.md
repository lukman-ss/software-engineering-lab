## Connection Pooling Flow

```text
[Aplikasi Konkuren]
    │      │      │
(Banyak Goroutines / Threads)
    ▼      ▼      ▼
[ Client Pool (Limit) ] ---> [Timeout] jika pool kosong
    │      │
(Re-use Connections)
    ▼      ▼
[ Database Backend ] 
(Max Connections Limit)
```

## Leak / Starvation Flow

```text
Goroutine 1: 
Minta Koneksi -> [Pool Kosong? Tidak] -> Berhasil
  -> Eksekusi DB
  -> Panggil API Eksternal LAAAAAMAAAAA (Koneksi tertahan)

Goroutine 2:
Minta Koneksi -> [Pool Kosong? Ya] -> Menunggu 
  -> Timeout (Context Deadline Exceeded)
```