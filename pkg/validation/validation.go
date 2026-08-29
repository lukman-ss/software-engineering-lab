// Package validation provides request validation utilities.
package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	apperrors "github.com/lukman-ss/software-engineering-lab/pkg/errors"
)

// MaxBodySize is the default maximum request body size (1 MB).
const MaxBodySize = 1 << 20

// ReadJSON parses a JSON body with a maximum size and rejects unknown fields.
// Returns structured errors with categories:
// - MalformedJSON for parse errors
// - BodyTooLarge for size violations
// - Validation for missing/invalid fields
func ReadJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return apperrors.ErrInvalidRequest.WithCause(errors.New("empty body"))
	}

	// Limit reader
	r.Body = http.MaxBytesReader(nil, r.Body, MaxBodySize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return classifyJSONError(err)
	}

	// Reject trailing garbage
	if dec.More() {
		return apperrors.ErrInvalidRequest.WithCause(errors.New("unexpected trailing data"))
	}

	return nil
}

func classifyJSONError(err error) error {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		return apperrors.New("BODY_TOO_LARGE",
			fmt.Sprintf("request body exceeds maximum size of %d bytes", MaxBodySize),
			apperrors.Validation)
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return apperrors.New("MALFORMED_JSON",
			fmt.Sprintf("invalid JSON syntax at byte %d", syntaxErr.Offset),
			apperrors.Validation).WithCause(err)
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return apperrors.New("INVALID_FIELD_TYPE",
			fmt.Sprintf("field %q has invalid type", typeErr.Field),
			apperrors.Validation).WithCause(err)
	}

	if strings.HasPrefix(err.Error(), "json: unknown field ") {
		field := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return apperrors.New("UNKNOWN_FIELD",
			fmt.Sprintf("unknown field %s", field),
			apperrors.Validation).WithCause(err)
	}

	if errors.Is(err, io.EOF) {
		return apperrors.New("MALFORMED_JSON", "empty body", apperrors.Validation)
	}

	return apperrors.New("MALFORMED_JSON", "could not parse JSON", apperrors.Validation).WithCause(err)
}

// MissingField creates a validation error for a missing required field.
func MissingField(name string) error {
	return apperrors.New("MISSING_FIELD",
		fmt.Sprintf("missing required field %q", name),
		apperrors.Validation)
}

// InvalidValue creates a validation error for an invalid value.
func InvalidValue(name, reason string) error {
	return apperrors.New("INVALID_VALUE",
		fmt.Sprintf("invalid value for %q: %s", name, reason),
		apperrors.Validation)
}
