# Code Snippets

Dokumen ini memuat potongan kode yang bersumber langsung dari implementasi terverifikasi pada lab `labs/17-architecture-decision-record`.

---

## Snippet 1 — Model Status dan Entitas Record

Source File: `internal/adr/models.go:3-29`
Purpose: Mendefinisikan status siklus hidup ADR dan struktur data penampung metadata keputusan arsitektur.

```go
type Status string

const (
	StatusProposed   Status = "Proposed"
	StatusAccepted   Status = "Accepted"
	StatusSuperseded Status = "Superseded"
	StatusDeprecated Status = "Deprecated"
	StatusRejected   Status = "Rejected"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusProposed, StatusAccepted, StatusSuperseded, StatusDeprecated, StatusRejected:
		return true
	default:
		return false
	}
}

type Record struct {
	ID           int
	Title        string
	Status       Status
	SupersededBy int // 0 if not superseded
	Supersedes   int // 0 if does not supersede
	Content      string
}
```

Explanation:
Mendefinisikan 5 status siklus hidup formal ADR sesuai panduan *AWS Prescriptive Guidance*. Metode `IsValid()` memvalidasi masukan status secara ketat, sementara struct `Record` menyimpan relasi silsilah penggantian (`SupersededBy` dan `Supersedes`).

---

## Snippet 2 — Parsing Metadata ADR Berbasis Regex

Source File: `internal/adr/parser.go:17-83`
Purpose: Memindai teks berkas Markdown untuk mengekstrak ID, Judul, Status, dan relasi silsilah keputusan.

```go
func Parse(content string) (*Record, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	record := &Record{
		Content: content,
	}

	foundTitle := false
	foundStatus := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if !foundTitle {
			matches := titleRegex.FindStringSubmatch(line)
			if len(matches) == 3 {
				id, err := strconv.Atoi(matches[1])
				if err != nil {
					return nil, fmt.Errorf("invalid record id: %w", err)
				}
				record.ID = id
				record.Title = matches[2]
				foundTitle = true
				continue
			}
		}

		if strings.HasPrefix(strings.ToLower(line), "status:") {
			matches := statusRegex.FindStringSubmatch(line)
			if len(matches) >= 2 {
				statusStr := strings.Title(strings.ToLower(matches[1]))
				record.Status = Status(statusStr)
				foundStatus = true

				if len(matches) == 3 && matches[2] != "" {
					supersededBy, err := strconv.Atoi(matches[2])
					if err == nil {
						record.SupersededBy = supersededBy
					}
				}
			}
		}

		if strings.HasPrefix(strings.ToLower(line), "supersedes:") {
			matches := supersedesRegex.FindStringSubmatch(line)
			if len(matches) == 2 {
				supersedes, err := strconv.Atoi(matches[1])
				if err == nil {
					record.Supersedes = supersedes
				}
			}
		}
	}

	if !foundTitle {
		return nil, fmt.Errorf("title not found or malformed")
	}

	if !foundStatus {
		return nil, fmt.Errorf("status not found")
	}

	if !record.Status.IsValid() {
		return nil, fmt.Errorf("invalid status: %s", record.Status)
	}

	return record, nil
}
```

Explanation:
Fungsi `Parse` membaca dokumen baris demi baris menggunakan `bufio.Scanner`. Atribut judul dan ID diekstraksi dari heading `# 1. Title`, status dan target pengganti diekstraksi dari baris `Status: Superseded by 2`, dan referensi ke keputusan sebelumnya diekstraksi dari `Supersedes: 1`.

---

## Snippet 3 — Validasi Penomoran Monotonik

Source File: `internal/adr/linter.go:27-39`
Purpose: Memastikan urutan penomoran seluruh ADR bernilai sekuensial tanpa loncatan nomor.

```go
	// Validate IDs are monotonic
	ids := make([]int, 0, len(records))
	for id := range recordMap {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	for i := 0; i < len(ids); i++ {
		if ids[i] != i+1 {
			errs = append(errs, fmt.Errorf("non-monotonic numbering, expected %d but got %d", i+1, ids[i]))
			break // Only report once
		}
	}
```

Explanation:
Linter mengumpulkan seluruh ID ADR yang terdaftar, mengurutkannya dari terkecil ke terbesar, lalu memverifikasi bahwa indeks ke-$i$ bernilai tepat $i+1$. Jika terdapat nomor yang terlewat (misalnya ADR 1 lalu ADR 3), proses validasi menghasilkan galat `non-monotonic numbering`.

---

## Snippet 4 — Validasi Konkuren Relasi Timbal-Balik Graf Keputusan

Source File: `internal/adr/linter.go:41-79`
Purpose: Memvalidasi integritas dua arah silsilah keputusan arsitektur secara paralel menggunakan Goroutines.

```go
	// Validate graph in parallel
	var wg sync.WaitGroup
	for _, r := range records {
		wg.Add(1)
		go func(rec *Record) {
			defer wg.Done()
			var localErrs []error

			if rec.Status == StatusSuperseded {
				if rec.SupersededBy == 0 {
					localErrs = append(localErrs, fmt.Errorf("ADR %d is superseded but missing superseded_by reference", rec.ID))
				} else {
					replacement, exists := recordMap[rec.SupersededBy]
					if !exists {
						localErrs = append(localErrs, fmt.Errorf("ADR %d superseded by non-existent ADR %d", rec.ID, rec.SupersededBy))
					} else if replacement.Supersedes != rec.ID {
						localErrs = append(localErrs, fmt.Errorf("ADR %d superseded by ADR %d, but ADR %d does not declare it supersedes ADR %d", rec.ID, rec.SupersededBy, rec.SupersededBy, rec.ID))
					}
				}
			}

			if rec.Supersedes != 0 {
				old, exists := recordMap[rec.Supersedes]
				if !exists {
					localErrs = append(localErrs, fmt.Errorf("ADR %d supersedes non-existent ADR %d", rec.ID, rec.Supersedes))
				} else if old.Status != StatusSuperseded || old.SupersededBy != rec.ID {
					localErrs = append(localErrs, fmt.Errorf("ADR %d supersedes ADR %d, but ADR %d is not properly marked as superseded by ADR %d", rec.ID, rec.Supersedes, rec.Supersedes, rec.ID))
				}
			}

			if len(localErrs) > 0 {
				mu.Lock()
				errs = append(errs, localErrs...)
				mu.Unlock()
			}
		}(r)
	}

	wg.Wait()
```

Explanation:
Setiap record divalidasi dalam goroutine terpisah. Jika record berstatus `Superseded`, ia memverifikasi bahwa record penggantinya ada dan menyatakan `Supersedes` secara resiprokal. Sebaliknya, jika record menyatakan `Supersedes`, record lama harus berstatus `Superseded` dengan penunjuk ke ID baru. Akses penggabungan ke slice galat `errs` dilindungi oleh `sync.Mutex`.

---

## Snippet 5 — Pengujian Kegagalan Referensi Penggantian

Source File: `tests/linter_test.go:36-100`
Purpose: Menguji deteksi kegagalan struktural referensi rusak atau tautan sepihak.

```go
func TestLinter_BrokenReferences(t *testing.T) {
	tests := []struct {
		name    string
		records []*adr.Record
		wantErr string
	}{
		{
			name: "superseded points to non-existent",
			records: []*adr.Record{
				{
					ID:           1,
					Title:        "Modular Monolith",
					Status:       adr.StatusSuperseded,
					SupersededBy: 99,
				},
			},
			wantErr: "superseded by non-existent ADR 99",
		},
		{
			name: "supersedes points to non-existent",
			records: []*adr.Record{
				{
					ID:         1,
					Title:      "Microservices",
					Status:     adr.StatusAccepted,
					Supersedes: 99,
				},
			},
			wantErr: "supersedes non-existent ADR 99",
		},
		{
			name: "mismatched supersession link",
			records: []*adr.Record{
				{
					ID:           1,
					Title:        "Modular Monolith",
					Status:       adr.StatusSuperseded,
					SupersededBy: 2,
				},
				{
					ID:     2,
					Title:  "Microservices",
					Status: adr.StatusAccepted,
					// Forgot Supersedes: 1
				},
			},
			wantErr: "ADR 1 superseded by ADR 2, but ADR 2 does not declare it supersedes ADR 1",
		},
		{
			name: "non-monotonic numbering",
			records: []*adr.Record{
				{
					ID:     1,
					Title:  "One",
					Status: adr.StatusAccepted,
				},
				{
					ID:     3,
					Title:  "Three",
					Status: adr.StatusAccepted,
				},
			},
			wantErr: "non-monotonic numbering, expected 2 but got 3",
		},
	}
```

Explanation:
Uji tabel (*table-driven test*) yang memvalidasi bahwa linter secara akurat menghasilkan pesan galat yang sesuai pada empat skenario kegagalan struktural: target pengganti tidak ada, asal yang digantikan tidak ada, tautan sepihak yang tidak sinkron, dan penomoran tidak monotonik.
