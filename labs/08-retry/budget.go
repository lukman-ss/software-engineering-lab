package retry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// RetryBudget limits the total number of retries allowed in a time window.
// This prevents retry amplification during prolonged outages.
type RetryBudget struct {
	mu            sync.Mutex
	budget        int64
	window        time.Duration
	used          int64
	lastReplenish time.Time
}

// NewRetryBudget creates a new retry budget.
func NewRetryBudget(budget int64, window time.Duration) *RetryBudget {
	return &RetryBudget{
		budget:        budget,
		window:        window,
		lastReplenish: time.Now(),
	}
}

// TryConsume returns true if a retry is allowed by the budget.
func (rb *RetryBudget) TryConsume() bool {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	now := time.Now()
	if now.Sub(rb.lastReplenish) >= rb.window {
		// Replenish budget for new window
		rb.used = 0
		rb.lastReplenish = now
	}

	if rb.used < rb.budget {
		rb.used++
		return true
	}
	return false
}

// BudgetedClient wraps RetryableClient with budget enforcement.
type BudgetedClient struct {
	*RetryableClient
	budget *RetryBudget
}

// NewBudgetedClient creates a client that enforces a retry budget.
func NewBudgetedClient(client *RetryableClient, budget *RetryBudget) *BudgetedClient {
	return &BudgetedClient{
		RetryableClient: client,
		budget:          budget,
	}
}

// Do overrides the base Do method to check budget before retrying.
func (c *BudgetedClient) Do(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetry; attempt++ {
		// For attempt > 0 (actual retries), check budget first
		if attempt > 0 {
			if !c.budget.TryConsume() {
				return nil, fmt.Errorf("retry budget exhausted, aborting. last error: %w", lastErr)
			}
		}

		c.totalAttempts.Add(1)

		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, err
		}

		resp, err := c.client.Do(req)
		if err != nil {
			if attempt < c.maxRetry && isRetryableNetworkError(err) {
				waitDuration(c.baseDelay, attempt, c.jitterFactor)
				lastErr = err
				continue
			}
			return nil, err
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}

		if !IsTransientError(resp.StatusCode) {
			return resp, fmt.Errorf("permanent error: status %d", resp.StatusCode)
		}

		if attempt < c.maxRetry {
			waitDuration(c.baseDelay, attempt, c.jitterFactor)
			lastErr = fmt.Errorf("retryable error: status %d", resp.StatusCode)
			continue
		}

		return resp, lastErr
	}

	return nil, lastErr
}

// Get is a convenience method for GET requests.
func (c *BudgetedClient) Get(ctx context.Context, url string) (*http.Response, error) {
	return c.Do(ctx, http.MethodGet, url, nil)
}

// DemonstrateAmplification explains the retry amplification problem.
func DemonstrateAmplification() {
	fmt.Println("RETRY AMPLIFICATION FACTOR")
	fmt.Println("==========================")
	fmt.Println("")
	fmt.Println("The danger of uncoordinated retries across multiple service layers.")
	fmt.Println("")
	fmt.Println("Architecture: Client -> API Gateway -> Service A -> Service B -> DB")
	fmt.Println("Assuming DB goes down.")
	fmt.Println("")
	fmt.Println("Scenario 1: No Retry Budgets (Multiplicative Amplification)")
	fmt.Println("  - Service A retries Service B 3 times")
	fmt.Println("  - API Gateway retries Service A 3 times")
	fmt.Println("  - Client retries API Gateway 3 times")
	fmt.Println("")
	fmt.Println("  1 initial request from client becomes:")
	fmt.Println("    = 1 * 3 (Client -> Gateway)")
	fmt.Println("    = 3 * 3 (Gateway -> Svc A)")
	fmt.Println("    = 9 * 3 (Svc A -> Svc B)")
	fmt.Println("    = 27 requests hitting Service B!")
	fmt.Println("")
	fmt.Println("  If 100 concurrent users do this:")
	fmt.Println("  100 requests -> 2,700 requests hitting the failing backend.")
	fmt.Println("  When DB recovers, it is instantly crushed by the amplified load (Retry Storm).")
	fmt.Println("")
	fmt.Println("Scenario 2: With Retry Budgets")
	fmt.Println("  - Each layer allows max 10% of total traffic as retries")
	fmt.Println("  - Client: 100 reqs -> 10 retries allowed")
	fmt.Println("  - Gateway: receives 110 reqs -> 11 retries allowed")
	fmt.Println("  - Service A: receives 121 reqs -> 12 retries allowed")
	fmt.Println("  - Service B receives: ~133 requests total (not 2,700!)")
	fmt.Println("")
	fmt.Println("RECOMMENDATIONS:")
	fmt.Println("  1. Retry at only ONE layer (usually the layer closest to the failure or the client edge)")
	fmt.Println("  2. If multi-layer retries are necessary, use Retry Budgets (e.g., token bucket)")
	fmt.Println("  3. Pass \"retry-attempt: N\" headers so downstream services know not to retry a retry")
	fmt.Println("  4. Fail fast when budget is exhausted")
}

// AmplificationFactor demonstrates the math.
func AmplificationFactor(layers, retriesPerLayer int) int {
	result := 1
	for i := 0; i < layers; i++ {
		result *= (1 + retriesPerLayer)
	}
	return result
}
