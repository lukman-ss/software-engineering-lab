package tests

import (
	"bufio"
	"bytes"
	"labs/35-websocket-and-sse/internal/server"
	"labs/35-websocket-and-sse/internal/sse"
	"labs/35-websocket-and-sse/internal/ws"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestSSE_Formatting(t *testing.T) {
	evt := sse.Event{
		ID:    42,
		Event: "update",
		Data:  "hello\nworld",
		Retry: 1000,
	}
	expected := "id: 42\nevent: update\nretry: 1000\ndata: hello\ndata: world\n\n"
	if string(evt.Format()) != expected {
		t.Fatalf("expected:\n%q\ngot:\n%q", expected, string(evt.Format()))
	}
}

func TestSSE_LastEventID_Resumption(t *testing.T) {
	s := server.NewServer()
	ts := httptest.NewServer(s.Mux)
	defer ts.Close()

	// Broadcast 3 events
	s.SSEHub.Broadcast("msg", "first")
	s.SSEHub.Broadcast("msg", "second")
	s.SSEHub.Broadcast("msg", "third")

	req, _ := http.NewRequest("GET", ts.URL+"/sse", nil)
	req.Header.Set("Last-Event-ID", "1") // Should skip event 1 and replay 2 and 3

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var lines []string
	for i := 0; i < 8; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		trimmed := strings.TrimRight(line, "\r\n")
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}

	foundID2 := false
	foundID3 := false
	for _, l := range lines {
		if l == "id: 2" {
			foundID2 = true
		}
		if l == "id: 3" {
			foundID3 = true
		}
		if l == "id: 1" {
			t.Errorf("unexpected replayed event ID 1")
		}
	}

	if !foundID2 || !foundID3 {
		t.Fatalf("failed to resume from Last-Event-ID: 1, got lines: %v", lines)
	}
}

func TestWebSocket_TextAndBinary(t *testing.T) {
	s := server.NewServer()
	ts := httptest.NewServer(s.Mux)
	defer ts.Close()

	addr := ts.Listener.Addr().String()
	rawConn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("net dial failed: %v", err)
	}
	defer rawConn.Close()

	wsConn, err := ws.Dial("/ws", rawConn)
	if err != nil {
		t.Fatalf("ws dial failed: %v", err)
	}
	defer wsConn.Close()

	// 1. Text test
	sendText := "ping test"
	if err := wsConn.WriteFrame(ws.OpText, []byte(sendText), true); err != nil {
		t.Fatalf("write text frame failed: %v", err)
	}

	op, payload, err := wsConn.ReadFrame()
	if err != nil {
		t.Fatalf("read text frame failed: %v", err)
	}
	if op != ws.OpText {
		t.Fatalf("expected text opcode 0x1, got 0x%X", op)
	}
	if string(payload) != "echo: ping test" {
		t.Fatalf("unexpected echo payload: %s", string(payload))
	}

	// 2. Binary test
	sendBin := []byte{0x01, 0x02, 0x03, 0x04}
	if err := wsConn.WriteFrame(ws.OpBinary, sendBin, true); err != nil {
		t.Fatalf("write binary frame failed: %v", err)
	}

	op, payload, err = wsConn.ReadFrame()
	if err != nil {
		t.Fatalf("read binary frame failed: %v", err)
	}
	if op != ws.OpBinary {
		t.Fatalf("expected binary opcode 0x2, got 0x%X", op)
	}
	if !bytes.Equal(payload, []byte{0x04, 0x03, 0x02, 0x01}) {
		t.Fatalf("unexpected reversed binary payload: %v", payload)
	}
}

func TestSSE_Concurrency(t *testing.T) {
	s := server.NewServer()
	ts := httptest.NewServer(s.Mux)
	defer ts.Close()

	var wg sync.WaitGroup
	clientsCount := 10

	for i := 0; i < clientsCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req, _ := http.NewRequest("GET", ts.URL+"/sse", nil)
			client := &http.Client{Timeout: 500 * time.Millisecond}
			resp, err := client.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()
			buf := make([]byte, 256)
			_, _ = resp.Body.Read(buf)
		}()
	}

	time.Sleep(50 * time.Millisecond)
	for i := 0; i < 5; i++ {
		s.SSEHub.Broadcast("ping", "concurrent event")
	}

	wg.Wait()
}
