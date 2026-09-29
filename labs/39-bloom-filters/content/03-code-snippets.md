## Snippet 1 — Inisialisasi Filter Optimal

Source File: internal/bloom/bloom.go:22-31
Purpose: Menghitung ukuran bit array dan jumlah hash berdasarkan target FP rate dan jumlah elemen.

```go
func New(n uint, epsilon float64) *Filter {
	m := optimalM(n, epsilon)
	k := optimalK(m, n)
	words := (m + 63) / 64
	return &Filter{
		bits: make([]uint64, words),
		m:    m,
		k:    k,
	}
}
```

Explanation: Konstruktor menggunakan rumus tertutup dari paper Bloom (1970) dan Kirsch‑Mitzenmacher (2006) untuk menentukan `m` dan `k` yang optimal. Fungsi `optimalM` dan `optimalK` mengimplementasikan rumus matematis yang dibuktikan.

---

## Snippet 2 — Double Hashing Kirsch‑Mitzenmacher

Source File: internal/bloom/bloom.go:79-83
Purpose: Menghasilkan k posisi bit dari dua hash dasar untuk mengurangi overhead komputasi.

```go
func doubleHash(h1, h2, i, m uint64) uint64 {
	return (h1 + i*h2) % m
}
```

Explanation: Teknik ini membuktikan bahwa penggunaan `h1 + i*h2 (mod m)` memberikan distribusi yang serupa dengan k hash independen tanpa memanggil fungsi hash berulang kali.

---

## Snippet 3 — Pembuatan Hash Dasar (FNV-1a + Rotate‑XOR)

Source File: internal/bloom/bloom.go:65-77
Purpose: Menghasilkan dua hash 64-bit independen dari kunci tanpa library eksternal.

```go
func baseHashes(key []byte) (uint64, uint64) {
	h := fnv.New64a()
	_, _ = h.Write(key)
	h1 := h.Sum64()
	h2 := h1 ^ (h1>>17 | h1<<47)
	if h2 == 0 {
		h2 = 0xdeadbeefdeadbeef
	}
	return h1, h2
}
```

Explanation: Menggunakan FNV‑1a dari stdlib untuk `h1`. `h2` dihasilkan melalui rotasi XOR bitwise untuk memastikan independensi statistik sementara menghindari `h2 == 0` yang akan menyebabkan semua probe mengembalikan posisi yang sama.

---

## Snippet 4 — Penambahan Elemen ke Filter

Source File: internal/bloom/bloom.go:34-40
Purpose: Menyet bit pada k posisi yang dihitung oleh double hashing.

```go
func (f *Filter) Add(key []byte) {
	h1, h2 := baseHashes(key)
	for i := uint(0); i < f.k; i++ {
		pos := doubleHash(h1, h2, uint64(i), uint64(f.m))
		f.bits[pos/64] |= 1 << (pos % 64)
	}
}
```

Explanation: Untuk setiap probe `i`, hitung posisi bit, lalu set bit tersebut menggunakan operasi OR pada word `uint64` yang sesuai.

---

## Snippet 5 — Pengecekan Keanggotaan

Source File: internal/bloom/bloom.go:44-53
Purpose: Mengembalikan false jika setidaknya satu bit 0; true jika semua bit 1.

```go
func (f *Filter) Check(key []byte) bool {
	h1, h2 := baseHashes(key)
	for i := uint(0); i < f.k; i++ {
		pos := doubleHash(h1, h2, uint64(i), uint64(f.m))
		if f.bits[pos/64]&(1<<(pos%64)) == 0 {
			return false
		}
	}
	return true
}
```

Explanation: Jika ada satu bit yang tidak di‑set, kunci pasti tidak pernah ditambahkan (zero false negative). Jika semua bit di‑set, kunci mungkin ada (kemungkinan false positive).

---

## Snippet 6 — Pruning Segment LSM‑Tree

Source File: internal/store/lsm.go:56-66
Purpose: Melewati segment bila filter menyatakan kunci tidak ada, sehingga tidak increment counter disk read.

```go
func (s *LSMStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keyBytes := []byte(key)
	for i := len(s.segments) - 1; i >= 0; i-- {
		seg := s.segments[i]
		if seg.filter != nil && !seg.filter.Check(keyBytes) {
			continue
		}
		atomic.AddUint64(&s.diskReads, 1)
		if v, ok := seg.data[key]; ok {
			return v, true
		}
	}
	return "", false
}
```

Explanation: Pengecekan Bloom filter dilakukan sebelum simulasi disk read. Jika filter mengembalikan false, segment dilewati sepenuhnya.

---

## Snippet 7 — Cache dengan Admission Gate Bloom Filter

Source File: internal/store/cache.go:54-58
Purpose: Memblokir kunci absent sebelum mencapai backend.

```go
func (c *Cache) Get(key string) (string, bool) {
	if c.filter != nil && !c.filter.Check([]byte(key)) {
		return "", false
	}
	// ... lanjut ke in‑memory map, lalu backend jika miss
}
```

Explanation: Filter di‑check terlebih dahulu. Jika mengembalikan false, kunci tidak dapat ada di cache atau backend, sehingga request ditolak segera tanpa memanggil backend.

---

## Snippet 8 — Wrapper Thread‑Safe SyncFilter

Source File: internal/bloom/bloom.go:104-128
Purpose: Menyediakan akses concurrent dengan RWMutex.

```go
type SyncFilter struct {
	mu sync.RWMutex
	f  *Filter
}

func (sf *SyncFilter) Add(key []byte) {
	sf.mu.Lock()
	sf.f.Add(key)
	sf.mu.Unlock()
}

func (sf *SyncFilter) Check(key []byte) bool {
	sf.mu.RLock()
	ok := sf.f.Check(key)
	sf.mu.RUnlock()
	return ok
}
```

Explanation: `Add` menggunakan exclusive lock; `Check` menggunakan shared read lock untuk memaksimalkan throughput pembaca.