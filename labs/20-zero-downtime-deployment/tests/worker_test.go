package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"zero-downtime-deployment/internal/worker"
)

func TestWorkerConcurrency(t *testing.T) {
	w := worker.NewWorker(10)
	w.Start(3)

	for i := 0; i < 6; i++ {
		w.Enqueue(worker.Job{ID: fmt.Sprintf("job-%d", i), Duration: 20 * time.Millisecond})
	}

	w.Stop(500 * time.Millisecond)

	completed := w.GetCompletedJobs()
	if len(completed) != 6 {
		t.Fatalf("expected all 6 jobs to complete with concurrency=3, got %d", len(completed))
	}
}

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

func TestWorkerEnqueueAfterStop(t *testing.T) {
	w := worker.NewWorker(10)
	w.Start(1)
	w.Stop(100 * time.Millisecond)

	w.Enqueue(worker.Job{ID: "post-stop-job", Duration: 0})

	completed := w.GetCompletedJobs()
	if len(completed) != 0 {
		t.Fatalf("expected 0 completed jobs after Stop, got %d: %v", len(completed), completed)
	}
}

func TestWorkerConcurrentEnqueueStop(t *testing.T) {
	for i := 0; i < 50; i++ {
		w := worker.NewWorker(100)
		w.Start(2)

		var wg sync.WaitGroup
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				w.Enqueue(worker.Job{ID: fmt.Sprintf("job-%d-%d", i, id), Duration: 0})
			}(j)
		}

		w.Stop(200 * time.Millisecond)
		wg.Wait()
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
