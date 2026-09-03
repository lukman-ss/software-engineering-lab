package main

import (
	"os"
	"testing"
)

func TestResolveServiceName(t *testing.T) {
	originalValue := os.Getenv("OTEL_SERVICE_NAME")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("OTEL_SERVICE_NAME")
		} else {
			os.Setenv("OTEL_SERVICE_NAME", originalValue)
		}
	}()

	t.Run("unset uses default", func(t *testing.T) {
		os.Unsetenv("OTEL_SERVICE_NAME")
		if got := resolveServiceName(); got != "lab07-observability" {
			t.Fatalf("expected lab07-observability, got %s", got)
		}
	})

	t.Run("empty string uses default", func(t *testing.T) {
		os.Setenv("OTEL_SERVICE_NAME", "")
		if got := resolveServiceName(); got != "lab07-observability" {
			t.Fatalf("expected lab07-observability, got %s", got)
		}
	})

	t.Run("custom value is used", func(t *testing.T) {
		os.Setenv("OTEL_SERVICE_NAME", "custom-service")
		if got := resolveServiceName(); got != "custom-service" {
			t.Fatalf("expected custom-service, got %s", got)
		}
	})
}
