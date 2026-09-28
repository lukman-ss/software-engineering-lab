package ws

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
	"sync"
)

const (
	OpContinuation = 0x0
	OpText         = 0x1
	OpBinary       = 0x2
	OpClose        = 0x8
	OpPing         = 0x9
	OpPong         = 0xA

	wsGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

type Conn struct {
	rw     io.ReadWriter
	closer io.Closer
	mu     sync.Mutex
}

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

func (c *Conn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closer.Close()
}

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

func Dial(urlStr string, rawConn net.Conn) (*Conn, error) {
	req := "GET /ws HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"Upgrade: websocket\r\n" +
		"Connection: Upgrade\r\n" +
		"Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\n" +
		"Sec-WebSocket-Version: 13\r\n\r\n"

	if _, err := rawConn.Write([]byte(req)); err != nil {
		return nil, err
	}

	respBuf := make([]byte, 1024)
	n, err := rawConn.Read(respBuf)
	if err != nil {
		return nil, err
	}

	if string(respBuf[:12]) != "HTTP/1.1 101" {
		return nil, errors.New("failed handshake: " + string(respBuf[:n]))
	}

	return &Conn{
		rw:     rawConn,
		closer: rawConn,
	}, nil
}
