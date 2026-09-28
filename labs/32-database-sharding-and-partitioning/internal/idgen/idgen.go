package idgen

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync"
	"time"
)

// UUIDv7 generates an RFC 9562 compliant time-ordered UUIDv7.
// 48-bit timestamp | 4-bit ver (0111) | 12-bit rand_a | 2-bit var (10) | 62-bit rand_b
func NewUUIDv7() (string, error) {
	var uuid [16]byte

	ms := uint64(time.Now().UnixMilli())
	uuid[0] = byte(ms >> 40)
	uuid[1] = byte(ms >> 32)
	uuid[2] = byte(ms >> 24)
	uuid[3] = byte(ms >> 16)
	uuid[4] = byte(ms >> 8)
	uuid[5] = byte(ms)

	if _, err := rand.Read(uuid[6:]); err != nil {
		return "", err
	}

	// version 7: 0111
	uuid[6] = (uuid[6] & 0x0F) | 0x70
	// variant RFC 4122/9562: 10xxxxxx
	uuid[8] = (uuid[8] & 0x3F) | 0x80

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16]), nil
}

// SequenceBlockAllocator simulates Vitess Sequences block allocation for auto-increment IDs.
type SequenceBlockAllocator struct {
	mu        sync.Mutex
	blockSize int64
	current   int64
	max       int64
	fetcher   func(blockSize int64) (int64, error)
}

func NewSequenceBlockAllocator(blockSize int64, fetcher func(blockSize int64) (int64, error)) *SequenceBlockAllocator {
	return &SequenceBlockAllocator{
		blockSize: blockSize,
		fetcher:   fetcher,
	}
}

func (s *SequenceBlockAllocator) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.current >= s.max {
		base, err := s.fetcher(s.blockSize)
		if err != nil {
			return 0, err
		}
		s.current = base
		s.max = base + s.blockSize
	}

	id := s.current
	s.current++
	return id, nil
}

// MemoryCentralSequence acts as the central coordinator giving ID blocks.
type MemoryCentralSequence struct {
	mu     sync.Mutex
	cursor int64
}

func (m *MemoryCentralSequence) AllocateBlock(size int64) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	start := m.cursor
	m.cursor += size
	return start + 1, nil
}

// EncodeTimeFromUUIDv7 extracts millisecond timestamp from UUIDv7 string.
func ExtractTimeFromUUIDv7(u string) (time.Time, error) {
	var clean string
	for _, c := range u {
		if c != '-' {
			clean += string(c)
		}
	}
	if len(clean) != 32 {
		return time.Time{}, fmt.Errorf("invalid uuid length")
	}
	var high [8]byte
	n, err := fmt.Sscanf(clean[:12], "%12x", &high)
	if err != nil && n == 0 {
		// Parse byte by byte manually
		var b []byte
		for i := 0; i < 12; i += 2 {
			var val byte
			fmt.Sscanf(clean[i:i+2], "%02x", &val)
			b = append(b, val)
		}
		ms := uint64(binary.BigEndian.Uint32(b[0:4]))<<16 | uint64(binary.BigEndian.Uint16(b[4:6]))
		return time.UnixMilli(int64(ms)), nil
	}
	return time.Time{}, nil
}
