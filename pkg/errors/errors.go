// Package errors defines the application error model.
package errors

import "errors"

// Category classifies errors for HTTP mapping and observability.
type Category string

const (
	Validation       Category = "validation"
	NotFound         Category = "not_found"
	Conflict         Category = "conflict"
	Unauthorized     Category = "unauthorized"
	Forbidden        Category = "forbidden"
	RateLimited      Category = "rate_limited"
	Internal         Category = "internal"
	DependencyFailure Category = "dependency_failure"
	Timeout          Category = "timeout"
)

// AppError is the structured error returned to clients.
// Internal details stay in logs, not in the response.
type AppError struct {
	Code    string
	Message string
	Category Category
	RequestID string
	// Cause is the wrapped error (logged server-side only)
	cause error
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return e.Message + ": " + e.cause.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.cause
}

// WithCause wraps an underlying error for logging.
func (e *AppError) WithCause(cause error) *AppError {
	e.cause = cause
	return e
}

// New creates an AppError.
func New(code, message string, category Category) *AppError {
	return &AppError{
		Code:     code,
		Message:  message,
		Category: category,
	}
}

// FromError attempts to extract an AppError from any error.
// If not found, returns a generic Internal error.
func FromError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return New("INTERNAL", "internal server error", Internal).WithCause(err)
}

// HTTPStatus returns the standard HTTP status code for a category.
func HTTPStatus(c Category) int {
	switch c {
	case Validation:
		return 400
	case NotFound:
		return 404
	case Conflict:
		return 409
	case Unauthorized:
		return 401
	case Forbidden:
		return 403
	case RateLimited:
		return 429
	case Timeout:
		return 504
	case DependencyFailure:
		return 502
	default:
		return 500
	}
}

// Sentinel errors for domain use
var (
	ErrOrderNotFound   = New("ORDER_NOT_FOUND", "order not found", NotFound)
	ErrPaymentNotFound = New("PAYMENT_NOT_FOUND", "payment not found", NotFound)
	ErrInventoryNotFound = New("INVENTORY_NOT_FOUND", "inventory not found", NotFound)
	ErrWalletNotFound  = New("WALLET_NOT_FOUND", "wallet not found", NotFound)

	ErrInvalidStatus     = New("INVALID_STATUS", "invalid status transition", Conflict)
	ErrIdempotencyConflict = New("IDEMPOTENCY_CONFLICT", "idempotency key conflict", Conflict)
	ErrInsufficientStock = New("INSUFFICIENT_STOCK", "insufficient stock", Conflict)
	ErrInsufficientBalance = New("INSUFFICIENT_BALANCE", "insufficient balance", Conflict)
	ErrConcurrentUpdate  = New("CONCURRENT_UPDATE", "resource updated by another request", Conflict)

	ErrInvalidRequest   = New("INVALID_REQUEST", "invalid request", Validation)
	ErrPaymentFailed    = New("PAYMENT_FAILED", "payment processing failed", DependencyFailure)

	ErrRateLimited      = New("RATE_LIMITED", "too many requests", RateLimited)
	ErrUnauthorized     = New("UNAUTHORIZED", "unauthorized", Unauthorized)
	ErrForbidden        = New("FORBIDDEN", "forbidden", Forbidden)
)
