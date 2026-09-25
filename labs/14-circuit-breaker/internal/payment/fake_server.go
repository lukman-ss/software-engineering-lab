package payment

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"
)

type ServerMode string

const (
	ModeHealthy ServerMode = "HEALTHY"
	ModeSlow    ServerMode = "SLOW"
	ModeDown    ServerMode = "DOWN"
)

type FakeServer struct {
	server       *httptest.Server
	mode         atomic.Value
	requestCount atomic.Int64
	slowDelay    time.Duration
}

func NewFakeServer(slowDelay time.Duration) *FakeServer {
	fs := &FakeServer{
		slowDelay: slowDelay,
	}
	fs.mode.Store(ModeHealthy)

	mux := http.NewServeMux()
	mux.HandleFunc("/pay", func(w http.ResponseWriter, r *http.Request) {
		fs.requestCount.Add(1)
		mode := fs.mode.Load().(ServerMode)

		switch mode {
		case ModeSlow:
			time.Sleep(fs.slowDelay)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok_slow"}`))
		case ModeDown:
			http.Error(w, "internal payment server failure", http.StatusInternalServerError)
		default: // ModeHealthy
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
		}
	})

	fs.server = httptest.NewServer(mux)
	return fs
}

func (fs *FakeServer) URL() string {
	return fs.server.URL
}

func (fs *FakeServer) SetMode(m ServerMode) {
	fs.mode.Store(m)
}

func (fs *FakeServer) RequestCount() int64 {
	return fs.requestCount.Load()
}

func (fs *FakeServer) ResetCount() {
	fs.requestCount.Store(0)
}

func (fs *FakeServer) Close() {
	fs.server.Close()
}
