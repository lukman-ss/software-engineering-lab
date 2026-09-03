package observability

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/lukman-ss/software-engineering-lab/pkg/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// HTTPNotificationClient sends notification via HTTP and propagates W3C traceparent header.
type HTTPNotificationClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewHTTPNotificationClient(baseURL string, client *http.Client) *HTTPNotificationClient {
	if client == nil {
		return &HTTPNotificationClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout:   5 * time.Second,
				Transport: otelhttp.NewTransport(http.DefaultTransport),
			},
		}
	}

	clientCopy := *client
	transport := http.DefaultTransport
	if clientCopy.Transport != nil {
		transport = clientCopy.Transport
	}
	if _, ok := transport.(*otelhttp.Transport); !ok {
		clientCopy.Transport = otelhttp.NewTransport(transport)
	}
	if clientCopy.Timeout == 0 {
		clientCopy.Timeout = 5 * time.Second
	}

	return &HTTPNotificationClient{
		BaseURL:    baseURL,
		HTTPClient: &clientCopy,
	}
}

func (c *HTTPNotificationClient) Send(ctx context.Context, invoiceID string) error {
	reqURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("invalid base url: %w", err)
	}

	// Preserve existing query params if any
	q := reqURL.Query()
	q.Set("invoice_id", invoiceID)
	reqURL.Path = "/notifications"
	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if reqID := middleware.GetRequestID(ctx); reqID != "" {
		req.Header.Set(middleware.HeaderRequestID, reqID)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("notification client request error: %w", err)
	}

	_, copyErr := io.Copy(io.Discard, io.LimitReader(resp.Body, 1024*1024))
	closeErr := resp.Body.Close()

	if resp.StatusCode >= 400 {
		statusErr := fmt.Errorf("notification service returned %d", resp.StatusCode)
		if copyErr != nil || closeErr != nil {
			return errors.Join(statusErr, copyErr, closeErr)
		}
		return statusErr
	}

	if copyErr != nil || closeErr != nil {
		return errors.Join(copyErr, closeErr)
	}

	return nil
}
