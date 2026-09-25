package loadtest

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	URL         string
	Method      string
	Body        []byte
	ContentType string
	VUs         int
	Duration    time.Duration
}

type Runner struct {
	cfg    Config
	client *http.Client
}

func NewRunner(cfg Config) *Runner {
	if cfg.VUs <= 0 {
		cfg.VUs = 1
	}
	// Custom transport to prevent connection pooling limits from hiding server bottlenecks
	tr := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 1000,
	}
	return &Runner{
		cfg: cfg,
		client: &http.Client{
			Transport: tr,
			Timeout:   5 * time.Second,
		},
	}
}

func (r *Runner) Run(ctx context.Context) Result {
	var wg sync.WaitGroup
	start := time.Now()

	// Per-VU results to avoid lock contention
	type vuResult struct {
		latencies []time.Duration
		errors    int
	}
	results := make([]vuResult, r.cfg.VUs)

	ctx, cancel := context.WithTimeout(ctx, r.cfg.Duration)
	defer cancel()

	for i := 0; i < r.cfg.VUs; i++ {
		wg.Add(1)
		go func(vuID int) {
			defer wg.Done()
			var lats []time.Duration
			var errs int

			for {
				select {
				case <-ctx.Done():
					results[vuID] = vuResult{latencies: lats, errors: errs}
					return
				default:
					req, err := http.NewRequestWithContext(ctx, r.cfg.Method, r.cfg.URL, bytes.NewReader(r.cfg.Body))
					if err != nil {
						errs++
						continue
					}
					if r.cfg.ContentType != "" {
						req.Header.Set("Content-Type", r.cfg.ContentType)
					}

					reqStart := time.Now()
					resp, err := r.client.Do(req)
					if err != nil {
						// Only count as error if not a context cancellation
						if ctx.Err() == nil {
							errs++
						}
						continue
					}
					_ = resp.Body.Close()
					
					if resp.StatusCode >= 400 {
						errs++
					} else {
						lats = append(lats, time.Since(reqStart))
					}
				}
			}
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(start)

	var allLatencies []time.Duration
	var totalErrs int
	for _, res := range results {
		allLatencies = append(allLatencies, res.latencies...)
		totalErrs += res.errors
	}

	return CalculateMetrics(allLatencies, totalErrs, totalDuration)
}
