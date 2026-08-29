package errors_test

import (
	"errors"
	"testing"

	apperrors "github.com/lukman-ss/software-engineering-lab/pkg/errors"
)

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
