package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// HTTPNotificationClient sends notification via HTTP and propagates W3C traceparent header.
type HTTPNotificationClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewHTTPNotificationClient(baseURL string, client *http.Client) *HTTPNotificationClient {
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Second}
	}
	return &HTTPNotificationClient{
		BaseURL:    baseURL,
		HTTPClient: client,
	}
}

func (c *HTTPNotificationClient) Send(ctx context.Context, invoiceID string) error {
	url := fmt.Sprintf("%s/notifications?invoice_id=%s", c.BaseURL, invoiceID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	// Inject W3C trace context using official OpenTelemetry TextMapPropagator
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("notification service returned %d", resp.StatusCode)
	}
	return nil
}
