package services

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"api-key-generator/internal/config"
	"api-key-generator/internal/database"
	"api-key-generator/internal/models"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dbPath := ":memory:"
	db, err := database.Initialize(dbPath)
	require.NoError(t, err)
	// SQLite single-writer optimization
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func newTestService(t *testing.T) *APIKeyService {
	t.Helper()
	db := setupTestDB(t)
	cfg := &config.Config{
		APIKey: config.APIKeyConfig{
			Length:         16,
			ExpiryDuration: 24 * time.Hour,
			Prefix:         "ak_",
		},
	}
	svc, err := NewAPIKeyService(db, cfg)
	require.NoError(t, err, "Failed to create APIKeyService")
	t.Cleanup(func() { _ = svc.Close() })
	return svc
}

func TestHasAnyPermission(t *testing.T) {
	granted := []string{"recipes:read", "users:write"}
	require.True(t, hasAnyPermission(granted, []string{"recipes:read"}))
	require.True(t, hasAnyPermission(granted, []string{"RECIPES:READ"}))
	require.True(t, hasAnyPermission(granted, []string{"foo", "users:write"}))
	require.False(t, hasAnyPermission(granted, []string{"foo", "bar"}))
	require.True(t, hasAnyPermission(granted, nil))
	require.True(t, hasAnyPermission(granted, []string{}))
}

func TestHasAllPermissions(t *testing.T) {
	granted := []string{"recipes:read", "users:write"}
	require.True(t, hasAllPermissions(granted, []string{"recipes:read"}))
	require.True(t, hasAllPermissions(granted, []string{"RECIPES:READ"}))
	require.True(t, hasAllPermissions(granted, []string{"recipes:read", "users:write"}))
	require.False(t, hasAllPermissions(granted, []string{"recipes:read", "users:read"}))
	require.True(t, hasAllPermissions(granted, nil))
	require.True(t, hasAllPermissions(granted, []string{}))
}

func TestAuthorizeAny(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	// Create key with known permissions
	req := &models.CreateAPIKeyRequest{
		Name:        "test-any",
		Permissions: []string{"recipes:read", "users:write"},
	}
	key, err := svc.CreateAPIKey(ctx, req)
	require.NoError(t, err)

	// success when one matches
	_, ok, err := svc.AuthorizeAny(ctx, key.Key, []string{"recipes:read", "recipes:write"})
	require.NoError(t, err)
	require.True(t, ok)

	// forbidden when none match
	_, ok, err = svc.AuthorizeAny(ctx, key.Key, []string{"recipes:write"})
	require.NoError(t, err)
	require.False(t, ok)

	// invalid key
	_, ok, err = svc.AuthorizeAny(ctx, "invalid", []string{"recipes:read"})
	require.Error(t, err)
	require.False(t, ok)
}

func TestAuthorizeAll(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	req := &models.CreateAPIKeyRequest{
		Name:        "test-all",
		Permissions: []string{"recipes:read", "users:write"},
	}
	key, err := svc.CreateAPIKey(ctx, req)
	require.NoError(t, err)

	// success when all match
	_, ok, err := svc.AuthorizeAll(ctx, key.Key, []string{"recipes:read", "users:write"})
	require.NoError(t, err)
	require.True(t, ok)

	// forbidden when one missing
	_, ok, err = svc.AuthorizeAll(ctx, key.Key, []string{"recipes:read", "users:read"})
	require.NoError(t, err)
	require.False(t, ok)
}
