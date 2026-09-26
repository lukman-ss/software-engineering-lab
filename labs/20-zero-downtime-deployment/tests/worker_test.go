package tests

import (
	"testing"
	"time"

	"zero-downtime-deployment/internal/worker"
)

func TestWorkerGracefulShutdown(t *testing.T) {
	w := worker.NewWorker(10)
	w.Start(1)

	w.Enqueue(worker.Job{ID: "job-1", Duration: 20 * time.Millisecond})
	w.Enqueue(worker.Job{ID: "job-2", Duration: 20 * time.Millisecond})

	time.Sleep(5 * time.Millisecond)

	w.Stop(100 * time.Millisecond)

	completed := w.GetCompletedJobs()
	if len(completed) != 2 {
		t.Fatalf("expected all 2 jobs to complete gracefully during drain, got %d", len(completed))
	}
	if completed[0] != "job-1" || completed[1] != "job-2" {
		t.Fatalf("unexpected completed jobs: %v", completed)
	}
}

func TestWorkerShutdownTimeout(t *testing.T) {
	w := worker.NewWorker(10)
	w.Start(1)

	w.Enqueue(worker.Job{ID: "job-slow", Duration: 100 * time.Millisecond})
	w.Enqueue(worker.Job{ID: "job-dropped", Duration: 100 * time.Millisecond})

	time.Sleep(5 * time.Millisecond)

	// Stop with short timeout so job-dropped is abandoned
	w.Stop(20 * time.Millisecond)

	completed := w.GetCompletedJobs()
	if len(completed) != 1 {
		t.Fatalf("expected 1 job (job-slow) completed before timeout, got %d", len(completed))
	}
	if completed[0] != "job-slow" {
		t.Fatalf("expected job-slow to be completed, got %v", completed)
	}
}
