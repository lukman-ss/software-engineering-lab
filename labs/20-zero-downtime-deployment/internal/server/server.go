package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Server struct {
	srv         *http.Server
	ready       atomic.Bool
	activeCount atomic.Int32
	wg          sync.WaitGroup
	preStop     time.Duration
}

func NewServer(addr string, preStopDelay time.Duration) *Server {
	mux := http.NewServeMux()
	s := &Server{
		preStop: preStopDelay,
	}

	mux.HandleFunc("/healthz/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/healthz/ready", func(w http.ResponseWriter, r *http.Request) {
		if s.ready.Load() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("READY"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("NOT READY"))
		}
	})

	mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		s.activeCount.Add(1)
		s.wg.Add(1)
		defer s.wg.Done()
		defer s.activeCount.Add(-1)

		durationStr := r.URL.Query().Get("d")
		d, err := time.ParseDuration(durationStr)
		if err != nil {
			d = 50 * time.Millisecond
		}

		select {
		case <-time.After(d):
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("WORK COMPLETED"))
		case <-r.Context().Done():
			// Client disconnected early
			return
		}
	})

	s.srv = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return s
}

func (s *Server) Start() error {
	log.Printf("Server starting on %s", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) SetReady(ready bool) {
	s.ready.Store(ready)
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Server received shutdown request")

	s.SetReady(false)
	log.Println("Server marked unready, detached from load balancer")

	if s.preStop > 0 {
		log.Printf("Executing preStop sleep for %v to allow routing table updates...", s.preStop)
		time.Sleep(s.preStop)
	}

	log.Println("Initiating graceful shutdown of HTTP listeners...")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown error: %w", err)
	}

	s.wg.Wait()
	log.Println("All in-flight requests completed. Server stopped gracefully.")
	return nil
}

func (s *Server) ActiveRequests() int32 {
	return s.activeCount.Load()
}
