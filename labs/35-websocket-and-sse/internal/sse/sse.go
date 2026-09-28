package sse

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Event struct {
	ID    int
	Event string
	Data  string
	Retry int
}

func (e Event) Format() []byte {
	var sb strings.Builder
	if e.ID > 0 {
		fmt.Fprintf(&sb, "id: %d\n", e.ID)
	}
	if e.Event != "" {
		fmt.Fprintf(&sb, "event: %s\n", e.Event)
	}
	if e.Retry > 0 {
		fmt.Fprintf(&sb, "retry: %d\n", e.Retry)
	}
	lines := strings.Split(e.Data, "\n")
	for _, l := range lines {
		fmt.Fprintf(&sb, "data: %s\n", l)
	}
	sb.WriteString("\n")
	return []byte(sb.String())
}

type Hub struct {
	mu      sync.RWMutex
	history []Event
	clients map[chan Event]struct{}
	nextID  int
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[chan Event]struct{}),
		nextID:  1,
	}
}

func (h *Hub) Broadcast(eventType, data string) Event {
	h.mu.Lock()
	defer h.mu.Unlock()

	evt := Event{
		ID:    h.nextID,
		Event: eventType,
		Data:  data,
		Retry: 2000,
	}
	h.nextID++
	h.history = append(h.history, evt)

	for ch := range h.clients {
		select {
		case ch <- evt:
		default:
		}
	}
	return evt
}

func (h *Hub) Subscribe(lastID int) (chan Event, []Event, func()) {
	h.mu.Lock()
	defer h.mu.Unlock()

	var replay []Event
	for _, evt := range h.history {
		if evt.ID > lastID {
			replay = append(replay, evt)
		}
	}

	ch := make(chan Event, 64)
	h.clients[ch] = struct{}{}

	unsubscribe := func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		delete(h.clients, ch)
		close(ch)
	}

	return ch, replay, unsubscribe
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	lastID := 0
	if idHeader := r.Header.Get("Last-Event-ID"); idHeader != "" {
		if id, err := strconv.Atoi(idHeader); err == nil {
			lastID = id
		}
	}

	ch, replay, unsub := h.Subscribe(lastID)
	defer unsub()

	for _, evt := range replay {
		if _, err := w.Write(evt.Format()); err != nil {
			return
		}
	}
	flusher.Flush()

	ctx := r.Context()
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-ch:
			if !ok {
				return
			}
			if _, err := w.Write(evt.Format()); err != nil {
				return
			}
			flusher.Flush()
		case <-ticker.C:
			// SSE comment as keep-alive heartbeat
			if _, err := w.Write([]byte(": keep-alive\n\n")); err != nil {
				return
			}
			flusher.Flush()
		}
	}
}
