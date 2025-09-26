package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api-key-generator/internal/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testPathCreateProtocol = "/api/v1/supplements/protocols"
	testHeaderAPIKey       = "X-API-Key"
	scopeWrite             = "supplements:write"
	scopeRead              = "supplements:read"
)

// MockAPIKeyAuthorizer implements APIKeyAuthorizer for testing
type MockAPIKeyAuthorizer struct{ mock.Mock }

func (m *MockAPIKeyAuthorizer) AuthorizeAny(ctx context.Context, rawKey string, scopes []string) (*models.APIKey, bool, error) {
	args := m.Called(ctx, rawKey, scopes)
	if ak := args.Get(0); ak != nil {
		return ak.(*models.APIKey), args.Bool(1), args.Error(2)
	}
	return (*models.APIKey)(nil), args.Bool(1), args.Error(2)
}

func (m *MockAPIKeyAuthorizer) RecordUsage(ctx context.Context, apiKeyID, endpoint, method string, status int, ip, userAgent string) error {
	args := m.Called(ctx, apiKeyID, endpoint, method, status, ip, userAgent)
	return args.Error(0)
}

// MockSuppService implements SupplementServiceContract for testing
type MockSuppService struct{ mock.Mock }

func (m *MockSuppService) GenerateSupplementProtocol(ctx context.Context, req *models.SupplementRequest) (*models.SupplementProtocol, error) {
	args := m.Called(ctx, req)
	if sp := args.Get(0); sp != nil {
		return sp.(*models.SupplementProtocol), args.Error(1)
	}
	return (*models.SupplementProtocol)(nil), args.Error(1)
}

func (m *MockSuppService) GetSupplementProtocol(ctx context.Context, id string) (*models.SupplementProtocol, error) {
	args := m.Called(ctx, id)
	if sp := args.Get(0); sp != nil {
		return sp.(*models.SupplementProtocol), args.Error(1)
	}
	return (*models.SupplementProtocol)(nil), args.Error(1)
}

// noopValidator implements echo.Validator for tests
type noopValidator struct{}

func (n *noopValidator) Validate(i interface{}) error { return nil }

func TestGenerateSupplementProtocolUnauthorizedWithoutKey(t *testing.T) {
	e := echo.New()
	e.Validator = &noopValidator{}
	mockAuth := new(MockAPIKeyAuthorizer)
	h := NewSupplementHandler(nil, mockAuth)

	req := httptest.NewRequest(http.MethodPost, testPathCreateProtocol, strings.NewReader(`{}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockAuth.On("AuthorizeAny", mock.Anything, "", []string{scopeWrite}).Return((*models.APIKey)(nil), false, assert.AnError)

	err := h.GenerateSupplementProtocol(c)
	assert.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusUnauthorized, httpErr.Code)
}

func TestGenerateSupplementProtocolSuccess(t *testing.T) {
	e := echo.New()
	e.Validator = &noopValidator{}
	mockAuth := new(MockAPIKeyAuthorizer)
	mockSupp := new(MockSuppService)
	h := NewSupplementHandler(mockSupp, mockAuth)

	reqBody := `{"user_id":"u1","workout_goal":"build_muscle","workout_type":"gym","duration_minutes":60,"intensity":"high"}`
	req := httptest.NewRequest(http.MethodPost, testPathCreateProtocol, strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(testHeaderAPIKey, "ak_valid")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	apiKey := &models.APIKey{ID: "k1", Permissions: []string{"supplements:write"}}
	mockAuth.On("AuthorizeAny", mock.Anything, "ak_valid", []string{scopeWrite}).Return(apiKey, true, nil)
	mockSupp.On("GenerateSupplementProtocol", mock.Anything, mock.AnythingOfType("*models.SupplementRequest")).Return(&models.SupplementProtocol{ID: "p1", UserID: "u1"}, nil)
	mockAuth.On("RecordUsage", mock.Anything, "k1", testPathCreateProtocol, http.MethodPost, http.StatusOK, mock.Anything, mock.Anything).Return(nil)

	err := h.GenerateSupplementProtocol(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp models.SuccessResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.True(t, resp.Success)
}

func TestGetSupplementProtocolSuccess(t *testing.T) {
	e := echo.New()
	e.Validator = &noopValidator{}
	mockAuth := new(MockAPIKeyAuthorizer)
	mockSupp := new(MockSuppService)
	h := NewSupplementHandler(mockSupp, mockAuth)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/supplements/protocols/p1", nil)
	req.Header.Set(testHeaderAPIKey, "ak_valid")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("p1")

	apiKey := &models.APIKey{ID: "k1", Permissions: []string{"supplements:read"}}
	mockAuth.On("AuthorizeAny", mock.Anything, "ak_valid", []string{scopeRead}).Return(apiKey, true, nil)
	mockSupp.On("GetSupplementProtocol", mock.Anything, "p1").Return(&models.SupplementProtocol{ID: "p1"}, nil)
	mockAuth.On("RecordUsage", mock.Anything, "k1", "/api/v1/supplements/protocols/p1", http.MethodGet, http.StatusOK, mock.Anything, mock.Anything).Return(nil).Maybe()

	err := h.GetSupplementProtocol(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetSupplementSafetyInfoSuccess(t *testing.T) {
	e := echo.New()
	e.Validator = &noopValidator{}
	mockAuth := new(MockAPIKeyAuthorizer)
	h := NewSupplementHandler(nil, mockAuth)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/supplements/safety/caffeine", nil)
	req.Header.Set(testHeaderAPIKey, "ak_valid")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("name")
	c.SetParamValues("caffeine")

	apiKey := &models.APIKey{ID: "k1", Permissions: []string{"supplements:read"}}
	mockAuth.On("AuthorizeAny", mock.Anything, "ak_valid", []string{scopeRead}).Return(apiKey, true, nil)
	mockAuth.On("RecordUsage", mock.Anything, "k1", "/api/v1/supplements/safety/caffeine", http.MethodGet, http.StatusOK, mock.Anything, mock.Anything).Return(nil).Maybe()

	err := h.GetSupplementSafetyInfo(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
