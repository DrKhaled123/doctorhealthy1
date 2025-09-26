package database

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RunMigrations applies idempotent SQL migrations. It never fatally exits; errors are logged and returned to caller.
func RunMigrations(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("nil db")
	}

	if err := ensureMigrationsTable(db); err != nil {
		return fmt.Errorf("ensure migrations table: %w", err)
	}

	// Single workout schema migration derived from workout_schema.sql
	applied, err := isMigrationApplied(db, "001_workout_schema")
	if err != nil {
		return fmt.Errorf("check migration: %w", err)
	}
	if !applied {
		if err := applyWorkoutSchema(db); err != nil {
			// Log and proceed (non-fatal startup policy)
			log.Printf("Workout schema migration failed: %v", err)
		} else {
			if err := recordApplied(db, "001_workout_schema"); err != nil {
				log.Printf("Failed to record migration 001_workout_schema: %v", err)
			}
			// Seed minimal data for indexes/views to apply correctly and for dev testing
			if err := seedWorkoutMinimal(db); err != nil {
				log.Printf("Workout minimal seed failed: %v", err)
			}
			// Re-apply indexes/views now that seed exists
			if err := applyWorkoutIndexesAndViews(db); err != nil {
				log.Printf("Re-applying workout indexes/views failed: %v", err)
			}
		}
	}

	return nil
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT UNIQUE NOT NULL,
        applied_at DATETIME NOT NULL
    );`)
	return err
}

func isMigrationApplied(db *sql.DB, name string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(1) FROM schema_migrations WHERE name = ?", name).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func recordApplied(db *sql.DB, name string) error {
	_, err := db.Exec("INSERT OR IGNORE INTO schema_migrations (name, applied_at) VALUES (?, ?)", name, time.Now().UTC())
	return err
}

func applyWorkoutSchema(db *sql.DB) error {
	schemaPath := filepath.Join("internal", "database", "workout_schema.sql") // nosec G304 - hardcoded path to internal schema file
	f, err := os.Open(schemaPath)
	if err != nil {
		return fmt.Errorf("open schema: %w", err)
	}
	defer func() { _ = f.Close() }()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	stmts := splitSQLStatements(string(bytes))

	// Phase 1: CREATE TABLE statements
	if err := execFiltered(db, stmts, func(s string) bool {
		up := strings.ToUpper(strings.TrimSpace(s))
		return strings.HasPrefix(up, "CREATE TABLE")
	}); err != nil {
		return err
	}

	// Phase 2: all remaining statements, skipping indexes/views referencing missing tables
	if err := execFiltered(db, stmts, func(s string) bool {
		up := strings.ToUpper(strings.TrimSpace(s))
		return !strings.HasPrefix(up, "CREATE TABLE")
	}); err != nil {
		return err
	}
	return nil
}

func splitSQLStatements(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		s := strings.TrimSpace(p)
		if s == "" || strings.HasPrefix(s, "--") {
			continue
		}
		out = append(out, s)
	}
	return out
}

func execFiltered(db *sql.DB, stmts []string, predicate func(string) bool) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	if _, err := tx.Exec("PRAGMA foreign_keys=ON"); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("pragma fk: %w", err)
	}
	for _, stmt := range stmts {
		if !predicate(stmt) {
			continue
		}
		s := strings.TrimSpace(stmt)
		if s == "" {
			continue
		}
		// Skip index/view creation if base table missing (safety)
		lower := strings.ToLower(s)
		if strings.HasPrefix(lower, "create index") || strings.HasPrefix(lower, "create view") {
			if tbl := referencedTable(lower); tbl != "" {
				if !tableExists(tx, tbl) {
					log.Printf("Skipping statement due to missing table %s: %.120s", tbl, s)
					continue
				}
			}
		}
		if _, err := tx.Exec(s); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("exec stmt failed: %v; stmt: %.120s", err, s)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func tableExists(tx *sql.Tx, name string) bool {
	var x int
	err := tx.QueryRow("SELECT 1 FROM sqlite_master WHERE type='table' AND name=?", name).Scan(&x)
	return err == nil
}

func referencedTable(stmtLower string) string {
	// naive extraction: look for " on <table>(" or " from <table>"
	for _, kw := range []string{" on ", " from ", " view ", " table "} {
		idx := strings.Index(stmtLower, kw)
		if idx >= 0 {
			rest := stmtLower[idx+len(kw):]
			// token until space or paren
			for i := 0; i < len(rest); i++ {
				if rest[i] == ' ' || rest[i] == '(' || rest[i] == '\n' || rest[i] == '\r' || rest[i] == '\t' {
					return strings.Trim(rest[:i], "`\"[]")
				}
			}
			return strings.Trim(rest, "`\"[] ")
		}
	}
	return ""
}

func seedWorkoutMinimal(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Create minimal rows only if tables exist
	tables := []string{"exercises", "workout_programs", "weekly_plans", "daily_workouts"}
	for _, t := range tables {
		if !tableExistsTx(db, t) {
			// Skip if base table missing; migration already logged
			return nil
		}
	}

	// Seed exercise
	_, _ = tx.Exec(`INSERT OR IGNORE INTO exercises (id, name_en, category, primary_muscles, difficulty, instructions_en) VALUES ('ex1','Bodyweight Squat','strength','["quads","glutes"]','beginner','Squat down and stand up')`)
	// Seed program
	_, _ = tx.Exec(`INSERT OR IGNORE INTO workout_programs (id, name_en, level, split_type, description_en, duration_weeks, sessions_per_week, goals) VALUES ('prog1','Starter Program','beginner','full_body','Beginner full body',4,3,'["weight_loss"]')`)
	// Seed week
	_, _ = tx.Exec(`INSERT OR IGNORE INTO weekly_plans (id, program_id, week_number, intensity_level, volume_increase_percent) VALUES ('w1','prog1',1,'low',0)`)
	// Seed day
	_, _ = tx.Exec(`INSERT OR IGNORE INTO daily_workouts (id, weekly_plan_id, day_number, focus_en, duration_minutes, is_rest_day) VALUES ('d1','w1',1,'Full Body A',45,0)`)

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func tableExistsTx(db *sql.DB, name string) bool {
	var x int
	err := db.QueryRow("SELECT 1 FROM sqlite_master WHERE type='table' AND name=?", name).Scan(&x)
	return err == nil
}

// applyWorkoutIndexesAndViews re-applies only indexes and views from the schema file.
func applyWorkoutIndexesAndViews(db *sql.DB) error {
	schemaPath := filepath.Join("internal", "database", "workout_schema.sql") // nosec G304 - hardcoded path to internal schema file
	f, err := os.Open(schemaPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	bytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	stmts := splitSQLStatements(string(bytes))
	// Only non-table statements
	return execFiltered(db, stmts, func(s string) bool {
		up := strings.ToUpper(strings.TrimSpace(s))
		return strings.HasPrefix(up, "CREATE INDEX") || strings.HasPrefix(up, "CREATE VIEW")
	})
}
