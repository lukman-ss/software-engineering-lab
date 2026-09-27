# Code Snippets

Source file references use real file paths from `labs/26-contract-testing`. Code is reproduced exactly from the approved implementation (with comments preserved), not re-implemented.

---

## Snippet 1 — Consumer Contract Definition

Source File: `internal/consumer/client.go:82-111`
Purpose: Build the consumer-driven contract as executable expectations. Generates a single `Interaction` capturing the minimal subset required by the Mobile App consumer: `GET /v1/orders/ORD-123`, expected `id`, `status` (`IN_PROGRESS`), `customer.name`, and `total` (`json.Number` to preserve integer semantics).

```go
func GenerateMobileContract() *contract.Contract {
	return &contract.Contract{
		Consumer: "MobileApp",
		Provider: "OrderService",
		Interactions: []contract.Interaction{
			{
				Description:   "A request for order details by ID",
				ProviderState: "Order ORD-123 exists and is IN_PROGRESS",
				Request: contract.RequestDefinition{
					Method: http.MethodGet,
					Path:   "/v1/orders/ORD-123",
				},
				Response: contract.ResponseDefinition{
					Status: http.StatusOK,
					Headers: map[string]string{
						"Content-Type": "application/json",
					},
					Body: map[string]interface{}{
						"id":     "ORD-123",
						"status": "IN_PROGRESS",
						"customer": map[string]interface{}{
							"name": "Budi Santoso",
						},
						"total": json.Number("150000"),
					},
				},
			},
		},
	}
}
```

Explanation: The consumer, not the provider, defines the contract. `json.Number` is used for `total` so the verifier can distinguish integer `150000` from string `"150000"` — capturing a semantic type change that JSON Schema alone would accept.

---

## Snippet 2 — Provider V1 (Contract Compliant)

Source File: `internal/provider/server.go:12-44`
Purpose: Implement the original V1 provider whose response body exactly satisfies the consumer contract, passing verification.

```go
// ProviderV1 implements original API server fulfilling V1 contract.
type ProviderV1 struct{}

func NewProviderV1() *ProviderV1 {
	return &ProviderV1{}
}

func (p *ProviderV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/v1/orders/") {
		id := strings.TrimPrefix(r.URL.Path, "/v1/orders/")
		resp := model.OrderResponseV1{
			ID:     id,
			Status: "IN_PROGRESS",
			Customer: model.CustomerResponseV1{
				ID:   "cust-100",
				Name: "Budi Santoso",
			},
			Total: 150000,
			Notes: "Internal note not used by mobile app",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	http.NotFound(w, r)
}
```

Explanation: `Total: 150000` (int64) matches the consumer contract's `json.Number("150000")`. `Notes` exists in the response but is **not** part of the contract — demonstrating subset verification: consumer-declared fields only.

---

## Snippet 3 — Provider Breaking (Contract Violation)

Source File: `internal/provider/server.go:46-80`
Purpose: Implement the provider that introduces exactly the three breaking changes the contract gate is designed to detect.

```go
// ProviderBreaking implements API server with breaking schema changes.
type ProviderBreaking struct{}

func NewProviderBreaking() *ProviderBreaking {
	return &ProviderBreaking{}
}

func (p *ProviderBreaking) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/v1/orders/") {
		id := strings.TrimPrefix(r.URL.Path, "/v1/orders/")
		// Introduces 3 breaking changes: status casing, field rename customer.full_name, total string
		resp := model.OrderResponseBreaking{
			ID:     id,
			Status: "in_progress",
			Customer: model.CustomerResponseBreaking{
				ID:       "cust-100",
				FullName: "Budi Santoso",
			},
			Total: "150000",
			Notes: "Breaking change version",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	http.NotFound(w, r)
}
```

Explanation: Three breaking mutations vs the contract:
1. `"IN_PROGRESS"` → `"in_progress"` — enum casing.
2. `customer.name` → `customer.full_name` — field rename, `name` becomes missing.
3. `Total: 150000` (int64) → `Total: "150000"` (string) — primitive type mutation.

---

## Snippet 4 — Provider Dual (V1 + V2 Safe Evolution)

Source File: `internal/provider/server.go:82-131`
Purpose: Implement the dual-versioned provider that preserves V1 contract compatibility while introducing a V2 schema, enabling independent consumer migration.

```go
// ProviderDual supports V1 (contract compliant) alongside V2 (new feature schema).
type ProviderDual struct{}

func NewProviderDual() *ProviderDual {
	return &ProviderDual{}
}

func (p *ProviderDual) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/v1/orders/") {
		id := strings.TrimPrefix(r.URL.Path, "/v1/orders/")
		respV1 := model.OrderResponseV1{
			ID:     id,
			Status: "IN_PROGRESS",
			Customer: model.CustomerResponseV1{
				ID:   "cust-100",
				Name: "Budi Santoso",
			},
			Total: 150000,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(respV1)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/v2/orders/") {
		id := strings.TrimPrefix(r.URL.Path, "/v2/orders/")
		respV2 := model.OrderResponseV2{
			ID:     id,
			Status: "in_progress",
			Customer: model.CustomerResponseV2{
				ID:       "cust-100",
				FullName: "Budi Santoso",
			},
			Total:    "150000",
			Currency: "IDR",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(respV2)
		return
	}

	http.NotFound(w, r)
}
```

Explanation: `/v1/orders/{id}` continues to serve the compliant V1 schema (`IN_PROGRESS`, `customer.name`, integer `total`). `/v2/orders/{id}` introduces the evolved schema (`in_progress`, `customer.full_name`, string `total`, added `currency`) as a parallel, opt-in endpoint — the expand phase of expand/contract. **Note:** V2 endpoint exists but has no associated verification test; only V1 contract verification is exercised in the lab (see GAP-02 in `engineering-audit-opensource/06-verdict.md`).

---

## Snippet 5 — Consumer Client Field Validation

Source File: `internal/consumer/client.go:51-67`
Purpose: Demonstrate how the consumer enforces contract assumptions at runtime — failing when `customer.name` is missing or `status` has unknown value.

```go
	var raw struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		Customer struct {
			Name string `json:"name"`
		} `json:"customer"`
		Total int64 `json:"total"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse contract schema: %w", err)
	}

	if raw.Customer.Name == "" {
		return nil, fmt.Errorf("contract violation: customer.name is missing")
	}

	if raw.Status != "IN_PROGRESS" && raw.Status != "COMPLETED" {
		return nil, fmt.Errorf("contract violation: unknown status %q", raw.Status)
	}
```

Explanation: `json.Unmarshal` into a struct with `json:"name"` will silently default to `""` when the field is absent (`full_name` instead of `name`). The explicit check then surfaces the contract violation with a clear message — mirroring how a real mobile client would surface schema breakage.

---

## Snippet 6 — Verifier Core Engine (Verify method)

Source File: `internal/contract/verifier.go:57-115`
Purpose: Execute every interaction of a contract against a live provider base URL; collect precise diffs per field path.

```go
func (v *Verifier) Verify(baseURL string, c *Contract) VerificationResult {
	result := VerificationResult{Passed: true}

	for _, interaction := range c.Interactions {
		targetURL := strings.TrimRight(baseURL, "/") + interaction.Request.Path
		req, err := http.NewRequest(interaction.Request.Method, targetURL, nil)
		if err != nil {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("[%s] failed to build request: %v", interaction.Description, err))
			continue
		}

		for k, val := range interaction.Request.Headers {
			req.Header.Set(k, val)
		}

		resp, err := v.Client.Do(req)
		if err != nil {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("[%s] request failed: %v", interaction.Description, err))
			continue
		}

		if resp.StatusCode != interaction.Response.Status {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("[%s] status code mismatch: expected %d, got %d",
				interaction.Description, interaction.Response.Status, resp.StatusCode))
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("[%s] failed to read body: %v", interaction.Description, err))
			continue
		}

		var actualBody map[string]interface{}
		decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
		decoder.UseNumber()
		if err := decoder.Decode(&actualBody); err != nil {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("[%s] response is not valid JSON object: %v", interaction.Description, err))
			continue
		}

		// Validate minimal subset expected by consumer
		diffs := diffValues("", interaction.Response.Body, actualBody)
		if len(diffs) > 0 {
			result.Passed = false
			for _, d := range diffs {
				result.Errors = append(result.Errors, fmt.Sprintf("[%s] %s", interaction.Description, d))
			}
		}
	}

	return result
}
```

Explanation: The verifier performs **subset verification** — it walks every field declared in the contract's `Body` and recurses only into those. Provider-side fields absent from the contract are never inspected.

---

## Snippet 7 — Verifier Recursive Diff Engine

Source File: `internal/contract/verifier.go:117-176`
Purpose: Compare expected contract body against provider body field-by-field, distinguishing missing fields, type mismatches (int vs string via `json.Number`), and value mismatches (enum casing).

```go
func diffValues(path string, expected, actual interface{}) []string {
	var diffs []string

	if expected == nil {
		return diffs
	}

	expMap, isExpMap := expected.(map[string]interface{})
	actMap, isActMap := actual.(map[string]interface{})

	if isExpMap {
		if !isActMap {
			return []string{fmt.Sprintf("path '%s': expected object, got %T", path, actual)}
		}
		for key, expVal := range expMap {
			childPath := key
			if path != "" {
				childPath = path + "." + key
			}
			actVal, exists := actMap[key]
			if !exists {
				diffs = append(diffs, fmt.Sprintf("missing expected field '%s'", childPath))
				continue
			}
			diffs = append(diffs, diffValues(childPath, expVal, actVal)...)
		}
		return diffs
	}

	// Compare primitives
	expNum, isExpNum := expected.(json.Number)
	actNum, isActNum := actual.(json.Number)

	if isExpNum || isActNum {
		if !(isExpNum && isActNum) {
			diffs = append(diffs, fmt.Sprintf("path '%s': type mismatch (expected %v [%T], got %v [%T])",
				path, expected, expected, actual, actual))
			return diffs
		}
		if expNum.String() != actNum.String() {
			diffs = append(diffs, fmt.Sprintf("path '%s': value mismatch (expected %s, got %s)",
				path, expNum.String(), actNum.String()))
		}
		return diffs
	}

	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		diffs = append(diffs, fmt.Sprintf("path '%s': type mismatch (expected %v [%T], got %v [%T])",
			path, expected, expected, actual, actual))
		return diffs
	}

	if expected != actual {
		diffs = append(diffs, fmt.Sprintf("path '%s': value mismatch (expected %q, got %q)",
			path, expected, actual))
	}

	return diffs
}
```

Explanation: The critical design decision is `decoder.UseNumber()`, which causes all JSON numbers to decode as `json.Number` instead of `float64`. When `total` arrives as the string `"150000"` (provider `OrderResponseBreaking.Total` is `string`), `actual` is a `string` not a `json.Number`, producing `type mismatch (expected 150000 [json.Number], got 150000 [string])`. This captures the semantic type change that a plain schema validator would miss.

---

## Snippet 8 — Model DTOs (V1 / V2 / Breaking)

Source File: `internal/model/order.go`
Purpose: Define the three provider response DTO representations as concrete Go types, each annotated with JSON tags encoding the breaking/migrating/evolved variants.

```go
// OrderResponseV1 matches original V1 consumer contract expectations.
type OrderResponseV1 struct {
	ID       string             `json:"id"`
	Status   string             `json:"status"` // "IN_PROGRESS", "COMPLETED"
	Customer CustomerResponseV1 `json:"customer"`
	Total    int64              `json:"total"` // Integer amount in IDR
	Notes    string             `json:"notes,omitempty"`
}

type CustomerResponseV1 struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// OrderResponseBreaking represents provider modifications that violate contract.
type OrderResponseBreaking struct {
	ID       string                   `json:"id"`
	Status   string                   `json:"status"` // "in_progress" (BREAKING: enum casing)
	Customer CustomerResponseBreaking `json:"customer"`
	Total    string                   `json:"total"` // "150000" (BREAKING: primitive type int->str)
	Notes    string                   `json:"notes,omitempty"`
}

type CustomerResponseBreaking struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"` // (BREAKING: renamed name -> full_name)
}

// OrderResponseV2 represents safe evolutionary schema.
type OrderResponseV2 struct {
	ID       string             `json:"id"`
	Status   string             `json:"status"` // "in_progress"
	Customer CustomerResponseV2 `json:"customer"`
	Total    string             `json:"total"` // "150000"
	Currency string             `json:"currency"`
}

type CustomerResponseV2 struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}
```

Explanation: The type definitions are the *source of truth* for what each provider version emits. `OrderResponseBreaking.Total` is `string` — this is the breaking mutation. `OrderResponseV2` adds `Currency` and uses `full_name`/`in_progress` — the evolved schema exposed on a separate `/v2` path.

---

## Snippet 9 — Demo Orchestrator (CI Gate Stages)

Source File: `cmd/demo/main.go:14-68`
Purpose: Drive the four-stage demo that illustrates the contract testing lifecycle and CI gate (allow/block) decisions.

```go
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
```

Explanation: Each `verifier.Verify()` call is treated as a CI gate decision: `Passed=true` → `ALLOWED`; `Passed=false` → errors enumerated and deployment `BLOCKED`. The demo exits non-zero (`os.Exit(1)`) on unexpected results, enabling pipeline failure integration.
