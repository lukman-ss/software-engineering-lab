package main

import (
	"os"
	"testing"
)

func TestResolveServiceName(t *testing.T) {
	originalValue, existed := os.LookupEnv("OTEL_SERVICE_NAME")

	t.Cleanup(func() {
		if existed {
			_ = os.Setenv("OTEL_SERVICE_NAME", originalValue)
			return
		}

		_ = os.Unsetenv("OTEL_SERVICE_NAME")
	})

	t.Run("defaults when OTEL_SERVICE_NAME is unset", func(t *testing.T) {
		_ = os.Unsetenv("OTEL_SERVICE_NAME")
		if got := resolveServiceName(); got != "lab07-observability" {
			t.Fatalf("expected lab07-observability, got %s", got)
		}
	})

	t.Run("defaults when OTEL_SERVICE_NAME is empty", func(t *testing.T) {
		t.Setenv("OTEL_SERVICE_NAME", "")
		if got := resolveServiceName(); got != "lab07-observability" {
			t.Fatalf("expected lab07-observability, got %s", got)
		}
	})

	t.Run("uses custom OTEL_SERVICE_NAME", func(t *testing.T) {
		t.Setenv("OTEL_SERVICE_NAME", "custom-service")
		if got := resolveServiceName(); got != "custom-service" {
			t.Fatalf("expected custom-service, got %s", got)
		}
	})
}
