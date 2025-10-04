package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_ValidConfig(t *testing.T) {
	// Set required environment variables
	os.Setenv("JWT_SECRET", "this-is-a-very-long-jwt-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Verify JWT secret is set
	if cfg.JWT.Secret == "" {
		t.Error("JWT.Secret is empty")
	}

	if len(cfg.JWT.Secret) < 32 {
		t.Errorf("JWT.Secret length is %d, expected >= 32", len(cfg.JWT.Secret))
	}
}

func TestLoad_MissingJWTSecret(t *testing.T) {
	// Ensure JWT_SECRET is not set
	os.Unsetenv("JWT_SECRET")

	cfg, err := Load()
	if err == nil {
		t.Fatal("Load() should return error when JWT_SECRET is missing")
	}

	if cfg != nil {
		t.Error("Load() should return nil config when JWT_SECRET is missing")
	}

	expectedErr := "JWT_SECRET environment variable is required"
	if err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestLoad_ShortJWTSecret(t *testing.T) {
	// Set JWT_SECRET that is too short
	os.Setenv("JWT_SECRET", "short")
	defer os.Unsetenv("JWT_SECRET")

	cfg, err := Load()
	if err == nil {
		t.Fatal("Load() should return error when JWT_SECRET is too short")
	}

	if cfg != nil {
		t.Error("Load() should return nil config when JWT_SECRET is too short")
	}

	// Check error message contains expected text
	if err != nil && err.Error()[:37] != "JWT_SECRET must be at least 32 charac" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	// Set only required JWT_SECRET
	os.Setenv("JWT_SECRET", "this-is-a-very-long-jwt-secret-key-for-testing")
	defer os.Unsetenv("JWT_SECRET")

	// Unset optional variables to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("HOST")
	os.Unsetenv("DB_PATH")
	os.Unsetenv("LOG_LEVEL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Test default values
	if cfg.Server.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Server.Port)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Expected default host 0.0.0.0, got %s", cfg.Server.Host)
	}

	if cfg.Database.Path != "./data/app.db" {
		t.Errorf("Expected default DB path ./data/app.db, got %s", cfg.Database.Path)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("Expected default log level info, got %s", cfg.Logging.Level)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	// Set custom environment variables
	os.Setenv("JWT_SECRET", "this-is-a-very-long-jwt-secret-key-for-testing")
	os.Setenv("PORT", "9000")
	os.Setenv("HOST", "localhost")
	os.Setenv("DB_PATH", "/custom/path/db.sqlite")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("ENV", "development")

	defer func() {
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("PORT")
		os.Unsetenv("HOST")
		os.Unsetenv("DB_PATH")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("ENV")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Test custom values
	if cfg.Server.Port != "9000" {
		t.Errorf("Expected port 9000, got %s", cfg.Server.Port)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("Expected host localhost, got %s", cfg.Server.Host)
	}

	if cfg.Database.Path != "/custom/path/db.sqlite" {
		t.Errorf("Expected DB path /custom/path/db.sqlite, got %s", cfg.Database.Path)
	}

	if cfg.Logging.Level != "debug" {
		t.Errorf("Expected log level debug, got %s", cfg.Logging.Level)
	}

	if !cfg.Logging.EnableDebug {
		t.Error("Expected EnableDebug to be true in development mode")
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Returns default when env not set",
			key:          "TEST_VAR_NOT_SET",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
		{
			name:         "Returns env value when set",
			key:          "TEST_VAR_SET",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetEnvInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue int
		envValue     string
		expected     int
	}{
		{
			name:         "Returns default when env not set",
			key:          "TEST_INT_NOT_SET",
			defaultValue: 100,
			envValue:     "",
			expected:     100,
		},
		{
			name:         "Returns parsed int when valid",
			key:          "TEST_INT_VALID",
			defaultValue: 100,
			envValue:     "200",
			expected:     200,
		},
		{
			name:         "Returns default when invalid int",
			key:          "TEST_INT_INVALID",
			defaultValue: 100,
			envValue:     "invalid",
			expected:     100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnvInt(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetEnvDuration(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue time.Duration
		envValue     string
		expected     time.Duration
	}{
		{
			name:         "Returns default when env not set",
			key:          "TEST_DURATION_NOT_SET",
			defaultValue: 5 * time.Minute,
			envValue:     "",
			expected:     5 * time.Minute,
		},
		{
			name:         "Returns parsed duration when valid",
			key:          "TEST_DURATION_VALID",
			defaultValue: 5 * time.Minute,
			envValue:     "10m",
			expected:     10 * time.Minute,
		},
		{
			name:         "Returns default when invalid duration",
			key:          "TEST_DURATION_INVALID",
			defaultValue: 5 * time.Minute,
			envValue:     "invalid",
			expected:     5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			} else {
				os.Unsetenv(tt.key)
			}

			result := getEnvDuration(tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestLoggingConfig(t *testing.T) {
	tests := []struct {
		name                string
		envValue            string
		expectedEnableDebug bool
	}{
		{
			name:                "Development mode enables debug",
			envValue:            "development",
			expectedEnableDebug: true,
		},
		{
			name:                "Production mode disables debug",
			envValue:            "production",
			expectedEnableDebug: false,
		},
		{
			name:                "Empty env defaults to production",
			envValue:            "",
			expectedEnableDebug: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("JWT_SECRET", "this-is-a-very-long-jwt-secret-key-for-testing")
			defer os.Unsetenv("JWT_SECRET")

			if tt.envValue != "" {
				os.Setenv("ENV", tt.envValue)
				defer os.Unsetenv("ENV")
			} else {
				os.Unsetenv("ENV")
			}

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() returned error: %v", err)
			}

			if cfg.Logging.EnableDebug != tt.expectedEnableDebug {
				t.Errorf("Expected EnableDebug=%v, got %v", tt.expectedEnableDebug, cfg.Logging.EnableDebug)
			}
		})
	}
}
