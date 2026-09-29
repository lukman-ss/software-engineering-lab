# Code Snippets: Property-Based Testing

Koleksi snippet kode yang diekstrak langsung dari implementasi dan test suite lab `labs/40-property-based-testing`.

---

## Snippet 1 — Robust Currency Struct & Formatter

Source File: `internal/currency/currency.go`  
Purpose: Mengelola representasi moneter menggunakan integer cents (`int64`) untuk mencegah kehilangan presisi biner (IEEE 754 floating point) serta memformatnya ke format `$Dollars.Cents` dengan dukungan nilai negatif.

```go
// RobustAmount stores currency as integer cents to preserve exact values.
type RobustAmount struct {
	Cents int64
}

func NewRobustAmount(cents int64) RobustAmount {
	return RobustAmount{Cents: cents}
}

func (r RobustAmount) Format() string {
	absCents := r.Cents
	sign := ""
	if absCents < 0 {
		sign = "-"
		absCents = -absCents
	}
	dollars := absCents / 100
	cents := absCents % 100
	return fmt.Sprintf("%s$%d.%02d", sign, dollars, cents)
}
```

Explanation:
Operasi pembagian integer `/ 100` dan modulo `% 100` menjamin nilai pecahan tetap tepat 2 desimal tanpa efek pembulatan biner tidak terduga.

---

## Snippet 2 — Custom Biased Generator & Roundtrip Property Test

Source File: `internal/currency/currency_test.go`  
Purpose: Menyiapkan generator acak berbasis boundary (`testing/quick.Generator`) dan memverifikasi Invariant Roundtrip `ParseRobust(amount.Format()) == amount` sebanyak 1.000 iterasi.

```go
// Custom generator for RobustAmount to exercise broad integer ranges including boundaries.
func (RobustAmount) Generate(r *rand.Rand, size int) reflect.Value {
	// Bias towards 0, negative values, large numbers, small cent values
	choice := r.Intn(5)
	var cents int64
	switch choice {
	case 0:
		cents = 0
	case 1:
		cents = int64(r.Intn(200) - 100) // -100 to 99 cents
	case 2:
		cents = int64(r.Intn(1000000))   // $0.00 to $10,000.00
	case 3:
		cents = -int64(r.Intn(1000000))  // -$0.00 to -$10,000.00
	default:
		cents = r.Int63n(1000000000000) - 500000000000
	}
	return reflect.ValueOf(RobustAmount{Cents: cents})
}

// Property-based test proving Roundtrip Invariant on RobustAmount:
// ParseRobust(amount.Format()) == amount
func TestPropertyRobustAmountRoundtrip(t *testing.T) {
	prop := func(a RobustAmount) bool {
		formatted := a.Format()
		parsed, err := ParseRobust(formatted)
		if err != nil {
			return false
		}
		return parsed.Cents == a.Cents
	}

	cfg := &quick.Config{
		MaxCount: 1000,
	}

	if err := quick.Check(prop, cfg); err != nil {
		t.Fatalf("property violation in RobustAmount roundtrip: %v", err)
	}
}
```

Explanation:
Generator sengaja mengarahkan distribusi nilai ke angka 0, sub-dollar negatif, dan bilangan multi-miliar cent untuk memvalidasi batas atas dan bawah sistem.

---

## Snippet 3 — Robust Interval Merging Implementation

Source File: `internal/interval/interval.go`  
Purpose: Menggabungkan irisan waktu/interval tertutup `[Start, End]` secara robust dengan melakukan pengurutan terlebih dahulu.

```go
// RobustMerge merges intervals correctly regardless of input order.
func RobustMerge(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}
	sorted := make([]Interval, len(intervals))
	copy(sorted, intervals)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Start != sorted[j].Start {
			return sorted[i].Start < sorted[j].Start
		}
		return sorted[i].End < sorted[j].End
	})

	out := []Interval{sorted[0]}
	for _, iv := range sorted[1:] {
		last := &out[len(out)-1]
		if iv.Start <= last.End+1 {
			if iv.End > last.End {
				last.End = iv.End
			}
		} else {
			out = append(out, iv)
		}
	}
	return out
}
```

Explanation:
Kloning slice dan sorting dua tingkat (`Start` lalu `End`) mencegah mutasi in-place pada input caller serta menangani overlap contiguous (`iv.Start <= last.End+1`).

---

## Snippet 4 — Counterexample Shrinking Algorithm

Source File: `internal/shrinker/shrinker.go`  
Purpose: Mereduksi slice counterexample yang gagal secara sistematis melalui binary sectioning, element removal, dan value reduction.

```go
// FindAndShrink searches for an input array that violates `predicate`,
// then performs counterexample shrinking (list binary sectioning & element reduction).
func FindAndShrink(predicate func([]int) bool, rng *rand.Rand, maxAttempts int) (*ShrinkResult, error) {
	// ... pencarian input awal yang gagal ...

	// Strategy 1: Try halving / chunk removal
	if len(current) > 1 {
		mid := len(current) / 2
		left := current[:mid]
		// Evaluasi apakah 'left' masih memicu kegagalan
		if !predicate(left) {
			current = append([]int(nil), left...)
			changed = true
			continue
		}
		// Evaluasi apakah 'right' masih memicu kegagalan
		right := current[mid:]
		if !predicate(right) {
			current = append([]int(nil), right...)
			changed = true
			continue
		}
	}

	// Strategy 2: Remove individual elements
	for i := 0; i < len(current); i++ {
		reduced := append(current[:i], current[i+1:]...)
		if !predicate(reduced) {
			current = reduced
			changed = true
			break
		}
	}

	// Strategy 3: Reduce element magnitude toward zero
	for i := 0; i < len(current); i++ {
		for _, cand := range []int{0, current[i] / 2} {
			// Evaluasi pengurangan nilai ke 0 atau setengahnya
		}
	}
	return res, nil
}
```

Explanation:
Algoritma shrinking mencoba strategi dari yang paling agresif (memotong 50% data) ke yang paling granular (menghapus satu elemen dan membagi nilai angka) hingga invariant tetap gagal pada ukuran minimum.
