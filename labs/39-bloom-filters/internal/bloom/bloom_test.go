package bloom

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestNoFalseNegatives(t *testing.T) {
	const n = 10000
	const eps = 0.01

	bf := New(n, eps)
	keys := make([]string, n)
	for i := 0; i < n; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
		bf.Add([]byte(keys[i]))
	}

	for _, k := range keys {
		if !bf.Check([]byte(k)) {
			t.Fatalf("False negative detected for key: %s", k)
		}
	}
}

func TestFalsePositiveRate(t *testing.T) {
	const n = 10000
	const eps = 0.01
	const queries = 50000

	bf := New(n, eps)
	for i := 0; i < n; i++ {
		bf.Add([]byte(fmt.Sprintf("inserted-key-%d", i)))
	}

	falsePositives := 0
	for i := 0; i < queries; i++ {
		query := fmt.Sprintf("absent-key-%d", i)
		if bf.Check([]byte(query)) {
			falsePositives++
		}
	}

	measuredRate := float64(falsePositives) / float64(queries)
	t.Logf("Expected FP rate: %.4f, Measured FP rate: %.4f (%d / %d)", eps, measuredRate, falsePositives, queries)

	// Allow empirical rate up to 2.5x the theoretical rate due to random noise / hash characteristics
	if measuredRate > eps*2.5 {
		t.Fatalf("Measured FP rate %.4f exceeds tolerance (%.4f)", measuredRate, eps*2.5)
	}
}

func TestOptimalSizing(t *testing.T) {
	const n = 1000
	const eps = 0.01

	bf := New(n, eps)

	// Theoretical m = ceil(-1000 * ln(0.01) / (ln 2)^2) = ceil(1000 * 4.60517 / 0.48045) = 9586 bits
	expectedM := uint(math.Ceil(-float64(n) * math.Log(eps) / (math.Log(2) * math.Log(2))))
	expectedK := uint(math.Round(float64(expectedM) / float64(n) * math.Log(2)))

	if bf.M() != expectedM {
		t.Errorf("Expected m=%d, got %d", expectedM, bf.M())
	}
	if bf.K() != expectedK {
		t.Errorf("Expected k=%d, got %d", expectedK, bf.K())
	}

	// Bits per element should be ~9.585 for eps=0.01
	bitsPerElem := float64(bf.M()) / float64(n)
	if bitsPerElem < 9.5 || bitsPerElem > 10.5 {
		t.Errorf("Expected ~9.6 bits per element, got %.2f", bitsPerElem)
	}
}

func TestEmptyFilter(t *testing.T) {
	bf := New(100, 0.01)
	if bf.Check([]byte("any-key")) {
		t.Errorf("Empty filter should not return true for arbitrary key")
	}
}

func TestSyncFilterConcurrency(t *testing.T) {
	const n = 5000
	sf := NewSync(n, 0.01)

	var wg sync.WaitGroup
	const writers = 8
	const readers = 16

	for w := 0; w < writers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < n/writers; i++ {
				sf.Add([]byte(fmt.Sprintf("w%d-key-%d", workerID, i)))
			}
		}(w)
	}

	for r := 0; r < readers; r++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				sf.Check([]byte(fmt.Sprintf("r%d-query-%d", workerID, i)))
			}
		}(r)
	}

	wg.Wait()
}

func BenchmarkFilterAdd(b *testing.B) {
	bf := New(uint(b.N), 0.01)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Add([]byte(fmt.Sprintf("key-%d", i)))
	}
}

func BenchmarkFilterCheck(b *testing.B) {
	bf := New(10000, 0.01)
	for i := 0; i < 10000; i++ {
		bf.Add([]byte(fmt.Sprintf("key-%d", i)))
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Check([]byte(fmt.Sprintf("query-%d", rng.Intn(20000))))
	}
}
