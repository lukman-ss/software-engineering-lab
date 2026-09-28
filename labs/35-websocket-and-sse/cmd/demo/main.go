package main

import (
	"bufio"
	"bytes"
	"fmt"
	"labs/35-websocket-and-sse/internal/server"
	"labs/35-websocket-and-sse/internal/ws"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func main() {
	s := server.NewServer()
	ts := httptest.NewServer(s.Mux)
	defer ts.Close()

	fmt.Println("=== DEMO: Server-Sent Events (SSE) vs WebSocket ===")
	fmt.Printf("Server listening on %s\n\n", ts.URL)

	// 1. SSE Demonstration
	fmt.Println("--- 1. SSE: Unidirectional Server Push with Last-Event-ID Resumption ---")
	
	// Broadcast initial events
	s.SSEHub.Broadcast("news", "Breaking: Go 1.22 released")
	s.SSEHub.Broadcast("news", "Update: Faster routing in net/http")

	// Client reconnects specifying Last-Event-ID: 1
	req, _ := http.NewRequest("GET", ts.URL+"/sse", nil)
	req.Header.Set("Last-Event-ID", "1")

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(resp.Body)
	fmt.Println("Client connected with Last-Event-ID: 1. Received resumed events:")
	
	// Read until first non-empty event
	for i := 0; i < 4; i++ {
		line, _ := reader.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		if line != "" {
			fmt.Printf("  [SSE Stream] %s\n", line)
		}
	}
	resp.Body.Close()

	// 2. WebSocket Demonstration
	fmt.Println("\n--- 2. WebSocket: Full-Duplex Text and Binary Protocol ---")
	addr := ts.Listener.Addr().String()
	rawConn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer rawConn.Close()

	wsConn, err := ws.Dial("/ws", rawConn)
	if err != nil {
		panic(err)
	}

	// Send Text Frame
	textMsg := "Hello RFC 6455"
	fmt.Printf("Client sending WS text frame: %q\n", textMsg)
	_ = wsConn.WriteFrame(ws.OpText, []byte(textMsg), true)

	op, payload, err := wsConn.ReadFrame()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server responded with Opcode 0x%X (Text): %q\n", op, string(payload))

	// Send Binary Frame
	binMsg := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	fmt.Printf("Client sending WS binary frame: %X\n", binMsg)
	_ = wsConn.WriteFrame(ws.OpBinary, binMsg, true)

	op, payload, err = wsConn.ReadFrame()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Server responded with Opcode 0x%X (Binary - Reversed): %X\n", op, payload)

	if bytes.Equal(payload, []byte{0xEF, 0xBE, 0xAD, 0xDE}) {
		fmt.Println("Binary frame payload verified successfully.")
	}

	fmt.Println("\n=== Demo Complete ===")
}
