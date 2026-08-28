// Package retry demonstrates robust retry mechanisms for handling transient failures.
// This lab covers:
// - Mock provider with controllable failure modes
// - Exponential backoff with jitter
// - Retry storm prevention
// - Retry budget limiting
package retry

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

// HTTPStatus represents an HTTP status code.
type HTTPStatus int

const (
	StatusOK                  HTTPStatus = 200
	StatusBadRequest          HTTPStatus = 400
	StatusUnauthorized        HTTPStatus = 401
	StatusForbidden           HTTPStatus = 403
	StatusNotFound            HTTPStatus = 404
	StatusTooManyRequests     HTTPStatus = 429
	StatusInternalServerError HTTPStatus = 500
	StatusServiceUnavailable  HTTPStatus = 503
)

// MockProvider simulates an external API provider with controllable failure modes.
// It can be configured to:
// - Fail N times before succeeding
// - Return 500 (internal server error)
// - Return 429 (rate limit)
// - Timeout
// - Succeed
type MockProvider struct {
	config      mockConfig
	failedCount atomic.Int32
	totalReqs   atomic.Int32
}

type mockConfig struct {
	maxFailures   int
	statusCode    HTTPStatus
	delay         time.Duration
	shouldTimeout bool
	failCount     atomic.Int32
}

// NewMockProvider creates a new mock provider with specified behavior.
func NewMockProvider() *MockProvider {
	return &MockProvider{
		config: mockConfig{
			statusCode:    StatusOK,
			delay:         0,
			shouldTimeout: false,
		},
	}
}

// WithMaxFailures makes the provider fail N times before succeeding.
func (p *MockProvider) WithMaxFailures(n int) *MockProvider {
	p.config.maxFailures = n
	return p
}

// WithStatusCode sets the HTTP status code to return.
func (p *MockProvider) WithStatusCode(code HTTPStatus) *MockProvider {
	p.config.statusCode = code
	return p
}

// WithDelay adds latency to each request.
func (p *MockProvider) WithDelay(d time.Duration) *MockProvider {
	p.config.delay = d
	return p
}

// WithTimeout causes requests to timeout.
func (p *MockProvider) WithTimeout(timeout bool) *MockProvider {
	p.config.shouldTimeout = timeout
	return p
}

// Do implements the HTTP client interface.
func (p *MockProvider) Do(req *http.Request) (*http.Response, error) {
	p.totalReqs.Add(1)

	// Simulate timeout
	if p.config.shouldTimeout {
		return nil, errors.New("connection timeout")
	}

	// Simulate latency
	if p.config.delay > 0 {
		select {
		case <-time.After(p.config.delay):
		case <-req.Context().Done():
			return nil, req.Context().Err()
		}
	}

	// Check failure count
	failures := p.config.failCount.Load()
	if failures < int32(p.config.maxFailures) {
		p.config.failCount.Add(1)
		p.failedCount.Add(1)
		body, _ := json.Marshal(map[string]string{"error": "simulated failure"})
		return &http.Response{
			StatusCode: int(p.config.statusCode),
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	}

	// Success
	body, _ := json.Marshal(map[string]string{"status": "ok"})
	return &http.Response{
		StatusCode: int(StatusOK),
		Body:       io.NopCloser(bytes.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

// Get makes a GET request.
func (p *MockProvider) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return p.Do(req)
}

// Reset resets the failure counter.
func (p *MockProvider) Reset() {
	p.config.failCount.Store(0)
	p.failedCount.Store(0)
	p.totalReqs.Store(0)
}

// FailedCount returns the number of failed requests.
func (p *MockProvider) FailedCount() int32 {
	return p.failedCount.Load()
}

// TotalRequests returns the total number of requests.
func (p *MockProvider) TotalRequests() int32 {
	return p.totalReqs.Load()
}

// ErrPermanent is a permanent error that should not be retried.
type ErrPermanent struct {
	Msg string
}

func (e *ErrPermanent) Error() string { return e.Msg }

// IsPermanent returns true if the error should never be retried.
func IsPermanent(err error) bool {
	var e *ErrPermanent
	return errors.As(err, &e)
}

// IsTransientError returns true if the error should be retried based on HTTP status.
func IsTransientError(statusCode int) bool {
	switch HTTPStatus(statusCode) {
	case StatusBadRequest:
		// 400 could be bad request or validation error - typically not retried
		return false
	case StatusUnauthorized:
		// 401 - credentials issue, not transient
		return false
	case StatusForbidden:
		// 403 - permission issue, not transient
		return false
	case StatusNotFound:
		// 404 - resource doesn't exist
		return false
	case StatusTooManyRequests:
		// 429 - rate limited, should retry with backoff
		return true
	case StatusInternalServerError:
		// 500 - server error, typically transient
		return true
	case StatusServiceUnavailable:
		// 503 - service unavailable, transient
		return true
	default:
		return statusCode >= 500
	}
}

// ShouldRetry determines if an error should be retried.
func ShouldRetry(err error, statusCode int) bool {
	if IsPermanent(err) {
		return false
	}
	return IsTransientError(statusCode)
}

// Provider represents an HTTP client provider for making requests.
type Provider interface {
	Do(req *http.Request) (*http.Response, error)
}

// RetryableClient implements exponential backoff with jitter.
type RetryableClient struct {
	client        Provider
	maxRetry      int
	baseDelay     time.Duration
	jitterFactor  float64
	totalAttempts atomic.Int32
}

// RetryableClientOption configures a RetryableClient.
type RetryableClientOption func(*RetryableClient)

// WithMaxRetry sets the maximum number of retry attempts.
func WithRetryAttempts(n int) RetryableClientOption {
	return func(c *RetryableClient) {
		c.maxRetry = n
	}
}

// WithBaseDelay sets the base delay for exponential backoff.
func WithBaseDelay(d time.Duration) RetryableClientOption {
	return func(c *RetryableClient) {
		c.baseDelay = d
	}
}

// WithJitter sets the jitter factor (0.0 to 1.0).
func WithJitter(factor float64) RetryableClientOption {
	return func(c *RetryableClient) {
		c.jitterFactor = factor
	}
}

// NewRetryableClient creates a new retry client.
func NewRetryableClient(provider Provider, opts ...RetryableClientOption) *RetryableClient {
	c := &RetryableClient{
		client:       provider,
		maxRetry:     3,
		baseDelay:    100 * time.Millisecond,
		jitterFactor: 0.5,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Do executes the request with retry logic.
func (c *RetryableClient) Do(ctx context.Context, method, url string, body io.Reader) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetry; attempt++ {
		c.totalAttempts.Add(1)

		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, err
		}

		resp, err := c.client.Do(req)
		if err != nil {
			if attempt < c.maxRetry && isRetryableNetworkError(err) {
				waitDuration(c.baseDelay, attempt, c.jitterFactor)
				continue
			}
			return nil, err
		}

		// Check response status
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}

		// Non-retryable status codes
		if !IsTransientError(resp.StatusCode) {
			return resp, fmt.Errorf("permanent error: status %d", resp.StatusCode)
		}

		// Retryable status codes
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
func (c *RetryableClient) Get(ctx context.Context, url string) (*http.Response, error) {
	return c.Do(ctx, http.MethodGet, url, nil)
}

// TotalAttempts returns the total number of attempts made.
func (c *RetryableClient) TotalAttempts() int32 {
	return c.totalAttempts.Load()
}

// ResetAttempts resets the attempt counter.
func (c *RetryableClient) ResetAttempts() {
	c.totalAttempts.Store(0)
}

// isRetryableNetworkError checks if a network error should be retried.
func isRetryableNetworkError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	return true
}

// waitDuration implements exponential backoff with full jitter.
func waitDuration(baseDelay time.Duration, attempt int, jitterFactor float64) time.Duration {
	// Exponential backoff: base * 2^attempt
	delay := float64(baseDelay) * float64(1<<uint(attempt))

	// Add jitter: multiply by random factor between (1-jitter) and (1+jitter)
	jitter := 1.0 + jitterFactor*(rand.Float64()*2-1)
	return time.Duration(float64(delay) * jitter)
}
