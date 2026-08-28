// Package util provides shared utilities.
package util

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// newID generates a ULID-like identifier.
// For production, use a proper ULID library or the system's uuid.
func NewID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b[:])[:16] // 16-char hex ID
}

// NewTxID generates a transaction ID.
func NewTxID() string {
	return "tx_" + NewID()
}

// Now returns current time.
func Now() time.Time {
	return time.Now()
}

// NewOrderID generates an order ID.
func NewOrderID() string {
	return "ord_" + NewID()
}

// NewPaymentID generates a payment ID.
func NewPaymentID() string {
	return "pay_" + NewID()
}

// NewInventoryID generates an inventory item ID.
func NewInventoryID() string {
	return "inv_" + NewID()
}

// NewWalletID generates a wallet ID.
func NewWalletID() string {
	return "wal_" + NewID()
}

// NewNotificationID generates a notification ID.
func NewNotificationID() string {
	return "notif_" + NewID()
}
