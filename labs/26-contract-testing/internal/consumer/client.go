package consumer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"labs/26-contract-testing/internal/contract"
)

// MobileOrderSummary is the consumer view of an Order (subset of provider payload).
type MobileOrderSummary struct {
	ID           string
	Status       string
	CustomerName string
	Total        int64
}

// MobileOrderClient interacts with the Order Provider service.
type MobileOrderClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewMobileOrderClient(baseURL string) *MobileOrderClient {
	return &MobileOrderClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// FetchOrder fetches order details and maps to consumer model.
func (c *MobileOrderClient) FetchOrder(orderID string) (*MobileOrderSummary, error) {
	url := fmt.Sprintf("%s/v1/orders/%s", c.BaseURL, orderID)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parsing raw JSON into consumer DTO
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

	return &MobileOrderSummary{
		ID:           raw.ID,
		Status:       raw.Status,
		CustomerName: raw.Customer.Name,
		Total:        raw.Total,
	}, nil
}

// GenerateMobileContract builds the consumer contract for order retrieval.
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
