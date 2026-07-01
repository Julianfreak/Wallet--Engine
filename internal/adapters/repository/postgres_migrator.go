package repository

import (
	"database/sql"
	_ "embed" // <-- Requerido para usar go:embed
	"fmt"
)

// Indicamos a Go que guarde el contenido del SQL en la variable schemaSQL en tiempo de compilación
//
//go:embed migrations/schema.sql
var schemaSQL string

// RunMigrations ejecuta el script incrustado para asegurar la estructura de la base de datos
func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("fallo en migraciones automáticas: %w", err)
	}
	return nil
}
