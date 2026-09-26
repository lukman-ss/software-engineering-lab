package tests

import (
	"testing"
	"time"

	"zero-downtime-deployment/internal/worker"
)

func TestWorkerGracefulShutdown(t *testing.T) {
	w := worker.NewWorker(10)
	w.Start(1)

	w.Enqueue(worker.Job{ID: "job-1", Duration: 50 * time.Millisecond})
	w.Enqueue(worker.Job{ID: "job-2", Duration: 50 * time.Millisecond})

	time.Sleep(10 * time.Millisecond)

	w.Stop()

	completed := w.GetCompletedJobs()
	if len(completed) < 1 {
		t.Fatalf("expected at least 1 job to complete gracefully, got %d", len(completed))
	}
}
