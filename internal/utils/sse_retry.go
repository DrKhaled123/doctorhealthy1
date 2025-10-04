package utils

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	"api-key-generator/internal/models"
)

// SSERetryConfig holds configuration for SSE retry logic
type SSERetryConfig struct {
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	BackoffFactor     float64
	MaxRetries        int
	JitterPercent     float64
	RetryableErrors   []string
	HeartbeatInterval time.Duration
}

// DefaultSSERetryConfig returns a sensible default configuration
func DefaultSSERetryConfig() *SSERetryConfig {
	return &SSERetryConfig{
		InitialDelay:      1 * time.Second,
		MaxDelay:          30 * time.Second,
		BackoffFactor:     2.0,
		MaxRetries:        10,
		JitterPercent:     0.1, // 10% jitter
		RetryableErrors:   []string{"network_error", "timeout", "server_error", "temporary_failure"},
		HeartbeatInterval: 30 * time.Second,
	}
}

// SSERetryHandler manages SSE connections with retry logic
type SSERetryHandler struct {
	config      *SSERetryConfig
	url         string
	eventSource *EventSource
	retryCount  int
	isConnected bool
	ctx         context.Context
	cancel      context.CancelFunc
}

// EventSource represents a simplified SSE client
type EventSource struct {
	URL         string
	Headers     map[string]string
	LastEventID string
	IsConnected bool
}

// NewSSERetryHandler creates a new SSE retry handler
func NewSSERetryHandler(url string, config *SSERetryConfig) *SSERetryHandler {
	if config == nil {
		config = DefaultSSERetryConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &SSERetryHandler{
		config:      config,
		url:         url,
		retryCount:  0,
		isConnected: false,
		ctx:         ctx,
		cancel:      cancel,
		eventSource: &EventSource{
			URL:     url,
			Headers: make(map[string]string),
		},
	}
}

// Connect establishes SSE connection with retry logic
func (h *SSERetryHandler) Connect() error {
	return h.connectWithRetry()
}

// Disconnect closes the SSE connection
func (h *SSERetryHandler) Disconnect() {
	h.cancel()
	if h.eventSource != nil {
		h.eventSource.IsConnected = false
	}
}

// IsConnected returns the connection status
func (h *SSERetryHandler) IsConnected() bool {
	return h.isConnected
}

// SetHeader sets a header for the SSE request
func (h *SSERetryHandler) SetHeader(key, value string) {
	h.eventSource.Headers[key] = value
}

// GetLastEventID returns the last event ID
func (h *SSERetryHandler) GetLastEventID() string {
	return h.eventSource.LastEventID
}

// SetLastEventID sets the last event ID
func (h *SSERetryHandler) SetLastEventID(eventID string) {
	h.eventSource.LastEventID = eventID
}

// connectWithRetry attempts to connect with exponential backoff retry logic
func (h *SSERetryHandler) connectWithRetry() error {
	var err error
	for h.retryCount = 0; h.retryCount <= h.config.MaxRetries; h.retryCount++ {
		// Check if context was cancelled
		select {
		case <-h.ctx.Done():
			return h.ctx.Err()
		default:
		}

		// Attempt to connect
		if err = h.attemptConnection(); err == nil {
			h.isConnected = true
			log.Printf("SSE connected successfully to %s", h.url)
			return nil
		}

		// Log retry attempt
		log.Printf("SSE connection attempt %d/%d failed: %v", h.retryCount+1, h.config.MaxRetries+1, err)

		// Check if error is retryable
		if !h.isRetryableError(err) {
			return fmt.Errorf("non-retryable error: %w", err)
		}

		// Wait before retry
		if h.retryCount < h.config.MaxRetries {
			delay := h.calculateDelay()
			log.Printf("Retrying SSE connection in %v", delay)

			select {
			case <-time.After(delay):
				// Continue to next retry
			case <-h.ctx.Done():
				return h.ctx.Err()
			}
		}
	}

	return models.NewSecurityError(
		fmt.Sprintf("Failed to connect after %d attempts", h.config.MaxRetries+1),
		"sse_connection_failed",
		"high",
	)
}

// attemptConnection tries to establish a single SSE connection
func (h *SSERetryHandler) attemptConnection() error {
	// Create HTTP request
	req, err := http.NewRequestWithContext(h.ctx, "GET", h.url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for key, value := range h.eventSource.Headers {
		req.Header.Set(key, value)
	}

	// Add Last-Event-ID if available
	if h.eventSource.LastEventID != "" {
		req.Header.Set("Last-Event-ID", h.eventSource.LastEventID)
	}

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	if contentType != "text/event-stream" {
		return fmt.Errorf("unexpected content type: %s", contentType)
	}

	// Connection successful
	h.eventSource.IsConnected = true
	return nil
}

// isRetryableError determines if an error should trigger a retry
func (h *SSERetryHandler) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()

	// Check against configured retryable errors
	for _, retryable := range h.config.RetryableErrors {
		if contains(errMsg, retryable) {
			return true
		}
	}

	// Check for common network-related errors
	networkErrors := []string{
		"connection refused",
		"connection reset",
		"timeout",
		"network is unreachable",
		"no such host",
		"server misbehaving",
	}

	for _, networkErr := range networkErrors {
		if contains(errMsg, networkErr) {
			return true
		}
	}

	return false
}

// calculateDelay computes the delay for the next retry with exponential backoff and jitter
func (h *SSERetryHandler) calculateDelay() time.Duration {
	// Exponential backoff
	delay := float64(h.config.InitialDelay)
	delay *= math.Pow(h.config.BackoffFactor, float64(h.retryCount))

	// Cap at max delay
	if delay > float64(h.config.MaxDelay) {
		delay = float64(h.config.MaxDelay)
	}

	// Add jitter
	jitterRange := delay * h.config.JitterPercent
	jitter := rand.Float64() * jitterRange
	actualDelay := time.Duration(delay + jitter)

	return actualDelay
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsAt(s, substr))))
}

// containsAt checks if substring exists at any position in string
func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// SSEWithRetry creates and starts an SSE connection with retry logic
func SSEWithRetry(url string, config *SSERetryConfig) (*SSERetryHandler, error) {
	handler := NewSSERetryHandler(url, config)

	if err := handler.Connect(); err != nil {
		return nil, err
	}

	return handler, nil
}

// MonitorSSEHealth provides health monitoring for SSE connections
func (h *SSERetryHandler) MonitorSSEHealth() {
	ticker := time.NewTicker(h.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check connection health
			if !h.isConnected {
				log.Println("SSE connection lost, attempting to reconnect...")
				if err := h.connectWithRetry(); err != nil {
					log.Printf("SSE reconnection failed: %v", err)
				}
			}
		case <-h.ctx.Done():
			return
		}
	}
}

// StartHealthMonitoring starts the health monitoring goroutine
func (h *SSERetryHandler) StartHealthMonitoring() {
	go h.MonitorSSEHealth()
}

// SSEConnectionManager manages multiple SSE connections
type SSEConnectionManager struct {
	connections map[string]*SSERetryHandler
	config      *SSERetryConfig
}

// NewSSEConnectionManager creates a new connection manager
func NewSSEConnectionManager(config *SSERetryConfig) *SSEConnectionManager {
	return &SSEConnectionManager{
		connections: make(map[string]*SSERetryHandler),
		config:      config,
	}
}

// AddConnection adds a new SSE connection
func (m *SSEConnectionManager) AddConnection(name, url string) error {
	handler := NewSSERetryHandler(url, m.config)

	if err := handler.Connect(); err != nil {
		return fmt.Errorf("failed to connect to %s: %w", url, err)
	}

	handler.StartHealthMonitoring()
	m.connections[name] = handler

	log.Printf("Added SSE connection: %s -> %s", name, url)
	return nil
}

// RemoveConnection removes an SSE connection
func (m *SSEConnectionManager) RemoveConnection(name string) {
	if handler, exists := m.connections[name]; exists {
		handler.Disconnect()
		delete(m.connections, name)
		log.Printf("Removed SSE connection: %s", name)
	}
}

// GetConnection retrieves a connection by name
func (m *SSEConnectionManager) GetConnection(name string) (*SSERetryHandler, bool) {
	handler, exists := m.connections[name]
	return handler, exists
}

// DisconnectAll closes all connections
func (m *SSEConnectionManager) DisconnectAll() {
	for name, handler := range m.connections {
		handler.Disconnect()
		log.Printf("Disconnected SSE connection: %s", name)
	}
	m.connections = make(map[string]*SSERetryHandler)
}

// GetConnectionCount returns the number of active connections
func (m *SSEConnectionManager) GetConnectionCount() int {
	return len(m.connections)
}
