package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"

	"labs/26-contract-testing/internal/consumer"
	"labs/26-contract-testing/internal/contract"
	"labs/26-contract-testing/internal/provider"
)

func main() {
	fmt.Println("=== Contract Testing Lab: Consumer-Driven Contracts & CI Verification ===")

	// 1. Generate Consumer Contract
	fmt.Println("\n[Stage 1] Consumer generates contract:")
	c := consumer.GenerateMobileContract()
	rawJSON, _ := json.MarshalIndent(c, "", "  ")
	fmt.Printf("Generated Contract (%s -> %s):\n%s\n", c.Consumer, c.Provider, string(rawJSON))

	verifier := contract.NewVerifier()

	// 2. Test against Provider V1 (Compliant)
	fmt.Println("\n[Stage 2] Running Provider V1 Contract Verification:")
	srvV1 := httptest.NewServer(provider.NewProviderV1())
	defer srvV1.Close()

	resV1 := verifier.Verify(srvV1.URL, c)
	if resV1.Passed {
		fmt.Println("Result: PASSED. Provider V1 satisfies Mobile consumer contract.")
		fmt.Println("CI Deployment Gate: ALLOWED.")
	} else {
		fmt.Printf("Result: FAILED. Errors: %v\n", resV1.Errors)
		os.Exit(1)
	}

	// 3. Test against Breaking Provider
	fmt.Println("\n[Stage 3] Running Breaking Provider Contract Verification:")
	srvBreaking := httptest.NewServer(provider.NewProviderBreaking())
	defer srvBreaking.Close()

	resBreaking := verifier.Verify(srvBreaking.URL, c)
	if !resBreaking.Passed {
		fmt.Println("Result: BLOCKED! Breaking changes detected before deployment:")
		for i, errStr := range resBreaking.Errors {
			fmt.Printf("  %d. %s\n", i+1, errStr)
		}
		fmt.Println("CI Deployment Gate: PREVENTED PRODUCTION OUTAGE.")
	} else {
		fmt.Println("Result: Unexpected pass! Verification logic failure.")
		os.Exit(1)
	}

	// 4. Test against Dual Provider (V1 + V2 Evolutionary Path)
	fmt.Println("\n[Stage 4] Running Dual Provider (V1 + V2) Verification:")
	srvDual := httptest.NewServer(provider.NewProviderDual())
	defer srvDual.Close()

	resDual := verifier.Verify(srvDual.URL, c)
	if resDual.Passed {
		fmt.Println("Result: PASSED. Dual provider maintains backwards-compatible V1 contract.")
		fmt.Println("CI Deployment Gate: ALLOWED for independent canary/migration.")
	} else {
		fmt.Printf("Result: FAILED. Errors: %v\n", resDual.Errors)
		os.Exit(1)
	}

	fmt.Println("\n=== Contract Testing Demonstration Complete ===")
}
