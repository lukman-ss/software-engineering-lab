## Snippet 1 — Domain Models (Legacy & Modern DTOs)

**Source File:** `internal/compat/model.go:1-39`

**Purpose:** Define structured types for legacy (1:1) and modern (1:N) phone representations. The `UserResponse` uses additive JSON tags to serve both client versions from a single endpoint.

```go
type PhoneEntry struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Number    string `json:"number"`
	IsPrimary bool   `json:"is_primary"`
}

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

**Explanation:** The `omitempty` on `Phone` ensures legacy field is omitted when empty (post-contract), while `Phones` array is always present. This single payload shape satisfies both client generations without versioned endpoints for reads.

---

## Snippet 2 — In-Memory Store with Dual-Write and Contract Drop

**Source File:** `internal/compat/store.go:53-85, 204-212`

**Purpose:** Simulate relational tables with atomic dual-write during Expand/Migrate and legacy column drop during Contract.

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

func (s *MemoryStore) ApplyContractDropLegacyColumn() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.legacyDropped = true
	for _, u := range s.users {
		u.Phone = nil
	}
}
```

**Explanation:** `CreateDual` atomically writes to both `users.phone` (unless already contracted) and `user_phones` table under a single mutex, simulating a transactional dual-write. `ApplyContractDropLegacyColumn` simulates `ALTER TABLE ... DROP COLUMN` by clearing the legacy pointer and setting a guard flag.

---

## Snippet 3 — Feature Flags Controlling Migration Phases

**Source File:** `internal/compat/flags.go:7-54`

**Purpose:** Atomic feature flags that gate write/read modes and contract state, enabling canary rollout and instant rollback.

```go
type WriteMode int
type ReadMode int

const (
	WriteLegacyOnly WriteMode = iota
	WriteDual
	WriteNewOnly
)

const (
	ReadLegacyOnly ReadMode = iota
	ReadFallback
	ReadNewOnly
)

type FeatureFlags struct {
	writeMode       atomic.Value // WriteMode
	readMode        atomic.Value // ReadMode
	contractApplied atomic.Bool  // true if legacy column/endpoints dropped
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

**Explanation:** Using `atomic.Value` and `atomic.Bool` ensures lock-free, thread-safe flag reads during concurrent request processing. Phases map directly: `WriteLegacyOnly` (Baseline), `WriteDual` (Expand/Migrate), `WriteNewOnly` (Post-Contract).

---

## Snippet 4 — Fallback Read with Lazy Backfill

**Source File:** `internal/compat/service.go:108-140`

**Purpose:** Read path that serves data from legacy column when modern table is empty, and lazily backfills during fallback reads.

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

**Explanation:** In `ReadFallback` mode, if `user_phones` is empty but `users.phone` exists, the legacy value is returned AND lazily written to `user_phones` via `SavePhoneEntry`. This ensures reads succeed during backfill without requiring full historical migration upfront.

---

## Snippet 5 — Idempotent Resumable Backfill Worker

**Source File:** `internal/compat/backfill.go:35-101`

**Purpose:** Batch-process legacy records with checkpoint persistence and idempotent insert logic.

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

func (b *BackfillWorker) RunAll(ctx context.Context) (int, error) {
	total := 0
	for {
		migrated, done, err := b.RunBatch(ctx)
		if err != nil {
			return total, err
		}
		total += migrated
		if done {
			break
		}
	}
	return total, nil
}
```

**Explanation:** Checkpoint uses `LastProcessedID` cursor for resumability. Idempotency achieved by checking `len(phones) == 0` before insert and `SavePhoneEntry` checking for existing number (see `store.go:155-160`). Batch size configurable for throttling.

---

## Snippet 6 — Data Reconciliation for Drift Detection

**Source File:** `internal/compat/service.go:186-226`

**Purpose:** Compare legacy `users.phone` against `user_phones` primary entry to detect dual-write desynchronization.

```go
func (s *Service) ReconcileData() (int, error) {
	if s.flags.IsContractApplied() {
		return 0, nil
	}

	ids := s.store.GetUserIDs(0, s.store.TotalUsers()+10)
	drifts := 0

	for _, id := range ids {
		user, err := s.store.GetUser(id)
		if err != nil {
			continue
		}
		if user.Phone == nil {
			continue
		}

		phones, err := s.store.GetPhones(id)
		if err != nil || len(phones) == 0 {
			drifts++
			s.obs.DriftDetected.Add(1)
			continue
		}

		var primaryNumber string
		for _, p := range phones {
			if p.IsPrimary {
				primaryNumber = p.Number
				break
			}
		}

		if primaryNumber != *user.Phone {
			drifts++
			s.obs.DriftDetected.Add(1)
		}
	}

	return drifts, nil
}
```

**Explanation:** Iterates all users, compares legacy `phone` pointer against primary entry in `user_phones`. Increments drift counter for missing entries or value mismatches. Guarded by `IsContractApplied()` to skip after legacy column dropped.

---

## Snippet 7 — Contract Enforcement with Zero-Traffic Guard

**Source File:** `internal/compat/service.go:228-245`

**Purpose:** Apply Contract phase only when observability confirms zero legacy traffic (unless forced).

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

**Explanation:** The `force` parameter allows bypassing the zero-traffic check (used in demo/tests to simulate post-observation-window). In production, `force=false` enforces the safety gate. Sequence: stop dual-write → stop legacy reads → mark contract → drop column.

---

## Snippet 8 — HTTP Deprecation and Sunset Headers

**Source File:** `internal/compat/handler.go:20-41`

**Purpose:** RFC 8594-compliant headers signaling deprecation to legacy consumers.

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

**Explanation:** Legacy endpoint `/api/v1/users` returns `Deprecation: true` and `Sunset` date headers. After contract, returns `410 Gone`. Modern endpoint `/api/v2/users` (lines 44-55) has no deprecation headers.

---

## Snippet 9 — Observability Metrics Collector

**Source File:** `internal/compat/metrics.go:7-28`

**Purpose:** Thread-safe atomic counters for migration monitoring.

```go
type Observability struct {
	LegacyReadHits    atomic.Int64
	NewReadHits       atomic.Int64
	DualWriteCount    atomic.Int64
	DualWriteErrors   atomic.Int64
	BackfillProcessed atomic.Int64
	DriftDetected     atomic.Int64
}

func (o *Observability) Snapshot() map[string]int64 {
	return map[string]int64{
		"legacy_reads":      o.LegacyReadHits.Load(),
		"new_reads":         o.NewReadHits.Load(),
		"dual_writes":       o.DualWriteCount.Load(),
		"dual_write_errors": o.DualWriteErrors.Load(),
		"backfilled":        o.BackfillProcessed.Load(),
		"drift_detected":    o.DriftDetected.Load(),
	}
}
```

**Explanation:** Uses `sync/atomic` for lock-free incrementing under concurrent load. `Snapshot()` provides a point-in-time view for dashboards and contract gates.