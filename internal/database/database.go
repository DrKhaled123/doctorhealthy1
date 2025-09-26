package database

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Initialize initializes the database connection and creates tables
func Initialize(dbPath string) (*sql.DB, error) {
	// Validate and sanitize database path
	cleanPath := filepath.Clean(dbPath)
	if strings.Contains(cleanPath, "..") {
		return nil, fmt.Errorf("invalid database path: path traversal detected")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(cleanPath), 0600); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Open database connection
	db, err := sql.Open("sqlite3", cleanPath+"?_journal_mode=WAL&_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings (optimized for SQLite)
	db.SetMaxOpenConns(1) // SQLite single-writer optimization
	db.SetMaxIdleConns(5)

	// Create tables
	if err := createTables(db); err != nil {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Failed to close database connection: %v", closeErr)
		}
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Apply non-fatal migrations (deterministic, idempotent)
	if err := RunMigrations(db); err != nil {
		log.Printf("Migrations encountered errors: %v", err)
	}

	return db, nil
}

// createTables creates all necessary tables
func createTables(db *sql.DB) error {
	queries := []string{
		createAPIKeysTable,
		createAPIKeyUsageTable,
		createAPIKeysIndex,
		createUsageIndex,
		createUsersTable,
		createNutritionPlansTable,
		createWorkoutPlansTable,
		createHealthPlansTable,
		createRecipesTable,
		createHealthIndexes,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}

const createAPIKeysTable = `
CREATE TABLE IF NOT EXISTS api_keys (
    id TEXT PRIMARY KEY,
    key TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    user_id TEXT,
    permissions TEXT NOT NULL, -- JSON array
    is_active BOOLEAN NOT NULL DEFAULT 1,
    expires_at DATETIME,
    last_used_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    usage_count INTEGER NOT NULL DEFAULT 0,
    rate_limit INTEGER,
    rate_limit_used INTEGER NOT NULL DEFAULT 0
);`

const createAPIKeyUsageTable = `
CREATE TABLE IF NOT EXISTS api_key_usage (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    api_key_id TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    method TEXT NOT NULL,
    status INTEGER NOT NULL,
    timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ip_address TEXT,
    user_agent TEXT,
    FOREIGN KEY (api_key_id) REFERENCES api_keys(id) ON DELETE CASCADE
);`

// #nosec G101 - SQL index creation, not credentials
const createAPIKeysIndex = `
CREATE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys(key);
CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_is_active ON api_keys(is_active);
CREATE INDEX IF NOT EXISTS idx_api_keys_expires_at ON api_keys(expires_at);`

const createUsageIndex = `
CREATE INDEX IF NOT EXISTS idx_usage_api_key_id ON api_key_usage(api_key_id);
CREATE INDEX IF NOT EXISTS idx_usage_timestamp ON api_key_usage(timestamp);
CREATE INDEX IF NOT EXISTS idx_usage_endpoint ON api_key_usage(endpoint);`

// Health management tables
const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    age INTEGER NOT NULL,
    weight REAL NOT NULL,
    height REAL NOT NULL,
    gender TEXT NOT NULL,
    activity_level TEXT NOT NULL,
    metabolic_rate TEXT NOT NULL,
    goal TEXT NOT NULL,
    food_dislikes TEXT, -- JSON array
    allergies TEXT, -- JSON array
    diseases TEXT, -- JSON array
    medications TEXT, -- JSON array
    preferred_cuisine TEXT,
    language TEXT DEFAULT 'en',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

const createNutritionPlansTable = `
CREATE TABLE IF NOT EXISTS nutrition_plans (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    calories_per_day INTEGER NOT NULL,
    protein_grams REAL NOT NULL,
    carbs_grams REAL NOT NULL,
    fats_grams REAL NOT NULL,
    plan_type TEXT NOT NULL,
    calculation_method TEXT,
    disclaimer TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createWorkoutPlansTable = `
CREATE TABLE IF NOT EXISTS workout_plans (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    goal TEXT NOT NULL,
    workout_type TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createHealthPlansTable = `
CREATE TABLE IF NOT EXISTS health_plans (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    disclaimer TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`

const createRecipesTable = `
CREATE TABLE IF NOT EXISTS recipes (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cuisine TEXT NOT NULL,
    prep_time INTEGER,
    cook_time INTEGER,
    servings INTEGER,
    calories INTEGER,
    difficulty TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);`

const createHealthIndexes = `
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_nutrition_plans_user_id ON nutrition_plans(user_id);
CREATE INDEX IF NOT EXISTS idx_workout_plans_user_id ON workout_plans(user_id);
CREATE INDEX IF NOT EXISTS idx_health_plans_user_id ON health_plans(user_id);
CREATE INDEX IF NOT EXISTS idx_recipes_cuisine ON recipes(cuisine);
CREATE INDEX IF NOT EXISTS idx_recipes_difficulty ON recipes(difficulty);
CREATE INDEX IF NOT EXISTS idx_recipes_created_at ON recipes(created_at);`

// executeWorkoutSchema executes the workout database schema
func executeWorkoutSchema(db *sql.DB) error {
	// Read the workout schema file
	schemaPath := filepath.Join("internal", "database", "workout_schema.sql")

	// Check if file exists
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		log.Printf("Workout schema file not found at %s, skipping workout tables creation", schemaPath)
		return nil
	}

	// Read the schema file
	file, err := os.Open(schemaPath) // nosec G304 - hardcoded path to internal schema file
	if err != nil {
		return fmt.Errorf("failed to open workout schema file: %w", err)
	}
	defer func() { _ = file.Close() }()

	schemaBytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read workout schema file: %w", err)
	}

	schema := string(schemaBytes)

	// Split schema into individual statements and enforce order: TABLEs -> remaining
	raw := strings.Split(schema, ";")
	var tableStmts []string
	var otherStmts []string

	for _, s := range raw {
		stmt := strings.TrimSpace(s)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}
		upper := strings.ToUpper(stmt)
		if strings.HasPrefix(upper, "CREATE TABLE") {
			tableStmts = append(tableStmts, stmt)
		} else {
			otherStmts = append(otherStmts, stmt)
		}
	}

	// Execute tables in first transaction
	tx1, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	if _, err := tx1.Exec("PRAGMA foreign_keys=ON"); err != nil {
		_ = tx1.Rollback()
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	for _, stmt := range tableStmts {
		// Log table being created (best-effort parse of table name)
		tbl := ""
		parts := strings.Fields(stmt)
		if len(parts) >= 3 {
			// CREATE TABLE [IF] [NOT] EXISTS <name>
			for i := 0; i < len(parts); i++ {
				if strings.ToUpper(parts[i]) == "EXISTS" && i+1 < len(parts) {
					tbl = strings.Trim(parts[i+1], "`\"[]")
					break
				}
			}
		}
		if tbl != "" {
			log.Printf("Creating table: %s", tbl)
		}
		if _, err := tx1.Exec(stmt); err != nil {
			_ = tx1.Rollback()
			return fmt.Errorf("failed to execute workout TABLE statement: %v; stmt: %.120s", err, stmt)
		}
	}
	if err := tx1.Commit(); err != nil {
		return fmt.Errorf("failed to commit workout tables: %w", err)
	}

	// Check critical tables exist before proceeding with indexes/views
	exercisesExists := false
	if err := db.QueryRow("SELECT 1 FROM sqlite_master WHERE type='table' AND name='exercises'").Scan(new(int)); err == nil {
		exercisesExists = true
	}
	if !exercisesExists {
		log.Printf("Exercises table not found after creation phase")
	}

	// Execute remaining statements in separate transaction
	tx2, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for indexes/views: %w", err)
	}
	if _, err := tx2.Exec("PRAGMA foreign_keys=ON"); err != nil {
		_ = tx2.Rollback()
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	for _, stmt := range otherStmts {
		// Skip indexes/views referencing exercises if table not present (safety)
		if !exercisesExists && (strings.Contains(strings.ToLower(stmt), " on exercises(") || strings.Contains(strings.ToLower(stmt), "from exercises")) {
			log.Printf("Skipping statement due to missing exercises table: %.120s", stmt)
			continue
		}
		if _, err := tx2.Exec(stmt); err != nil {
			_ = tx2.Rollback()
			return fmt.Errorf("failed to execute workout schema statement: %v; stmt: %.120s", err, stmt)
		}
	}
	if err := tx2.Commit(); err != nil {
		return fmt.Errorf("failed to commit workout indexes/views: %w", err)
	}

	return nil
}
