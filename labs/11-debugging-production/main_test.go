package main

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestReproduceBug(t *testing.T) {
	handler := &Handler{
		svc: &Service{
			repo: &SlowRepository{},
		},
	}

	res, err := handler.HandleGetOrders(context.Background(), "1y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.TotalDuration < 4*time.Second {
		t.Fatalf("expected request to be slow (> 4s), but took %v", res.TotalDuration)
	}
}

func TestCollectEvidence(t *testing.T) {
	handler := &Handler{
		svc: &Service{
			repo: &SlowRepository{},
		},
	}

	res, err := handler.HandleGetOrders(context.Background(), "1y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fmt.Printf("request_id=abc123\n")
	fmt.Printf("handler_duration=%v\n", res.HandlerDuration)
	fmt.Printf("service_duration=%v\n", res.ServiceDuration)
	fmt.Printf("repository_duration=%v\n", res.RepositoryDuration)
	fmt.Printf("total_duration=%v\n", res.TotalDuration)
}

func TestVerifyFix(t *testing.T) {
	handler := &Handler{
		svc: &Service{
			repo: &FastRepository{},
		},
	}

	res, err := handler.HandleGetOrders(context.Background(), "1y")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.TotalDuration > 100*time.Millisecond {
		t.Fatalf("expected request to be fast (< 100ms), but took %v", res.TotalDuration)
	}
}
