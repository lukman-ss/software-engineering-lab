package observability

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// HTTPNotificationClient sends notification via HTTP and propagates W3C traceparent header.
type HTTPNotificationClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewHTTPNotificationClient(baseURL string, client *http.Client) *HTTPNotificationClient {
	if client == nil {
		client = &http.Client{
			Timeout:   5 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}
	} else if client.Transport == nil {
		client.Transport = otelhttp.NewTransport(http.DefaultTransport)
	} else {
		// Wrap existing transport if not already wrapped
		if _, ok := client.Transport.(*otelhttp.Transport); !ok {
			client.Transport = otelhttp.NewTransport(client.Transport)
		}
	}
	return &HTTPNotificationClient{
		BaseURL:    baseURL,
		HTTPClient: client,
	}
}

func (c *HTTPNotificationClient) Send(ctx context.Context, invoiceID string) error {
	reqURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("invalid base url: %w", err)
	}
	reqURL.Path = "/notifications"
	q := reqURL.Query()
	q.Set("invoice_id", invoiceID)
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("notification client request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("notification service returned %d", resp.StatusCode)
	}
	return nil
}
