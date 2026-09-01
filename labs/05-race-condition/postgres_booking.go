//go:build integration
// +build integration

package race

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

// PostgresBookingRepository implements BookingRepository using PostgreSQL.
type PostgresBookingRepository struct {
	db *sql.DB
}

func NewPostgresBookingRepository(db *sql.DB) *PostgresBookingRepository {
	return &PostgresBookingRepository{db: db}
}

// CreateBooking attempts to insert a new booking.
// It relies on the database UNIQUE(branch_id, service_date, slot_time) constraint.
func (r *PostgresBookingRepository) CreateBooking(ctx context.Context, branchID, slotTime string) error {
	// Parse slotTime assuming it contains both date and time if it's longer than 5 chars,
	// otherwise assume today's date with the given time.
	// For simplicity in this lab, we'll hardcode date and time parsing based on expected input.
	var serviceDate string
	var parsedTime string

	if len(slotTime) > 5 && strings.Contains(slotTime, " ") {
		parts := strings.Split(slotTime, " ")
		serviceDate = parts[0]
		parsedTime = parts[1]
	} else {
		serviceDate = time.Now().Format("2006-01-02")
		parsedTime = slotTime
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO service_bookings (branch_id, customer_id, service_date, slot_time)
		 VALUES ($1, $2, $3, $4)`,
		branchID, "customer-test", serviceDate, parsedTime)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// 23505 is unique_violation
			if pqErr.Code == "23505" {
				return ErrDuplicateKey
			}
		}
		return fmt.Errorf("create booking: %w", err)
	}

	return nil
}

// CountBookings counts bookings for a specific branch and slot.
func (r *PostgresBookingRepository) CountBookings(ctx context.Context, branchID, slotTime string) (int, error) {
	var serviceDate string
	var parsedTime string

	if len(slotTime) > 5 && strings.Contains(slotTime, " ") {
		parts := strings.Split(slotTime, " ")
		serviceDate = parts[0]
		parsedTime = parts[1]
	} else {
		serviceDate = time.Now().Format("2006-01-02")
		parsedTime = slotTime
	}

	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM service_bookings
		 WHERE branch_id = $1 AND service_date = $2 AND slot_time = $3`,
		branchID, serviceDate, parsedTime).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("count bookings: %w", err)
	}

	return count, nil
}
