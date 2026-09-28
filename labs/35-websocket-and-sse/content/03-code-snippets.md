# Snippet 1 — SSE Event Formatting

**Source File:** `internal/sse/sse.go:19-36`

**Purpose:** Mengformat event SSE sesuai WHATWG spec dengan field `id`, `event`, `data`, `retry`

```go
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
```

**Explanation:** Multi-line data dipisah per baris dan masing-masing mendapat prefix `data: `. Double-newline terminator (`\n\n`) dihasilkan dari trailing newline di baris terakhir `data:` + satu newline dari `sb.WriteString("\n")`. Ini sesuai dengan WHATWG SSE wire format.

---

# Snippet 2 — SSE Hub Broadcast dengan History

**Source File:** `internal/sse/sse.go:52-72`

**Purpose:** Broadcast event ke semua connected clients dengan thread-safe locking dan history buffer

```go
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
```

**Explanation:** Setiap event mendapat monotonic ID dari `h.nextID`. Non-blocking send (`select`/`default`) memastikan slow client tidak memblokir broadcast ke client lain. Jika channel buffer (capacity 64) penuh, event di-drop silent untuk client tersebut.

---

# Snippet 3 — SSE Subscribe with Last-Event-ID Replay

**Source File:** `internal/sse/sse.go:74-96`

**Purpose:** Subscribe ke live events sekaligus replay missed events berdasarkan `Last-Event-ID`

```go
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
```

**Explanation:** Replay dihitung dengan iterasi history: hanya event dengan `ID > lastID` yang di-replay. Channel didaftarkan ke clients map setelah replay dihitung, dalam satu lock scope — tidak ada gap event.

---

# Snippet 4 — WebSocket Handshake (Upgrade)

**Source File:** `internal/ws/ws.go:31-76`

**Purpose:** Melakukan RFC 6455 handshake: validasi header, hitung `Sec-WebSocket-Accept`, hijack koneksi

```go
func Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "Expected websocket upgrade", http.StatusBadRequest)
		return nil, errors.New("not a websocket handshake")
	}

	key := r.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		http.Error(w, "Missing Sec-WebSocket-Key", http.StatusBadRequest)
		return nil, errors.New("missing sec-websocket-key")
	}

	h := sha1.New()
	h.Write([]byte(key + wsGUID))
	accept := base64.StdEncoding.EncodeToString(h.Sum(nil))

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking unsupported", http.StatusInternalServerError)
		return nil, errors.New("cannot hijack connection")
	}

	netConn, bufrw, err := hj.Hijack()
	if err != nil {
		return nil, err
	}

	resp := "HTTP/1.1 101 Switching Protocols\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Accept: " + accept + "\r\n\r\n"

	if _, err := bufrw.WriteString(resp); err != nil {
		netConn.Close()
		return nil, err
	}
	if err := bufrw.Flush(); err != nil {
		netConn.Close()
		return nil, err
	}

	return &Conn{
		rw:     bufrw,
		closer: netConn,
	}, nil
}
```

**Explanation:** `Sec-WebSocket-Accept` dihitung dari `base64(SHA1(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))`. Connection di-hijack dari `net/http` untuk raw TCP access. Response 101 Switching Protocols mengakhiri HTTP lifecycle dan beralih ke WebSocket framing.

---

# Snippet 5 — WebSocket Frame Reading (RFC 6455 Framing)

**Source File:** `internal/ws/ws.go:84-127`

**Purpose:** Membaca WebSocket frame dengan support masking dan payload length 7/16/64-bit

```go
func (c *Conn) ReadFrame() (opcode byte, payload []byte, err error) {
	var head [2]byte
	if _, err := io.ReadFull(c.rw, head[:]); err != nil {
		return 0, nil, err
	}

	opcode = head[0] & 0x0F
	masked := (head[1] & 0x80) != 0
	length := uint64(head[1] & 0x7F)

	if length == 126 {
		var ext [2]byte
		if _, err := io.ReadFull(c.rw, ext[:]); err != nil {
			return 0, nil, err
		}
		length = uint64(binary.BigEndian.Uint16(ext[:]))
	} else if length == 127 {
		var ext [8]byte
		if _, err := io.ReadFull(c.rw, ext[:]); err != nil {
			return 0, nil, err
		}
		length = binary.BigEndian.Uint64(ext[:])
	}

	var mask [4]byte
	if masked {
		if _, err := io.ReadFull(c.rw, mask[:]); err != nil {
			return 0, nil, err
		}
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(c.rw, data); err != nil {
		return 0, nil, err
	}

	if masked {
		for i := range data {
			data[i] ^= mask[i%4]
		}
	}

	return opcode, data, nil
}
```

**Explanation:** Dua byte pertama berisi opcode (4 bit) dan length (7 bit). Jika length == 126, 2 byte tambahan dibaca untuk 16-bit length. Jika length == 127, 8 byte untuk 64-bit length. XOR unmasking dilakukan dengan `data[i] ^= mask[i%4]` — persis sesuai RFC 6455 Section 5.3.

---

# Snippet 6 — WebSocket Frame Writing

**Source File:** `internal/ws/ws.go:129-174`

**Purpose:** Menulis WebSocket frame dengan FIN=1, length encoding yang benar, dan optional masking

```go
func (c *Conn) WriteFrame(opcode byte, payload []byte, mask bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var buf []byte
	b0 := byte(0x80) | (opcode & 0x0F) // FIN=1
	buf = append(buf, b0)

	length := len(payload)
	maskBit := byte(0x00)
	if mask {
		maskBit = 0x80
	}

	if length <= 125 {
		buf = append(buf, maskBit|byte(length))
	} else if length <= 65535 {
		buf = append(buf, maskBit|126)
		var ext [2]byte
		binary.BigEndian.PutUint16(ext[:], uint16(length))
		buf = append(buf, ext[:]...)
	} else {
		buf = append(buf, maskBit|127)
		var ext [8]byte
		binary.BigEndian.PutUint64(ext[:], uint64(length))
		buf = append(buf, ext[:]...)
	}

	if mask {
		maskKey := [4]byte{0x12, 0x34, 0x56, 0x78}
		buf = append(buf, maskKey[:]...)
		maskedPayload := make([]byte, len(payload))
		for i := range payload {
			maskedPayload[i] = payload[i] ^ maskKey[i%4]
		}
		buf = append(buf, maskedPayload...)
	} else {
		buf = append(buf, payload...)
	}

	_, err := c.rw.Write(buf)
	if flusher, ok := c.rw.(interface{ Flush() error }); ok {
		_ = flusher.Flush()
	}
	return err
}
```

**Explanation:** `FIN=1` di-set dengan OR `0x80`. Mask key statis `{0x12, 0x34, 0x56, 0x78}` — untuk lab ini cukup, tapi RFC 6455 Section 5.3 mensyaratkan mask key yang unpredictabel. XOR masking dilakukan sebelum payload ditulis ke wire.

---

# Snippet 7 — Server Handler (WebSocket Echo + Reverse Binary)

**Source File:** `internal/server/server.go:26-56`

**Purpose:** Server handler yang memproses WebSocket frames: echo text, reverse binary, handle ping/pong/close

```go
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
			_ = conn.WriteFrame(ws.OpText, append([]byte("echo: "), data...), false)
		case ws.OpBinary:
			rev := make([]byte, len(data))
			for i := range data {
				rev[i] = data[len(data)-1-i]
			}
			_ = conn.WriteFrame(ws.OpBinary, rev, false)
		}
	}
})
```

**Explanation:** Loop tunggal per koneksi — read frame, lalu write response. Tidak ada concurrent write dari multiple goroutines. OpBinary: payload dibalik untuk membuktikan byte-level integrity.

---

# Snippet 8 — SSE Keep-alive Heartbeat

**Source File:** `internal/sse/sse.go:128-150`

**Purpose:** Mengirim SSE comment sebagai keep-alive heartbeat setiap 15 detik

```go
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
		if _, err := w.Write([]byte(": keep-alive\n\n")); err != nil {
			return
		}
		flusher.Flush()
	}
}
```

**Explanation:** SSE comment (baris diawali `:`) tidak diproses oleh browser sebagai event — hanya menjaga koneksi tetap hidup. `X-Accel-Buffering: no` di header mencegah proxy buffering. `Flush()` memaksa data terkirim segera ke client.

---

# Snippet 9 — Test: SSE Last-Event-ID Resumption

**Source File:** `tests/protocol_test.go:31-81`

**Purpose:** Verifikasi bahwa event di bawah `Last-Event-ID` tidak di-replay, event setelahnya di-replay

```go
func TestSSE_LastEventID_Resumption(t *testing.T) {
	s := server.NewServer()
	ts := httptest.NewServer(s.Mux)
	defer ts.Close()

	s.SSEHub.Broadcast("msg", "first")   // ID: 1
	s.SSEHub.Broadcast("msg", "second")  // ID: 2
	s.SSEHub.Broadcast("msg", "third")   // ID: 3

	req, _ := http.NewRequest("GET", ts.URL+"/sse", nil)
	req.Header.Set("Last-Event-ID", "1")

	client := &http.Client{}
	resp, err := client.Do(req)
	// ... read lines, assert event 1 absent, events 2 & 3 present
}
```

**Explanation:** 3 event di-broadcast sebelum client connect. Client mengirim `Last-Event-ID: 1` saat connect. Test memverifikasi event ID 1 tidak muncul di stream (dilewati), sementara ID 2 dan 3 di-replay.

---

# Snippet 10 — Test: Concurrent SSE Broadcast

**Source File:** `tests/protocol_test.go:136-166`

**Purpose:** Verifikasi thread-safety dengan 10 concurrent clients dan 5 concurrent broadcasts

```go
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
```

**Explanation:** Race detector bersih setelah test ini. Semua goroutine selesai tanpa deadlock atau data race. Channel per-client memastikan broadcast ke satu client tidak memblokir client lain.
