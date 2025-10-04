package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"api-key-generator/internal/models"

	"github.com/labstack/echo/v4"
)

// MonitoringConfig holds monitoring service configuration
type MonitoringConfig struct {
	SentryDSN      string
	DataDogAPIKey  string
	GrafanaURL     string
	EnableMetrics  bool
	EnableTracing  bool
	EnableLogging  bool
	Environment    string
	ServiceName    string
	ServiceVersion string
}

// DefaultMonitoringConfig returns a default monitoring configuration
func DefaultMonitoringConfig() *MonitoringConfig {
	return &MonitoringConfig{
		EnableMetrics:  true,
		EnableTracing:  true,
		EnableLogging:  true,
		Environment:    "development",
		ServiceName:    "nutrition-app",
		ServiceVersion: "1.0.0",
	}
}

// ErrorTracker handles error tracking and monitoring
type ErrorTracker struct {
	config     *MonitoringConfig
	httpClient *http.Client
}

// NewErrorTracker creates a new error tracker
func NewErrorTracker(config *MonitoringConfig) *ErrorTracker {
	if config == nil {
		config = DefaultMonitoringConfig()
	}

	return &ErrorTracker{
		config: config,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// TrackError sends error information to monitoring services
func (et *ErrorTracker) TrackError(errorResp *models.EnhancedErrorResponse, context map[string]interface{}) {
	// Track in Sentry if configured
	if et.config.SentryDSN != "" {
		et.trackInSentry(errorResp, context)
	}

	// Track in DataDog if configured
	if et.config.DataDogAPIKey != "" {
		et.trackInDataDog(errorResp, context)
	}

	// Log locally
	et.logLocally(errorResp, context)
}

// trackInSentry sends error to Sentry
func (et *ErrorTracker) trackInSentry(errorResp *models.EnhancedErrorResponse, context map[string]interface{}) {
	// This is a simplified implementation
	// In practice, you would use the official Sentry Go SDK

	sentryEvent := map[string]interface{}{
		"exception": map[string]interface{}{
			"values": []map[string]interface{}{
				{
					"type":       errorResp.Error,
					"value":      errorResp.Message,
					"stacktrace": context["stacktrace"],
				},
			},
		},
		"level":       et.getSentryLevel(errorResp.Severity),
		"tags":        et.buildSentryTags(errorResp),
		"contexts":    et.buildSentryContexts(context),
		"release":     et.config.ServiceVersion,
		"environment": et.config.Environment,
	}

	et.sendToService("sentry", sentryEvent)
}

// trackInDataDog sends error to DataDog
func (et *ErrorTracker) trackInDataDog(errorResp *models.EnhancedErrorResponse, context map[string]interface{}) {
	dataDogEvent := map[string]interface{}{
		"title":      fmt.Sprintf("%s: %s", errorResp.Error, errorResp.Message),
		"text":       errorResp.Message,
		"alert_type": et.getDataDogAlertType(errorResp.Severity),
		"tags": []string{
			fmt.Sprintf("service:%s", et.config.ServiceName),
			fmt.Sprintf("version:%s", et.config.ServiceVersion),
			fmt.Sprintf("environment:%s", et.config.Environment),
			fmt.Sprintf("error_type:%s", errorResp.Error),
			fmt.Sprintf("severity:%s", errorResp.Severity),
		},
		"timestamp": time.Now().Unix(),
	}

	et.sendToService("datadog", dataDogEvent)
}

// logLocally logs error information locally
func (et *ErrorTracker) logLocally(errorResp *models.EnhancedErrorResponse, context map[string]interface{}) {
	log.Printf("ERROR TRACKED [%s]: %s | Message: %s | Severity: %s | TraceID: %s",
		et.config.ServiceName,
		errorResp.Error,
		errorResp.Message,
		errorResp.Severity,
		et.getTraceID(context),
	)
}

// getSentryLevel maps severity to Sentry level
func (et *ErrorTracker) getSentryLevel(severity string) string {
	switch severity {
	case "critical":
		return "fatal"
	case "high":
		return "error"
	case "medium":
		return "warning"
	case "low":
		return "info"
	default:
		return "error"
	}
}

// getDataDogAlertType maps severity to DataDog alert type
func (et *ErrorTracker) getDataDogAlertType(severity string) string {
	switch severity {
	case "critical", "high":
		return "error"
	case "medium":
		return "warning"
	case "low":
		return "info"
	default:
		return "error"
	}
}

// buildSentryTags builds Sentry tags from error response
func (et *ErrorTracker) buildSentryTags(errorResp *models.EnhancedErrorResponse) map[string]string {
	tags := map[string]string{
		"service":     et.config.ServiceName,
		"version":     et.config.ServiceVersion,
		"environment": et.config.Environment,
		"error_type":  errorResp.Error,
		"severity":    errorResp.Severity,
		"category":    errorResp.Category,
	}

	if errorResp.Context != nil {
		if errorResp.Context.UserID != "" {
			tags["user_id"] = errorResp.Context.UserID
		}
		if errorResp.Context.Endpoint != "" {
			tags["endpoint"] = errorResp.Context.Endpoint
		}
	}

	return tags
}

// buildSentryContexts builds Sentry contexts
func (et *ErrorTracker) buildSentryContexts(context map[string]interface{}) map[string]interface{} {
	contexts := map[string]interface{}{
		"service": map[string]interface{}{
			"name":    et.config.ServiceName,
			"version": et.config.ServiceVersion,
		},
	}

	// Add custom context
	for key, value := range context {
		contexts[key] = value
	}

	return contexts
}

// getTraceID extracts trace ID from context
func (et *ErrorTracker) getTraceID(context map[string]interface{}) string {
	if traceID, ok := context["trace_id"].(string); ok {
		return traceID
	}
	return ""
}

// sendToService sends data to external monitoring service
func (et *ErrorTracker) sendToService(service string, data interface{}) {
	// This is a simplified implementation
	// In practice, you would use the appropriate SDK for each service

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal %s data: %v", service, err)
		return
	}

	// Log the data that would be sent
	log.Printf("Would send to %s: %s", service, string(jsonData))
}

// MetricsCollector collects and sends application metrics
type MetricsCollector struct {
	config     *MonitoringConfig
	metrics    map[string]interface{}
	httpClient *http.Client
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector(config *MonitoringConfig) *MetricsCollector {
	if config == nil {
		config = DefaultMonitoringConfig()
	}

	return &MetricsCollector{
		config:  config,
		metrics: make(map[string]interface{}),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RecordMetric records a metric value
func (mc *MetricsCollector) RecordMetric(name string, value interface{}, tags map[string]string) {
	mc.metrics[name] = map[string]interface{}{
		"value":     value,
		"timestamp": time.Now().Unix(),
		"tags":      tags,
	}
}

// IncrementCounter increments a counter metric
func (mc *MetricsCollector) IncrementCounter(name string, tags map[string]string) {
	current, exists := mc.metrics[name]
	if !exists {
		mc.metrics[name] = map[string]interface{}{
			"value":     1,
			"timestamp": time.Now().Unix(),
			"tags":      tags,
		}
	} else {
		if metric, ok := current.(map[string]interface{}); ok {
			if val, ok := metric["value"].(int); ok {
				metric["value"] = val + 1
				metric["timestamp"] = time.Now().Unix()
			}
		}
	}
}

// RecordHistogram records a histogram value
func (mc *MetricsCollector) RecordHistogram(name string, value float64, tags map[string]string) {
	metric := mc.metrics[name]
	if metric == nil {
		mc.metrics[name] = map[string]interface{}{
			"values":    []float64{value},
			"timestamp": time.Now().Unix(),
			"tags":      tags,
		}
	} else {
		if m, ok := metric.(map[string]interface{}); ok {
			if values, ok := m["values"].([]float64); ok {
				m["values"] = append(values, value)
				m["timestamp"] = time.Now().Unix()
			}
		}
	}
}

// Flush sends all collected metrics to monitoring services
func (mc *MetricsCollector) Flush() {
	if len(mc.metrics) == 0 {
		return
	}

	// Send to DataDog if configured
	if mc.config.DataDogAPIKey != "" {
		mc.sendToDataDog()
	}

	// Send to Grafana if configured
	if mc.config.GrafanaURL != "" {
		mc.sendToGrafana()
	}

	// Clear metrics after sending
	mc.metrics = make(map[string]interface{})
}

// sendToDataDog sends metrics to DataDog
func (mc *MetricsCollector) sendToDataDog() {
	// Convert metrics to DataDog format
	series := make([]map[string]interface{}, 0)

	for name, metric := range mc.metrics {
		if m, ok := metric.(map[string]interface{}); ok {
			series = append(series, map[string]interface{}{
				"metric": fmt.Sprintf("%s.%s", mc.config.ServiceName, name),
				"points": [][]interface{}{
					{m["timestamp"], m["value"]},
				},
				"tags": mc.buildDataDogTags(m["tags"]),
			})
		}
	}

	payload := map[string]interface{}{
		"series": series,
	}

	mc.sendToService("datadog", payload)
}

// sendToGrafana sends metrics to Grafana
func (mc *MetricsCollector) sendToGrafana() {
	// Convert metrics to Grafana format
	metrics := make([]map[string]interface{}, 0)

	for name, metric := range mc.metrics {
		if m, ok := metric.(map[string]interface{}); ok {
			metrics = append(metrics, map[string]interface{}{
				"name":      fmt.Sprintf("%s_%s", mc.config.ServiceName, name),
				"value":     m["value"],
				"timestamp": m["timestamp"],
				"tags":      m["tags"],
			})
		}
	}

	payload := map[string]interface{}{
		"metrics": metrics,
	}

	mc.sendToService("grafana", payload)
}

// buildDataDogTags builds DataDog-compatible tags
func (mc *MetricsCollector) buildDataDogTags(tags interface{}) []string {
	if tags == nil {
		return []string{}
	}

	if tagMap, ok := tags.(map[string]string); ok {
		result := make([]string, 0, len(tagMap))
		for key, value := range tagMap {
			result = append(result, fmt.Sprintf("%s:%s", key, value))
		}
		return result
	}

	return []string{}
}

// sendToService sends data to external service
func (mc *MetricsCollector) sendToService(service string, data interface{}) {
	// This is a simplified implementation
	// In practice, you would use the appropriate client for each service

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal %s metrics: %v", service, err)
		return
	}

	// Log the data that would be sent
	log.Printf("Would send metrics to %s: %s", service, string(jsonData))
}

// TracingService handles distributed tracing
type TracingService struct {
	config     *MonitoringConfig
	traces     []Trace
	httpClient *http.Client
}

// Trace represents a distributed trace
type Trace struct {
	TraceID   string                 `json:"trace_id"`
	SpanID    string                 `json:"span_id"`
	Operation string                 `json:"operation"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
	Duration  time.Duration          `json:"duration"`
	Tags      map[string]string      `json:"tags"`
	Context   map[string]interface{} `json:"context"`
}

// NewTracingService creates a new tracing service
func NewTracingService(config *MonitoringConfig) *TracingService {
	if config == nil {
		config = DefaultMonitoringConfig()
	}

	return &TracingService{
		config: config,
		traces: make([]Trace, 0),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// StartSpan starts a new trace span
func (ts *TracingService) StartSpan(operation string, tags map[string]string) *Span {
	span := &Span{
		TraceID:   generateTraceID(),
		SpanID:    generateSpanID(),
		Operation: operation,
		StartTime: time.Now(),
		Tags:      tags,
		Context:   make(map[string]interface{}),
		tracing:   ts,
	}

	return span
}

// Span represents an active trace span
type Span struct {
	TraceID   string                 `json:"trace_id"`
	SpanID    string                 `json:"span_id"`
	Operation string                 `json:"operation"`
	StartTime time.Time              `json:"start_time"`
	EndTime   time.Time              `json:"end_time"`
	Duration  time.Duration          `json:"duration"`
	Tags      map[string]string      `json:"tags"`
	Context   map[string]interface{} `json:"context"`
	tracing   *TracingService
}

// Finish completes the span and records the trace
func (s *Span) Finish() {
	s.EndTime = time.Now()
	s.Duration = s.EndTime.Sub(s.StartTime)

	// Record the trace
	s.tracing.traces = append(s.tracing.traces, Trace{
		TraceID:   s.TraceID,
		SpanID:    s.SpanID,
		Operation: s.Operation,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
		Duration:  s.Duration,
		Tags:      s.Tags,
		Context:   s.Context,
	})

	// Send to tracing service if configured
	s.sendToTracingService()
}

// SetTag sets a tag on the span
func (s *Span) SetTag(key, value string) {
	if s.Tags == nil {
		s.Tags = make(map[string]string)
	}
	s.Tags[key] = value
}

// SetContext sets context information on the span
func (s *Span) SetContext(key string, value interface{}) {
	s.Context[key] = value
}

// sendToTracingService sends trace data to external tracing service
func (s *Span) sendToTracingService() {
	// This is a simplified implementation
	// In practice, you would use a tracing service like Jaeger or DataDog APM

	traceData := map[string]interface{}{
		"trace_id":   s.TraceID,
		"span_id":    s.SpanID,
		"operation":  s.Operation,
		"duration":   s.Duration.Microseconds(),
		"start_time": s.StartTime.UnixMicro(),
		"end_time":   s.EndTime.UnixMicro(),
		"tags":       s.Tags,
		"context":    s.Context,
	}

	jsonData, err := json.Marshal(traceData)
	if err != nil {
		log.Printf("Failed to marshal trace data: %v", err)
		return
	}

	log.Printf("Would send trace to tracing service: %s", string(jsonData))
}

// generateTraceID generates a unique trace ID
func generateTraceID() string {
	return fmt.Sprintf("trace_%d", time.Now().UnixNano())
}

// generateSpanID generates a unique span ID
func generateSpanID() string {
	return fmt.Sprintf("span_%d", time.Now().UnixNano())
}

// Global monitoring instances
var GlobalErrorTracker = NewErrorTracker(DefaultMonitoringConfig())
var GlobalMetricsCollector = NewMetricsCollector(DefaultMonitoringConfig())
var GlobalTracingService = NewTracingService(DefaultMonitoringConfig())

// Convenience functions for global monitoring

// TrackError tracks an error using the global tracker
func TrackError(errorResp *models.EnhancedErrorResponse, context map[string]interface{}) {
	GlobalErrorTracker.TrackError(errorResp, context)
}

// RecordMetric records a metric using the global collector
func RecordMetric(name string, value interface{}, tags map[string]string) {
	GlobalMetricsCollector.RecordMetric(name, value, tags)
}

// IncrementCounter increments a counter using the global collector
func IncrementCounter(name string, tags map[string]string) {
	GlobalMetricsCollector.IncrementCounter(name, tags)
}

// RecordHistogram records a histogram value using the global collector
func RecordHistogram(name string, value float64, tags map[string]string) {
	GlobalMetricsCollector.RecordHistogram(name, value, tags)
}

// StartSpan starts a trace span using the global tracing service
func StartSpan(operation string, tags map[string]string) *Span {
	return GlobalTracingService.StartSpan(operation, tags)
}

// FlushMetrics flushes all collected metrics
func FlushMetrics() {
	GlobalMetricsCollector.Flush()
}

// MonitoringMiddleware provides HTTP middleware for automatic monitoring
type MonitoringMiddleware struct {
	errorTracker *ErrorTracker
	metrics      *MetricsCollector
	tracing      *TracingService
}

// NewMonitoringMiddleware creates new monitoring middleware
func NewMonitoringMiddleware() *MonitoringMiddleware {
	return &MonitoringMiddleware{
		errorTracker: GlobalErrorTracker,
		metrics:      GlobalMetricsCollector,
		tracing:      GlobalTracingService,
	}
}

// Middleware returns Echo middleware for monitoring
func (mm *MonitoringMiddleware) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Start trace span
			span := mm.tracing.StartSpan(
				fmt.Sprintf("%s %s", c.Request().Method, c.Request().URL.Path),
				map[string]string{
					"method":     c.Request().Method,
					"endpoint":   c.Request().URL.Path,
					"user_agent": c.Request().UserAgent(),
				},
			)

			// Add request context
			span.SetContext("request_id", c.Get("request_id"))
			span.SetContext("user_id", c.Get("user_id"))

			// Record request start
			start := time.Now()
			mm.metrics.IncrementCounter("http_requests_total", map[string]string{
				"method":   c.Request().Method,
				"endpoint": c.Request().URL.Path,
			})

			// Process request
			err := next(c)

			// Record response metrics
			duration := time.Since(start)
			mm.metrics.RecordHistogram("http_request_duration_seconds", duration.Seconds(), map[string]string{
				"method":   c.Request().Method,
				"endpoint": c.Request().URL.Path,
				"status":   fmt.Sprintf("%d", c.Response().Status),
			})

			// Track errors
			if err != nil {
				errorResp := models.NewEnhancedErrorResponse(
					"HTTP Request Error",
					err.Error(),
					"http_error",
					false,
				).WithContext("method", c.Request().Method).
					WithContext("endpoint", c.Request().URL.Path).
					WithContext("status_code", c.Response().Status).
					WithCategory("http").
					WithSeverity("medium")

				mm.errorTracker.TrackError(errorResp, map[string]interface{}{
					"method":     c.Request().Method,
					"endpoint":   c.Request().URL.Path,
					"status":     c.Response().Status,
					"user_agent": c.Request().UserAgent(),
				})
			}

			// Finish span
			span.Finish()

			return err
		}
	}
}

// Global monitoring middleware instance
var GlobalMonitoringMiddleware = NewMonitoringMiddleware()

// GetMonitoringMiddleware returns the global monitoring middleware
func GetMonitoringMiddleware() *MonitoringMiddleware {
	return GlobalMonitoringMiddleware
}
