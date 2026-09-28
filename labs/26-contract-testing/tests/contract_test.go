package tests

import (
	"net/http"
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

func TestVerifier_HeaderValidation_And_ErrorBranches(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/bad-header", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"1"}`))
	})
	mux.HandleFunc("/v1/bad-json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid-json`))
	})
	mux.HandleFunc("/v1/status-mismatch", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	srv := httptest.NewServer(mux)
	defer srv.Close()

	verifier := contract.NewVerifier()

	// Header mismatch test
	c1 := &contract.Contract{
		Interactions: []contract.Interaction{
			{
				Description: "Bad Header Test",
				Request:     contract.RequestDefinition{Method: "GET", Path: "/v1/bad-header"},
				Response: contract.ResponseDefinition{
					Status:  http.StatusOK,
					Headers: map[string]string{"Content-Type": "application/json"},
					Body:    map[string]interface{}{"id": "1"},
				},
			},
		},
	}
	res1 := verifier.Verify(srv.URL, c1)
	if res1.Passed || len(res1.Errors) == 0 {
		t.Fatalf("expected failure on header mismatch")
	}

	// Invalid JSON test
	c2 := &contract.Contract{
		Interactions: []contract.Interaction{
			{
				Description: "Bad JSON Test",
				Request:     contract.RequestDefinition{Method: "GET", Path: "/v1/bad-json"},
				Response: contract.ResponseDefinition{
					Status: http.StatusOK,
					Body:   map[string]interface{}{"id": "1"},
				},
			},
		},
	}
	res2 := verifier.Verify(srv.URL, c2)
	if res2.Passed || len(res2.Errors) == 0 {
		t.Fatalf("expected failure on bad JSON")
	}

	// Status code mismatch test
	c3 := &contract.Contract{
		Interactions: []contract.Interaction{
			{
				Description: "Status Mismatch Test",
				Request:     contract.RequestDefinition{Method: "GET", Path: "/v1/status-mismatch"},
				Response: contract.ResponseDefinition{
					Status: http.StatusOK,
					Body:   map[string]interface{}{},
				},
			},
		},
	}
	res3 := verifier.Verify(srv.URL, c3)
	if res3.Passed || len(res3.Errors) == 0 {
		t.Fatalf("expected failure on status mismatch")
	}
}

func TestProviderDual_V2Endpoint_DirectAssertion(t *testing.T) {
	srv := httptest.NewServer(provider.NewProviderDual())
	defer srv.Close()

	res, err := srv.Client().Get(srv.URL + "/v2/orders/ORD-123")
	if err != nil {
		t.Fatalf("failed to query /v2 endpoint: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK from /v2, got %d", res.StatusCode)
	}
}
