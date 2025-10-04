package logger

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
)

type Logger struct {
	*log.Logger
	level  LogLevel
	output io.Writer
}

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	UserID    string                 `json:"user_id,omitempty"`
	IP        string                 `json:"ip,omitempty"`
	UserAgent string                 `json:"user_agent,omitempty"`
	Method    string                 `json:"method,omitempty"`
	URI       string                 `json:"uri,omitempty"`
	Status    int                    `json:"status,omitempty"`
	Latency   string                 `json:"latency,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

func New() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", 0),
		level:  INFO,
		output: os.Stdout,
	}
}

func NewWithLevel(level LogLevel) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", 0),
		level:  level,
		output: os.Stdout,
	}
}

func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

func (l *Logger) Debug(message string, fields ...interface{}) {
	if l.level <= DEBUG {
		l.log(DEBUG, message, fields...)
	}
}

func (l *Logger) Info(message string, fields ...interface{}) {
	if l.level <= INFO {
		l.log(INFO, message, fields...)
	}
}

func (l *Logger) Warn(message string, fields ...interface{}) {
	if l.level <= WARN {
		l.log(WARN, message, fields...)
	}
}

func (l *Logger) Error(message string, fields ...interface{}) {
	if l.level <= ERROR {
		l.log(ERROR, message, fields...)
	}
}

func (l *Logger) log(level LogLevel, message string, fields ...interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level.String(),
		Message:   message,
		Fields:    make(map[string]interface{}),
	}

	// Parse fields (key-value pairs)
	for i := 0; i < len(fields)-1; i += 2 {
		if key, ok := fields[i].(string); ok {
			entry.Fields[key] = fields[i+1]
		}
	}

	// Convert to JSON and write
	jsonData, err := json.Marshal(entry)
	if err != nil {
		l.Logger.Printf("Failed to marshal log entry: %v", err)
		return
	}

	l.output.Write(jsonData)
	l.output.Write([]byte("\n"))
}

// HTTPLogger middleware for Echo
func (l *Logger) HTTPLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process request
			err := next(c)

			// Calculate latency
			latency := time.Since(start)

			// Get request info
			req := c.Request()
			res := c.Response()

			// Create log entry
			entry := LogEntry{
				Timestamp: start.UTC().Format(time.RFC3339),
				Level:     "INFO",
				Message:   "HTTP request",
				RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
				UserID:    c.Get("user_id").(string),
				IP:        c.RealIP(),
				UserAgent: req.UserAgent(),
				Method:    req.Method,
				URI:       req.RequestURI,
				Status:    res.Status,
				Latency:   latency.String(),
			}

			if err != nil {
				entry.Error = err.Error()
				entry.Level = "ERROR"
			}

			// Convert to JSON and write
			jsonData, _ := json.Marshal(entry)
			l.output.Write(jsonData)
			l.output.Write([]byte("\n"))

			return err
		}
	}
}
