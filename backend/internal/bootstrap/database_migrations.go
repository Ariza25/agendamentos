package bootstrap

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// applyDatabaseSQL executes the unified canonical SQL schema file (backend/sql/schema.sql).
// Since the product is pre-homologation / greenfield, schema migrations are unified
// into a single declarative, idempotent schema file.
func applyDatabaseSQL(db *sql.DB) error {
	schemaPath := resolveSchemaFile()
	if schemaPath == "" {
		log.Println("no schema.sql found, skipping database schema init")
		return nil
	}

	content, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema file %s: %w", schemaPath, err)
	}

	log.Printf("applying canonical database schema from %s", schemaPath)
	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("execute canonical schema: %w", err)
	}

	log.Println("database schema initialized successfully")
	return nil
}

func resolveSchemaFile() string {
	candidates := []string{
		"sql/schema.sql",
		"backend/sql/schema.sql",
		"../sql/schema.sql",
		"../../sql/schema.sql",
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}
