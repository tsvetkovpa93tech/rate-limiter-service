package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Event represents a webhook event
type Event struct {
	Type      string    `json:"type"`      // "blocked", "limit_exceeded"
	Key       string    `json:"key"`
	Algorithm string    `json:"algorithm"`
	Limit     int       `json:"limit"`
	Window    string    `json:"window"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

// Client sends webhook notifications
type Client struct {
	url     string
	timeout time.Duration
	client  *http.Client
	logger  *slog.Logger
}

// NewClient creates a new webhook client
func NewClient(url string, timeout time.Duration, logger *slog.Logger) *Client {
	if logger == nil {
		logger = slog.Default()
	}
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	return &Client{
		url:     url,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

// Send sends a webhook event
func (c *Client) Send(ctx context.Context, event Event) error {
	if c.url == "" {
		return nil // Webhook not configured
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "rate-limiter-service/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Warn("Failed to send webhook", "error", err, "url", c.url)
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		c.logger.Warn("Webhook returned error status", "status", resp.StatusCode, "url", c.url)
		return fmt.Errorf("webhook returned status %d", resp.StatusCode)
	}

	c.logger.Debug("Webhook sent successfully", "url", c.url, "type", event.Type)
	return nil
}

// SendAsync sends a webhook event asynchronously
func (c *Client) SendAsync(event Event) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
		defer cancel()

		if err := c.Send(ctx, event); err != nil {
			c.logger.Error("Async webhook send failed", "error", err)
		}
	}()
}

