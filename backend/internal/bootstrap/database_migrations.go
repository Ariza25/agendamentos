package bootstrap

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func applyDatabaseSQL(db *sql.DB) error {
	for _, directory := range []string{"sql/catalog", "sql/migrations"} {
		if err := applySQLDirectory(db, directory); err != nil {
			return err
		}
	}
	return nil
}

func applySQLDirectory(db *sql.DB, directory string) error {
	entries, err := os.ReadDir(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read sql directory %s: %w", directory, err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".sql" {
			continue
		}
		if directory == "sql/catalog" && strings.Contains(entry.Name(), "seed") {
			continue
		}
		files = append(files, filepath.Join(directory, entry.Name()))
	}
	sort.Strings(files)

	if directory == "sql/migrations" {
		if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS catalog.schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`); err != nil {
			return fmt.Errorf("ensure schema migrations table: %w", err)
		}
	}

	for _, file := range files {
		if directory == "sql/migrations" {
			if err := applyMigrationFile(db, file); err != nil {
				return err
			}
			continue
		}
		if err := applySQLFile(db, file); err != nil {
			return err
		}
	}
	return nil
}

func applySQLFile(db *sql.DB, file string) error {
	sqlBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read sql file %s: %w", file, err)
	}
	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("apply sql file %s: %w", file, err)
	}
	return nil
}

func applyMigrationFile(db *sql.DB, file string) (err error) {
	version := filepath.Base(file)
	var alreadyApplied bool
	if err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM catalog.schema_migrations WHERE version = $1)", version).Scan(&alreadyApplied); err != nil {
		return fmt.Errorf("check migration %s: %w", file, err)
	}
	if alreadyApplied {
		return nil
	}

	sqlBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("read sql file %s: %w", file, err)
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", file, err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec(string(sqlBytes)); err != nil {
		return fmt.Errorf("apply sql file %s: %w", file, err)
	}
	if _, err = tx.Exec("INSERT INTO catalog.schema_migrations (version) VALUES ($1)", version); err != nil {
		return fmt.Errorf("record migration %s: %w", file, err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", file, err)
	}
	return nil
}
