package worker

import (
	"context"
	"log"
	"sync"
	"time"
)

type Job struct {
	ID       string
	Duration time.Duration
}

type Worker struct {
	jobChan     chan Job
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	completed   []string
	completedMu sync.Mutex
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
	w.jobChan <- job
}

func (w *Worker) Stop() {
	log.Println("Worker receiving stop signal, no longer accepting new jobs...")
	close(w.jobChan)
	w.cancel()
	w.wg.Wait()
	log.Println("Worker gracefully stopped")
}

func (w *Worker) GetCompletedJobs() []string {
	w.completedMu.Lock()
	defer w.completedMu.Unlock()
	res := make([]string, len(w.completed))
	copy(res, w.completed)
	return res
}
