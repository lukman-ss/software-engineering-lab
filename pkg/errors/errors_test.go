package errors_test

import (
	"errors"
	"testing"

	apperrors "github.com/lukman-ss/software-engineering-lab/pkg/errors"
)

func TestAppError(t *testing.T) {
	tests := []struct {
		name    string
		err     *apperrors.AppError
		wantMsg string
		wantCat apperrors.Category
	}{
		{
			name:    "NotFound error",
			err:     apperrors.ErrOrderNotFound,
			wantMsg: "order not found",
			wantCat: apperrors.NotFound,
		},
		{
			name:    "Validation error",
			err:     apperrors.ErrInvalidRequest,
			wantMsg: "invalid request",
			wantCat: apperrors.Validation,
		},
		{
			name:    "Conflict error",
			err:     apperrors.ErrConcurrentUpdate,
			wantMsg: "concurrent update",
			wantCat: apperrors.Conflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Message != tt.wantMsg {
				t.Errorf("expected message %q, got %q", tt.wantMsg, tt.err.Message)
			}
			if tt.err.Category != tt.wantCat {
				t.Errorf("expected category %v, got %v", tt.wantCat, tt.err.Category)
			}
		})
	}
}

func TestWithError(t *testing.T) {
	cause := errors.New("database connection failed")
	err := apperrors.ErrPaymentFailed.WithCause(cause)

	if err.Error() != "payment processing failed: database connection failed" {
		t.Errorf("unexpected error string: %s", err.Error())
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != cause {
		t.Errorf("cause not preserved")
	}
}

func TestFromError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want apperrors.Category
	}{
		{"AppError", apperrors.ErrOrderNotFound, apperrors.NotFound},
		{"Generic error", errors.New("something"), apperrors.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extracted := apperrors.FromError(tt.err)
			if extracted.Category != tt.want {
				t.Errorf("expected category %v, got %v", tt.want, extracted.Category)
			}
		})
	}
}

func TestHTTPStatusForCategory(t *testing.T) {
	tests := []struct {
		cat  apperrors.Category
		want int
	}{
		{apperrors.Validation, 400},
		{apperrors.NotFound, 404},
		{apperrors.Conflict, 409},
		{apperrors.Unauthorized, 401},
		{apperrors.Forbidden, 403},
		{apperrors.RateLimited, 429},
		{apperrors.Internal, 500},
		{apperrors.DependencyFailure, 502},
		{apperrors.Timeout, 504},
	}

	for _, tt := range tests {
		t.Run(string(tt.cat), func(t *testing.T) {
			got := apperrors.HTTPStatus(tt.cat)
			if got != tt.want {
				t.Errorf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
