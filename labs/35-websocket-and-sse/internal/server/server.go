package server

import (
	"fmt"
	"labs/35-websocket-and-sse/internal/sse"
	"labs/35-websocket-and-sse/internal/ws"
	"net/http"
)

type Server struct {
	SSEHub *sse.Hub
	Mux    *http.ServeMux
}

func NewServer() *Server {
	hub := sse.NewHub()
	mux := http.NewServeMux()

	s := &Server{
		SSEHub: hub,
		Mux:    mux,
	}

	mux.Handle("/sse", hub)

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := ws.Upgrade(w, r)
		if err != nil {
			return
		}
		defer conn.Close()

		for {
			op, data, err := conn.ReadFrame()
			if err != nil {
				return
			}
			switch op {
			case ws.OpClose:
				_ = conn.WriteFrame(ws.OpClose, nil, false)
				return
			case ws.OpPing:
				_ = conn.WriteFrame(ws.OpPong, data, false)
			case ws.OpText:
				// Echo text
				_ = conn.WriteFrame(ws.OpText, append([]byte("echo: "), data...), false)
			case ws.OpBinary:
				// Echo binary (reversed payload)
				rev := make([]byte, len(data))
				for i := range data {
					rev[i] = data[len(data)-1-i]
				}
				_ = conn.WriteFrame(ws.OpBinary, rev, false)
			}
		}
	})

	mux.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("msg")
		if msg == "" {
			msg = "default notification"
		}
		evt := hub.Broadcast("notification", msg)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Published event ID %d\n", evt.ID)
	})

	return s
}
