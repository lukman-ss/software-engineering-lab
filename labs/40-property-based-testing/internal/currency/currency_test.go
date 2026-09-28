package currency

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

// Example-based test: passes hand-picked cases, giving false confidence.
func TestExampleBasedNaiveCurrency(t *testing.T) {
	naive := NaiveCurrency{}
	examples := []float64{10.50, 99.99, 1.00, 5.25}

	for _, ex := range examples {
		formatted := naive.Format(ex)
		parsed, err := naive.Parse(formatted)
		if err != nil {
			t.Fatalf("failed parsing %f: %v", ex, err)
		}
		if parsed != ex {
			t.Fatalf("expected %f, got %f", ex, parsed)
		}
	}
}

type NaiveFloat float64

func (NaiveFloat) Generate(r *rand.Rand, size int) reflect.Value {
	// Generate floating point with sub-cent precision, negatives, and normal values
	dollars := (r.Float64() - 0.5) * 10000
	return reflect.ValueOf(NaiveFloat(dollars))
}

// Property-based test demonstrating failure on naive float currency.
func TestPropertyNaiveCurrencyFails(t *testing.T) {
	naive := NaiveCurrency{}
	prop := func(val NaiveFloat) bool {
		dollars := float64(val)
		if math.IsNaN(dollars) || math.IsInf(dollars, 0) {
			return true
		}
		formatted := naive.Format(dollars)
		parsed, err := naive.Parse(formatted)
		if err != nil {
			return false
		}
		// Subcent precision is lost in %.2f format
		return parsed == dollars
	}

	cfg := &quick.Config{
		MaxCount: 100,
	}

	err := quick.Check(prop, cfg)
	if err == nil {
		t.Fatal("expected naive float currency property check to fail, but it passed")
	}
}

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
