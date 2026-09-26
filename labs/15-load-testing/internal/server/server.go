package server

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sync/atomic"
	"time"
)

type Config struct {
	MaxDBConnections int
	DBQueryDuration  time.Duration
}

type Server struct {
	cfg       Config
	semaphore chan struct{}
	activeReq int64
}

// ActiveConnections returns the number of currently acquired database slots.
func (s *Server) ActiveConnections() int {
	return len(s.semaphore)
}

type BookingRequest struct {
	VehicleID string `json:"vehicle_id"`
	Service   string `json:"service"`
}

type BookingResponse struct {
	Status    string `json:"status"`
	BookingID string `json:"booking_id"`
}

func New(cfg Config) *Server {
	if cfg.MaxDBConnections <= 0 {
		cfg.MaxDBConnections = 5
	}
	if cfg.DBQueryDuration <= 0 {
		cfg.DBQueryDuration = 10 * time.Millisecond
	}
	return &Server{
		cfg:       cfg,
		semaphore: make(chan struct{}, cfg.MaxDBConnections),
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/booking", s.handleBooking)
	return mux
}

func (s *Server) handleBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	atomic.AddInt64(&s.activeReq, 1)
	defer atomic.AddInt64(&s.activeReq, -1)

	// Acquire DB connection slot (simulates DB connection pool limit)
	select {
	case s.semaphore <- struct{}{}:
	case <-r.Context().Done():
		return
	}
	defer func() { <-s.semaphore }()

	dur := s.cfg.DBQueryDuration
	if atomic.LoadInt64(&s.activeReq) > int64(s.cfg.MaxDBConnections) {
		if rand.Float32() < 0.10 {
			dur = s.cfg.DBQueryDuration * 25
		}
	}
	t := time.NewTimer(dur)
	defer t.Stop()

	select {
	case <-t.C:
	case <-r.Context().Done():
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(BookingResponse{
		Status:    "confirmed",
		BookingID: "BK-1001",
	})
}
