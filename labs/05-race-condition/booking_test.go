package race

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// Testing unique violation error codes
var (
	ErrDuplicateKey  = errors.New("duplicate key violation")
	ErrAlreadyBooked = errors.New("slot already booked")
)

// BookingRepository contract untuk booking service
type BookingRepository interface {
	CreateBooking(ctx context.Context, serviceID, slotTime string) error
	CountBookings(ctx context.Context, serviceID, slotTime string) (int, error)
}

// MockBookingRepository implementasi mock untuk testing
type MockBookingRepository struct {
	mu       sync.RWMutex
	bookings map[string]bool // key = serviceID:slotTime
}

func NewMockBookingRepository() *MockBookingRepository {
	return &MockBookingRepository{
		bookings: make(map[string]bool),
	}
}

// CreateBooking mencoba membuat booking dengan UNIQUE constraint check
func (m *MockBookingRepository) CreateBooking(_ context.Context, serviceID, slotTime string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := serviceID + ":" + slotTime
	if m.bookings[key] {
		return ErrDuplicateKey
	}
	m.bookings[key] = true
	return nil
}

// CountBookings menghitung jumlah booking
func (m *MockBookingRepository) CountBookings(_ context.Context, serviceID, slotTime string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := serviceID + ":" + slotTime
	if m.bookings[key] {
		return 1, nil
	}
	return 0, nil
}

// Test500_ConcurrentBooking menguji 500 concurrent requests booking exclusive slot.
//
// Setup:
// - attempts = 500 goroutines
//
// Expected:
// - created = 1
// - conflict = 499
// - final_count = 1
//
// Assert:
// - COUNT(bookings for slot) == 1
func Test500_ConcurrentBooking(t *testing.T) {
	const attempts = 500

	repo := NewMockBookingRepository()
	ctx := context.Background()

	var createdCount int
	var conflictCount int
	var mu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(attempts)

	// All 500 goroutines compete for 1 exclusive slot
	for i := 0; i < attempts; i++ {
		go func(id int) {
			defer wg.Done()
			err := repo.CreateBooking(ctx, "service-01", "09:00")
			if err != nil {
				mu.Lock()
				conflictCount++
				mu.Unlock()
				return
			}
			mu.Lock()
			createdCount++
			mu.Unlock()
		}(i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout: potential goroutine leak")
	}

	// Verify final count
	finalCount, _ := repo.CountBookings(ctx, "service-01", "09:00")

	t.Logf("=== CONCURRENT BOOKING RESULTS ===")
	t.Logf("Created: %d", createdCount)
	t.Logf("Conflict/Already booked: %d", conflictCount)
	t.Logf("Final booking count: %d", finalCount)

	// Assertions
	if createdCount != 1 {
		t.Errorf("expected created = 1, got %d", createdCount)
	}

	if conflictCount != 499 {
		t.Errorf("expected conflict = 499, got %d", conflictCount)
	}

	if finalCount != 1 {
		t.Errorf("expected final_count = 1, got %d", finalCount)
	}

	// Main invariant: COUNT(bookings for slot) == 1
	t.Logf("✅ INVARIANT HOLDS: COUNT(bookings) = %d", finalCount)
}

// TestBooking_ErrorHandling verifies unique violation error handling
func TestBooking_ErrorHandling(t *testing.T) {
	repo := NewMockBookingRepository()
	ctx := context.Background()

	// First booking should succeed
	err := repo.CreateBooking(ctx, "service-01", "09:00")
	if err != nil {
		t.Errorf("first booking should succeed, got error: %v", err)
	}

	// Second booking should fail with meaningful error
	err = repo.CreateBooking(ctx, "service-01", "09:00")
	if err == nil {
		t.Error("second booking should fail")
	}
	if !errors.Is(err, ErrDuplicateKey) && err != ErrAlreadyBooked {
		t.Errorf("expected duplicate key or already booked error, got: %v", err)
	}

	t.Logf("✅ Error handling works: duplicate key -> %v", err)
}
