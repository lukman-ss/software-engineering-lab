package worker

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type Job struct {
	ID       string
	Duration time.Duration
}

// ponytail: in-memory Go channel queue ceiling; upgrade to Redis/RabbitMQ/Kafka when distributed worker pools required.
// ponytail: cooperative drain ceiling; in-flight active jobs run to completion while queued jobs are abandoned upon drain timeout. Upgrade to job-level context preemption or hard deadline escalation when arbitrary job runtimes require forced interruption.
type Worker struct {
	jobChan     chan Job
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	completed   []string
	completedMu sync.Mutex
	stopped     atomic.Bool
}

func NewWorker(bufferSize int) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		jobChan: make(chan Job, bufferSize),
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (w *Worker) Start(concurrency int) {
	for i := 0; i < concurrency; i++ {
		w.wg.Add(1)
		go func(workerID int) {
			defer w.wg.Done()
			for {
				select {
				case <-w.ctx.Done():
					return
				default:
				}

				select {
				case <-w.ctx.Done():
					return
				case job, ok := <-w.jobChan:
					if !ok {
						return
					}
					log.Printf("Worker %d starting job %s", workerID, job.ID)
					time.Sleep(job.Duration)
					w.completedMu.Lock()
					w.completed = append(w.completed, job.ID)
					w.completedMu.Unlock()
					log.Printf("Worker %d finished job %s", workerID, job.ID)
				}
			}
		}(i)
	}
}

func (w *Worker) Enqueue(job Job) {
	if w.stopped.Load() {
		log.Printf("Enqueue rejected: worker stopped, dropping job %s", job.ID)
		return
	}
	w.jobChan <- job
}

func (w *Worker) Stop(timeout time.Duration) {
	log.Println("Worker receiving stop signal, no longer accepting new jobs...")
	w.stopped.Store(true)
	close(w.jobChan)
	
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		log.Println("Worker drain timeout reached, cancelling context...")
		w.cancel()
		<-done
	}
	log.Println("Worker gracefully stopped")
}

func (w *Worker) GetCompletedJobs() []string {
	w.completedMu.Lock()
	defer w.completedMu.Unlock()
	res := make([]string, len(w.completed))
	copy(res, w.completed)
	return res
}
