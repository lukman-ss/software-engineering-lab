## Snippet 1 — Data Models and DTOs

Source File: internal/compat/model.go
Purpose: Definisikan struktur data untuk merepresentasikan schema legacy dan modern, beserta DTO untuk konsumen API V1 dan V2.

```go
type User struct {
    ID        int
    Name      string
    Phone     *string // Legacy field; nil once contracted
    CreatedAt time.Time
}

// UserResponse is an enriched additive response supporting both legacy and modern clients
type UserResponse struct {
    ID     int          `json:"id"`
    Name   string       `json:"name"`
    Phone  string       `json:"phone,omitempty"` // Legacy field maintained for v1 consumers
    Phones []PhoneEntry `json:"phones"`          // New field for v2 consumers
}

// LegacyConsumerDTO models an un-upgraded client expecting only string phone
type LegacyConsumerDTO struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Phone string `json:"phone"`
}

// ModernConsumerDTO models an upgraded client expecting phones array
type ModernConsumerDTO struct {
    ID     int          `json:"id"`
    Name   string       `json:"name"`
    Phones []PhoneEntry `json:"phones"`
}
```

Explanation:
- Struct `User` menggunakan pointer `*string` untuk field `Phone` sehingga dapat di-set ke `nil` setelah kontrak diterapkan (simulasi drop kolom).
- Struct `UserResponse` memberikan response yang "di-enrich" dengan baik field lama (`phone`) maupun field baru (`phones`) untuk mendukung kompatibilitas mundur dan maju dalam satu payload.
- DTO terpisah (`LegacyConsumerDTO` dan `ModernConsumerDTO`) menunjukkan kontrak eksplisit yang diharapkan oleh masing-masing versi klien, mempermudah validasi pada level marshalling/unmarshalling JSON.

## Snippet 2 — Dual-Write Storage Implementation

Source File: internal/compat/store.go
Purpose: Tampilkan mekanisme dual-write atomic yang menulis ke kedua tabel legacy dan modern dalam satu kritial section.

```go
func (s *MemoryStore) CreateDual(name, phone string) (*User, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.userSeq++
    id := s.userSeq

    var phonePtr *string
    if !s.legacyDropped {
        p := phone
        phonePtr = &p
    }

    u := &User{
        ID:        id,
        Name:      name,
        Phone:     phonePtr,
        CreatedAt: time.Now(),
    }
    s.users[id] = u

    // Write to modern user_phones table
    s.phoneSeq++
    entry := PhoneEntry{
        ID:        s.phoneSeq,
        UserID:    id,
        Number:    phone,
        IsPrimary: true,
    }
    s.userPhones[id] = append(s.userPhones[id], entry)

    return u, nil
}
```

Explanation:
- Fungsi ini hanya dipanggil ketika `WriteMode` diatur ke `WriteDual`.
- `sync.RWMutex` digunakan untuk menjamin atomicity penulisan ke kedua tabel dalam satu operasi kritial.
- Tulis ke tabel legacy terjadi jika `!s.legacyDropped` (kolom belum di-drop).
- Tulis ke tabel modern terjadi tanpa kondisi (selalu terjadi selama dual-write aktif).
- Jika sistem crash di tengah eksekusi, baik keduanya akan diperbarui atau keduanya tidak berubah, menjamin konsistensi data.

## Snippet 3 — Backfill Worker with Idempotent and Resumable Logic

Source File: internal/compat/backfill.go
Purpose: Tampilkan worker backfill yang dapat dihentikan dan dilanjutkan lagi dengan memanfaatkan checkpoint dan logika idempotent.

```go
func (b *BackfillWorker) RunBatch(ctx context.Context) (int, bool, error) {
    b.checkpoint.mu.Lock()
    defer b.checkpoint.mu.Unlock()

    if b.checkpoint.IsComplete {
        return 0, true, nil
    }

    ids := b.store.GetUserIDs(b.checkpoint.LastProcessedID, b.batchSize)
    if len(ids) == 0 {
        b.checkpoint.IsComplete = true
        return 0, true, nil
    }

    migratedInBatch := 0
    for _, id := range ids {
        select {
        case <-ctx.Done():
            return migratedInBatch, false, ctx.Err()
        default:
        }

        user, err := b.store.GetUser(id)
        if err != nil {
            continue
        }

        // Only backfill if legacy phone exists and user_phones is empty
        if user.Phone != nil && *user.Phone != "" {
            phones, _ := b.store.GetPhones(id)
            if len(phones) == 0 {
                _, err := b.store.SavePhoneEntry(id, *user.Phone, true)
                if err != nil {
                    return migratedInBatch, false, fmt.Errorf("failed backfilling user %d: %w", id, err)
                }
                migratedInBatch++
                b.obs.BackfillProcessed.Add(1)
            }
        }
        b.checkpoint.LastProcessedID = id
    }

    b.checkpoint.TotalMigrated += migratedInBatch

    if len(ids) < b.batchSize {
        b.checkpoint.IsComplete = true
    }

    return migratedInBatch, b.checkpoint.IsComplete, nil
}
```

Explanation:
- Checkpoint (`LastProcessedID` dan `IsComplete`) menyimpan state sehingga worker dapat dilanjutkan dari titik berhenti jika dihentikan.
- Loop memproses record dalam batch yang ditentukan oleh `batchSize`.
- Idempotency dicapai dengan memeriksa `len(phones) == 0` sebelum melakukan insert—jika sudah ada entri di `user_phones`, maka tidak dilakukan apa-apa.
- Setiap backfill yang berhasil meningkatkan counter observasi `BackfillProcessed`.
- Jika jumlah ID yang diterima kurang dari `batchSize`, maka diasumsikan bahwa seluruh data telah diproses dan `IsComplete` disetel ke true.

## Snippet 4 — Fallback Read (Dual Read) Mechanism

Source File: internal/compat/service.go
Purpose: Tampilkan mekanisme bacaan yang coba baca dari skema baru terlebih dahulu, lalu fallback ke skema lama jika diperlukan.

```go
func (s *Service) GetUser(id int) (UserResponse, error) {
    u, err := s.store.GetUser(id)
    if err != nil {
        return UserResponse{}, err
    }

    phones, _ := s.store.GetPhones(id)
    var phoneVal string

    readMode := s.flags.GetReadMode()

    switch readMode {
    case ReadLegacyOnly:
        if u.Phone != nil {
            phoneVal = *u.Phone
        }
    case ReadFallback:
        if len(phones) > 0 {
            // Find primary or first phone
            phoneVal = phones[0].Number
            for _, p := range phones {
                if p.IsPrimary {
                    phoneVal = p.Number
                    break
                }
            }
        } else if u.Phone != nil {
            // Fallback read from legacy column
            phoneVal = *u.Phone
            // Lazy backfill: write to new schema during fallback read
            _, _ = s.store.SavePhoneEntry(id, phoneVal, true)
        }
    case ReadNewOnly:
        if len(phones) > 0 {
            phoneVal = phones[0].Number
            for _, p := range phones {
                if p.IsPrimary {
                    phoneVal = p.Number
                    break
                }
            }
        }
    }

    return UserResponse{
        ID:     u.ID,
        Name:   u.Name,
        Phone:  phoneVal,
        Phones: phones,
    }, nil
}
```

Explanation:
- Pada mode `ReadFallback`, sistem pertama-tama mencoba membaca dari `user_phones` (skema baru).
- Jika hasilnya kosong dan data legacy masih tersedia (`u.Phone != nil`), maka nilai dari kolom legacy digunakan.
- Setelah fallback read berhasil, sistem melakukan "lazy backfill" dengan menulis nilai tersebut ke `user_phones` untuk mengurangi beban bacaan di masa depan.
- Mekanisme ini mencegah "data starvation" saat backfill belum selesai namun aplikasi sudah mulai membaca dari skema baru.

## Snippet 5 — Contract Enforcement with Traffic Guard

Source File: internal/compat/service.go
Purpose: Tampilkan bagaimana kontrak hanya dapat diterapkan ketika tidak ada trafic legacy terdeteksi.

```go
func (s *Service) ApplyContract(force bool) error {
    // Guard: Ensure legacy read hits are zero since last reset, unless forced
    if !force && s.obs.LegacyReadHits.Load() > 0 {
        return fmt.Errorf("%w: recorded %d legacy reads", ErrContractViolation, s.obs.LegacyReadHits.Load())
    }

    // 1. Switch write mode to NewOnly
    s.flags.SetWriteMode(WriteNewOnly)
    // 2. Switch read mode to NewOnly
    s.flags.SetReadMode(ReadNewOnly)
    // 3. Mark contract applied
    s.flags.SetContractApplied(true)
    // 4. Drop legacy column from storage
    s.store.ApplyContractDropLegacyColumn()

    return nil
}
```

Explanation:
- Fungsi ini memeriksa meter `LegacyReadHits` sebelum mengizinkan kontrak.
- Jika tidak dipaksa (`force == false`) dan masih ada baca legacy yang terdeteksi, maka fungsi mengembalikan error `ErrContractViolation`.
- Jika izin diberikan, maka:
  1. Mode tulis diatur ke `WriteNewOnly` (tidak lagi menulis ke legacy)
  2. Mode baca diatur ke `ReadNewOnly` (tidak lagi membaca dari legacy)
  3. Flag kontrak diatur ke `true`
  4. Operasi drop kolom legacy dipanggil pada storage (`ApplyContractDropLegacyColumn`)
- Pendekatan ini memastikan bahwa kontrak hanya diterapkan ketika kita yakin tidak ada klien lama yang masih aktif.

## Snippet 6 — HTTP Handler with Deprecation Headers

Source File: internal/compat/handler.go
Purpose: Tampilkan implementasi header Deprecation dan Sunset sesuai RFC 8594 untuk memberi tahu klien tentang pengangkatan.

```go
func (h *APIHandler) GetUserV1(w http.ResponseWriter, r *http.Request) {
    idStr := r.URL.Query().Get("id")
    id, _ := strconv.Atoi(idStr)

    dto, err := h.svc.GetLegacyUser(id)
    if err != nil {
        if errors.Is(err, ErrLegacyUnavailable) {
            w.WriteHeader(http.StatusGone) // 410 Gone after sunset
            w.Write([]byte(`{"error": "legacy v1 endpoint has been permanently removed"}`))
            return
        }
        w.WriteHeader(http.StatusNotFound)
        return
    }

    // Inject Deprecation Headers
    w.Header().Set("Deprecation", "true")
    w.Header().Set("Sunset", "Mon, 31 Dec 2026 23:59:59 GMT") // Sunset date
    w.Header().Set("Content-Type", "application/json")

    json.NewEncoder(w).Encode(dto)
}
```

Explanation:
- Header `Deprecation: true` memberitahu klien bahwa endpoint ini tidak lagi disarankan untuk digunakan.
- Header `Sunset: <tanggal>` memberitahu klien tanggal tepat ketika endpoint akan dihapus kembali (setelah tanggal ini, klien akan menerima HTTP 410 Gone).
- Implementasi ini mengikuti standar RFC 8594 yang digunakan oleh platform seperti GitHub API untuk pengangkatan versi API.
- Setelah kontrak diterapkan dan legacy endpoint diakses, klien akan menerima respons `410 Gone` dengan pesan error bahwa endpoint telah dihapus secara permanen.

## Snippet 7 — Feature Flags for Runtime Control

Source File: internal/compat/flags.go
Purpose: Tampilkan implementasi fitur flag yang memungkinkan kontrol runtime atas mode tulis dan baca tanpa perlu redeploy.

```go
type FeatureFlags struct {
    writeMode       atomic.Value // WriteMode
    readMode        atomic.Value // ReadMode
    contractApplied atomic.Bool  // true if legacy column/endpoints dropped
}

func NewFeatureFlags() *FeatureFlags {
    f := &FeatureFlags{}
    f.writeMode.Store(WriteLegacyOnly)
    f.readMode.Store(ReadLegacyOnly)
    return f
}

func (f *FeatureFlags) GetWriteMode() WriteMode {
    return f.writeMode.Load().(WriteMode)
}

func (f *FeatureFlags) SetWriteMode(m WriteMode) {
    f.writeMode.Store(m)
}

func (f *FeatureFlags) GetReadMode() ReadMode {
    return f.readMode.Load().(ReadMode)
}

func (f *FeatureFlags) SetReadMode(m ReadMode) {
    f.readMode.Store(m)
}

func (f *FeatureFlags) IsContractApplied() bool {
    return f.contractApplied.Load()
}

func (f *FeatureFlags) SetContractApplied(b bool) {
    f.contractApplied.Store(b)
}
```

Explanation:
- Menggunakan `atomic.Value` dan `atomic.Bool` untuk memastikan akses thread-safe ke nilai flag tanpa mutex eksplisit.
- Tiga mode yang dikoordinasikan:
  - `WriteMode`: Mengontrol kemana data baru ditulis (`WriteLegacyOnly`, `WriteDual`, `WriteNewOnly`)
  - `ReadMode`: Mengontrol darimana data dibaca (`ReadLegacyOnly`, `ReadFallback`, `ReadNewOnly`)
  - `ContractApplied`: Bendera boolean yang menunjukkan apakah kontrak (drop kolom legacy) telah diterapkan
- Fitur flag ini memungkinkan:
  - Canary rollout: mengalihkan sebagian kecil traffic ke mode baru terlebih dahulu
  - Instant rollback: kembali ke mode legacy dengan satu operasi jika diperlukan
  - Pengendalian fase migrasi: masing-masing fase (Expand, Migrate, Contract) dapat dikontrol independen melalui kombinasi flag yang berbeda