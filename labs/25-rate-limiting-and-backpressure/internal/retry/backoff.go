package retry

import (
	"math"
	"math/rand"
	"time"
)

type BackoffStrategy string

const (
	NoJitter           BackoffStrategy = "NoJitter"
	FullJitter         BackoffStrategy = "FullJitter"
	EqualJitter        BackoffStrategy = "EqualJitter"
	DecorrelatedJitter BackoffStrategy = "DecorrelatedJitter"
)

type Config struct {
	Base time.Duration
	Cap  time.Duration
}

// ComputeBackoff calculates sleep duration according to AWS Marc Brooker formulas.
func ComputeBackoff(strategy BackoffStrategy, attempt int, cfg Config, prevSleep time.Duration) time.Duration {
	baseFloat := float64(cfg.Base)
	capFloat := float64(cfg.Cap)

	// temp = min(cap, base * 2^attempt)
	expBackoff := baseFloat * math.Pow(2, float64(attempt))
	temp := math.Min(capFloat, expBackoff)

	switch strategy {
	case NoJitter:
		return time.Duration(temp)

	case FullJitter:
		// sleep = random_between(0, min(cap, base * 2^attempt))
		if temp <= 0 {
			return 0
		}
		sleep := rand.Float64() * temp
		return time.Duration(sleep)

	case EqualJitter:
		// sleep = min(cap, base * 2^attempt) / 2 + random_between(0, min(cap, base * 2^attempt) / 2)
		half := temp / 2.0
		sleep := half + rand.Float64()*half
		return time.Duration(sleep)

	case DecorrelatedJitter:
		// sleep = min(cap, random_between(base, prevSleep * 3))
		prevFloat := float64(prevSleep)
		if prevFloat < baseFloat {
			prevFloat = baseFloat
		}
		rangeMax := prevFloat * 3.0
		sleep := baseFloat + rand.Float64()*(rangeMax-baseFloat)
		return time.Duration(math.Min(capFloat, sleep))

	default:
		return time.Duration(temp)
	}
}
