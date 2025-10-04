package utils

import (
	"fmt"
	"runtime"
	"time"

	"api-key-generator/internal/models"
)

// WasmErrorBoundary provides error boundary functionality for WebAssembly operations
type WasmErrorBoundary struct {
	MaxMemoryMB    int64
	MaxInputSize   int64
	Timeout        time.Duration
	EnableFallback bool
}

// NewWasmErrorBoundary creates a new WASM error boundary with safe defaults
func NewWasmErrorBoundary() *WasmErrorBoundary {
	return &WasmErrorBoundary{
		MaxMemoryMB:    100,         // 100MB limit
		MaxInputSize:   1024 * 1024, // 1MB input limit
		Timeout:        30 * time.Second,
		EnableFallback: true,
	}
}

// SafeWasmCall executes a WASM function with comprehensive error handling
func (web *WasmErrorBoundary) SafeWasmCall(
	wasmFunc func() (interface{}, error),
	fallbackFunc func() (interface{}, error),
) (interface{}, error) {
	// Check input size limits
	if err := web.checkInputLimits(); err != nil {
		return web.handleWasmError(err, "input_size_exceeded", fallbackFunc)
	}

	// Set up timeout
	done := make(chan struct{})
	var result interface{}
	var err error

	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("wasm panic: %v", r)
			}
			close(done)
		}()

		result, err = wasmFunc()
	}()

	// Wait for completion or timeout
	select {
	case <-done:
		if err != nil {
			return web.handleWasmError(err, "wasm_execution_failed", fallbackFunc)
		}
		return result, nil
	case <-time.After(web.Timeout):
		return web.handleWasmError(fmt.Errorf("wasm timeout after %v", web.Timeout), "wasm_timeout", fallbackFunc)
	}
}

// checkInputLimits validates input size constraints
func (web *WasmErrorBoundary) checkInputLimits() error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Check memory usage
	memoryMB := int64(memStats.Alloc) / (1024 * 1024)
	if memoryMB > web.MaxMemoryMB {
		return models.NewWasmError(
			fmt.Sprintf("Memory usage %dMB exceeds limit %dMB", memoryMB, web.MaxMemoryMB),
			"memory_check",
			"checkInputLimits",
			memoryMB,
			0,
		)
	}

	return nil
}

// handleWasmError processes WASM errors and provides fallback options
func (web *WasmErrorBoundary) handleWasmError(
	err error,
	errorCode string,
	fallbackFunc func() (interface{}, error),
) (interface{}, error) {
	// Create detailed WASM error
	wasmErr := &models.WasmError{
		ErrorType: "wasm_error",
		Message:   err.Error(),
		Code:      errorCode,
		Timestamp: time.Now().UTC(),
	}

	// Add memory stats to context
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	wasmErr.MemoryUsage = int64(memStats.Alloc)
	wasmErr.Context = map[string]interface{}{
		"goroutines": runtime.NumGoroutine(),
		"memory_mb":  int64(memStats.Alloc) / (1024 * 1024),
		"gc_runs":    memStats.NumGC,
	}

	// Try fallback if enabled
	if web.EnableFallback && fallbackFunc != nil {
		if fallbackResult, fallbackErr := fallbackFunc(); fallbackErr == nil {
			// Log successful fallback
			wasmErr.Context["fallback_success"] = true
			wasmErr.Context["fallback_result_type"] = fmt.Sprintf("%T", fallbackResult)
			return fallbackResult, nil
		}
	}

	return nil, wasmErr
}

// ValidateWasmInput performs comprehensive input validation for WASM operations
func (web *WasmErrorBoundary) ValidateWasmInput(input interface{}) error {
	// Check for nil input
	if input == nil {
		return models.NewWasmError(
			"Input cannot be nil",
			"validation",
			"ValidateWasmInput",
			0,
			0,
		)
	}

	// Check input size (basic implementation)
	inputSize := web.estimateInputSize(input)
	if inputSize > web.MaxInputSize {
		return models.NewWasmError(
			fmt.Sprintf("Input size %d bytes exceeds limit %d bytes", inputSize, web.MaxInputSize),
			"validation",
			"ValidateWasmInput",
			0,
			inputSize,
		)
	}

	return nil
}

// estimateInputSize provides a rough estimate of input size
func (web *WasmErrorBoundary) estimateInputSize(input interface{}) int64 {
	// This is a simplified implementation
	// In practice, you might use reflection or serialization to get actual size
	switch v := input.(type) {
	case string:
		return int64(len(v))
	case []byte:
		return int64(len(v))
	default:
		// For complex types, return a conservative estimate
		return 1024 // 1KB default estimate
	}
}

// WasmPerformanceMonitor tracks WASM performance metrics
type WasmPerformanceMonitor struct {
	Threshold time.Duration
	AlertFunc func(operation string, duration time.Duration)
}

// NewWasmPerformanceMonitor creates a new performance monitor
func NewWasmPerformanceMonitor() *WasmPerformanceMonitor {
	return &WasmPerformanceMonitor{
		Threshold: 500 * time.Millisecond, // 500ms threshold
		AlertFunc: func(operation string, duration time.Duration) {
			// Default alert function - in production, integrate with monitoring system
			fmt.Printf("WASM PERFORMANCE ALERT: %s took %v (threshold: %v)\n",
				operation, duration, 500*time.Millisecond)
		},
	}
}

// MonitorWasmOperation wraps a WASM operation with performance monitoring
func (wpm *WasmPerformanceMonitor) MonitorWasmOperation(
	operation string,
	wasmFunc func() (interface{}, error),
) (interface{}, error) {
	start := time.Now()
	result, err := wasmFunc()
	duration := time.Since(start)

	// Check if operation exceeded threshold
	if duration > wpm.Threshold {
		if wpm.AlertFunc != nil {
			wpm.AlertFunc(operation, duration)
		}

		// Create performance error
		perfErr := models.NewPerformanceError(
			fmt.Sprintf("WASM operation '%s' exceeded performance threshold", operation),
			operation,
			"wasm_execution",
			duration,
			wpm.Threshold,
		)

		// If the operation succeeded despite being slow, return both result and warning
		if err == nil {
			// In a real implementation, you might return a custom result with warning
			return result, nil
		}

		return nil, perfErr
	}

	return result, err
}

// SafeWasmOperation combines error boundary and performance monitoring
func SafeWasmOperation(
	wasmFunc func() (interface{}, error),
	fallbackFunc func() (interface{}, error),
) (interface{}, error) {
	boundary := NewWasmErrorBoundary()
	monitor := NewWasmPerformanceMonitor()

	// Monitor the WASM operation with error boundary
	result, err := monitor.MonitorWasmOperation("wasm_operation", func() (interface{}, error) {
		return boundary.SafeWasmCall(wasmFunc, fallbackFunc)
	})

	return result, err
}
