package contract

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
)

// Interaction defines a single consumer-provider HTTP interaction expectation.
type Interaction struct {
	Description string            `json:"description"`
	ProviderState string          `json:"provider_state,omitempty"`
	Request     RequestDefinition `json:"request"`
	Response    ResponseDefinition `json:"response"`
}

type RequestDefinition struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers,omitempty"`
}

type ResponseDefinition struct {
	Status  int                    `json:"status"`
	Headers map[string]string      `json:"headers,omitempty"`
	Body    map[string]interface{} `json:"body"`
}

// Contract encapsulates all consumer expectations for a given provider.
type Contract struct {
	Consumer     string        `json:"consumer"`
	Provider     string        `json:"provider"`
	Interactions []Interaction `json:"interactions"`
}

// VerificationResult captures details of contract verification against a provider.
type VerificationResult struct {
	Passed bool     `json:"passed"`
	Errors []string `json:"errors"`
}

// Verifier executes contracts against a running HTTP server.
type Verifier struct {
	Client *http.Client
}

func NewVerifier() *Verifier {
	return &Verifier{
		Client: &http.Client{},
	}
}

// VerifyInteractions tests interactions against baseURL.
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
		// Both must be numbers
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
