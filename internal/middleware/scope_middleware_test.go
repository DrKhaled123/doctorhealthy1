package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"api-key-generator/internal/config"
	"api-key-generator/internal/database"
	"api-key-generator/internal/models"
	"api-key-generator/internal/services"
)

func newTestServices(t *testing.T) (*services.APIKeyService, echo.Context) {
	t.Helper()
	dbPath := t.TempDir() + "/test.db"
	db, err := database.Initialize(dbPath)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	svc, err := services.NewAPIKeyService(db, &config.Config{APIKey: config.APIKeyConfig{Length: 16, ExpiryDuration: 24 * time.Hour, Prefix: "ak_"}})
	require.NoError(t, err, "Failed to create APIKeyService")
	t.Cleanup(func() { _ = svc.Close() })

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return svc, c
}

func TestScopesAny_AllowsWhenAnyMatches(t *testing.T) {
	svc, c := newTestServices(t)
	key, err := svc.CreateAPIKey(c.Request().Context(), &models.CreateAPIKeyRequest{Name: "kk", Permissions: []string{"recipes:read"}})
	require.NoError(t, err)

	mw := ScopesAny(svc, "recipes:read", "recipes:write")
	handler := mw(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	c.Request().Header.Set("X-API-Key", key.Key)
	err = handler(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, c.Response().Status)
}

func TestScopesAny_RejectsWhenNoneMatch(t *testing.T) {
	svc, c := newTestServices(t)
	key, err := svc.CreateAPIKey(c.Request().Context(), &models.CreateAPIKeyRequest{Name: "kk", Permissions: []string{"recipes:read"}})
	require.NoError(t, err)

	mw := ScopesAny(svc, "recipes:write")
	handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	c.Request().Header.Set("X-API-Key", key.Key)
	err = handler(c)
	require.Error(t, err)
}

func TestScopesAll_AllRequired(t *testing.T) {
	svc, c := newTestServices(t)
	key, err := svc.CreateAPIKey(c.Request().Context(), &models.CreateAPIKeyRequest{Name: "kk", Permissions: []string{"recipes:read", "users:write"}})
	require.NoError(t, err)

	mw := ScopesAll(svc, "recipes:read", "users:write")
	handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	c.Request().Header.Set("X-API-Key", key.Key)
	err = handler(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, c.Response().Status)
}

func TestScopesAll_ForbiddenWhenMissing(t *testing.T) {
	svc, c := newTestServices(t)
	key, err := svc.CreateAPIKey(c.Request().Context(), &models.CreateAPIKeyRequest{Name: "kk", Permissions: []string{"recipes:read"}})
	require.NoError(t, err)

	mw := ScopesAll(svc, "recipes:read", "users:write")
	handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	c.Request().Header.Set("X-API-Key", key.Key)
	err = handler(c)
	require.Error(t, err)
}

func TestScopesAny_MissingAPIKeyHeader(t *testing.T) {
	svc, c := newTestServices(t)

	mw := ScopesAny(svc, "recipes:read")
	handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	// Don't set X-API-Key header
	err := handler(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
	require.Equal(t, "API key required", httpErr.Message)
}

func TestScopesAll_MissingAPIKeyHeader(t *testing.T) {
	svc, c := newTestServices(t)

	mw := ScopesAll(svc, "recipes:read")
	handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	// Don't set X-API-Key header
	err := handler(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
	require.Equal(t, "API key required", httpErr.Message)
}

func TestScopesAny_InvalidAPIKey(t *testing.T) {
	svc, c := newTestServices(t)

	mw := ScopesAny(svc, "recipes:read")
	handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	c.Request().Header.Set("X-API-Key", "invalid-key")
	err := handler(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
	require.Equal(t, "invalid API key", httpErr.Message)
}

func TestScopesAll_InvalidAPIKey(t *testing.T) {
	svc, c := newTestServices(t)

	mw := ScopesAll(svc, "recipes:read")
	handler := mw(func(c echo.Context) error { return c.String(http.StatusOK, "ok") })

	c.Request().Header.Set("X-API-Key", "invalid-key")
	err := handler(c)
	require.Error(t, err)
	httpErr, ok := err.(*echo.HTTPError)
	require.True(t, ok)
	require.Equal(t, http.StatusUnauthorized, httpErr.Code)
	require.Equal(t, "invalid API key", httpErr.Message)
}

func TestScopesAny_StoresAPIKeyRecord(t *testing.T) {
	svc, c := newTestServices(t)
	key, err := svc.CreateAPIKey(c.Request().Context(), &models.CreateAPIKeyRequest{Name: "test", Permissions: []string{"recipes:read"}})
	require.NoError(t, err)

	mw := ScopesAny(svc, "recipes:read")
	handler := mw(func(c echo.Context) error {
		record := c.Get("api_key_record")
		require.NotNil(t, record)
		apiKey, ok := record.(*models.APIKey)
		require.True(t, ok)
		require.Equal(t, key.ID, apiKey.ID)
		return c.String(http.StatusOK, "ok")
	})

	c.Request().Header.Set("X-API-Key", key.Key)
	err = handler(c)
	require.NoError(t, err)
}

func TestScopesAll_StoresAPIKeyRecord(t *testing.T) {
	svc, c := newTestServices(t)
	key, err := svc.CreateAPIKey(c.Request().Context(), &models.CreateAPIKeyRequest{Name: "test", Permissions: []string{"recipes:read", "users:write"}})
	require.NoError(t, err)

	mw := ScopesAll(svc, "recipes:read", "users:write")
	handler := mw(func(c echo.Context) error {
		record := c.Get("api_key_record")
		require.NotNil(t, record)
		apiKey, ok := record.(*models.APIKey)
		require.True(t, ok)
		require.Equal(t, key.ID, apiKey.ID)
		return c.String(http.StatusOK, "ok")
	})

	c.Request().Header.Set("X-API-Key", key.Key)
	err = handler(c)
	require.NoError(t, err)
}
