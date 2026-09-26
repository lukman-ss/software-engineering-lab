package tests

import (
	"net/http/httptest"
	"sync"
	"testing"

	"labs/26-contract-testing/internal/consumer"
	"labs/26-contract-testing/internal/contract"
	"labs/26-contract-testing/internal/provider"
)

func TestConsumerContractGeneration(t *testing.T) {
	c := consumer.GenerateMobileContract()
	if c.Consumer != "MobileApp" || c.Provider != "OrderService" {
		t.Fatalf("unexpected consumer/provider names: %s -> %s", c.Consumer, c.Provider)
	}
	if len(c.Interactions) != 1 {
		t.Fatalf("expected 1 interaction, got %d", len(c.Interactions))
	}
	interaction := c.Interactions[0]
	if interaction.Request.Path != "/v1/orders/ORD-123" {
		t.Errorf("expected path /v1/orders/ORD-123, got %s", interaction.Request.Path)
	}
}

func TestProviderV1_ContractVerification_Success(t *testing.T) {
	srv := httptest.NewServer(provider.NewProviderV1())
	defer srv.Close()

	c := consumer.GenerateMobileContract()
	verifier := contract.NewVerifier()
	result := verifier.Verify(srv.URL, c)

	if !result.Passed {
		t.Fatalf("expected contract verification to pass, got errors: %v", result.Errors)
	}

	// Also verify mobile client end-to-end against Provider V1
	client := consumer.NewMobileOrderClient(srv.URL)
	order, err := client.FetchOrder("ORD-123")
	if err != nil {
		t.Fatalf("client FetchOrder failed: %v", err)
	}
	if order.CustomerName != "Budi Santoso" || order.Status != "IN_PROGRESS" || order.Total != 150000 {
		t.Errorf("client parsed unexpected values: %+v", order)
	}
}

func TestProviderBreaking_ContractVerification_Fails(t *testing.T) {
	srv := httptest.NewServer(provider.NewProviderBreaking())
	defer srv.Close()

	c := consumer.GenerateMobileContract()
	verifier := contract.NewVerifier()
	result := verifier.Verify(srv.URL, c)

	if result.Passed {
		t.Fatalf("expected contract verification to FAIL for breaking provider, but it passed")
	}

	if len(result.Errors) < 3 {
		t.Fatalf("expected at least 3 breaking contract errors, got %d: %v", len(result.Errors), result.Errors)
	}

	// Assert Mobile client actually fails when calling breaking provider
	client := consumer.NewMobileOrderClient(srv.URL)
	_, err := client.FetchOrder("ORD-123")
	if err == nil {
		t.Fatalf("expected mobile client to fail on breaking provider schema, but got no error")
	}
}

func TestProviderDual_ContractVerification_Success(t *testing.T) {
	srv := httptest.NewServer(provider.NewProviderDual())
	defer srv.Close()

	c := consumer.GenerateMobileContract()
	verifier := contract.NewVerifier()
	result := verifier.Verify(srv.URL, c)

	if !result.Passed {
		t.Fatalf("expected contract verification to pass for dual provider V1 path, got errors: %v", result.Errors)
	}

	client := consumer.NewMobileOrderClient(srv.URL)
	order, err := client.FetchOrder("ORD-123")
	if err != nil {
		t.Fatalf("client FetchOrder failed on dual provider: %v", err)
	}
	if order.Status != "IN_PROGRESS" || order.CustomerName != "Budi Santoso" {
		t.Errorf("unexpected order response from dual provider: %+v", order)
	}
}

func TestConcurrentContractVerification(t *testing.T) {
	srv := httptest.NewServer(provider.NewProviderV1())
	defer srv.Close()

	c := consumer.GenerateMobileContract()
	verifier := contract.NewVerifier()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res := verifier.Verify(srv.URL, c)
			if !res.Passed {
				t.Errorf("concurrent contract verification failed: %v", res.Errors)
			}
		}()
	}
	wg.Wait()
}
