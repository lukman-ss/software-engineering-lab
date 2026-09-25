# Code Snippets

## Snippet 1 — Schema Evolution (Expand vs Contract)

Source File: `schema.sql`
Purpose: Mendefinisikan schema relasional awal (V1), perluasan struktur 1:N (V2 - Expand), dan penghapusan kolom lama (V3 - Contract).

```sql
-- V1: Baseline (1:1 relationship)
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone TEXT NOT NULL
);

-- V2: Expand (1:N relationship added without touching old column)
CREATE TABLE user_phones (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    number TEXT NOT NULL,
    is_primary BOOLEAN DEFAULT FALSE
);
CREATE INDEX idx_user_phones_user_id ON user_phones(user_id);

-- V3: Contract (executed ONLY after legacy consumer and dual-write are fully phased out)
-- ALTER TABLE users DROP COLUMN phone;
```

Explanation:
Tahap Expand menambahkan tabel `user_phones` tanpa mengubah atau menghapus kolom `phone` pada tabel `users`. Dengan cara ini, query versi lama tetap dapat berjalan.

---

## Snippet 2 — Additive Response Model

Source File: `internal/compat/model.go`
Purpose: Menyediakan DTO respons yang aditif agar kompatibel dengan consumer V1 maupun V2.

```go
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
Struct `UserResponse` menyertakan representasi lama (`Phone`) dan representasi baru (`Phones`), memungkinkan kedua jenis consumer membaca data dari endpoint yang sama tanpa deserialization error.

---

## Snippet 3 — Atomic Dual-Write Logic

Source File: `internal/compat/store.go`
Purpose: Menulis data baru ke representasi lama dan baru secara simultan.

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
Operasi `CreateDual` menulis data ke map `users` dan `userPhones` dalam satu proteksi lock, memastikan kedua tabel tetap sinkron saat data baru masuk.
