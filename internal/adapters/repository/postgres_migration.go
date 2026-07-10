package repository

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations utiliza golang-migrate para aplicar los cambios de esquema
// dbURL debe tener el formato: "postgres://usuario:password@localhost:5432/nombre_db?sslmode=disable"
func RunMigrations(dbURL string) error {
	// "file://migrations" apunta a tu carpeta en la raíz del proyecto
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return fmt.Errorf("error creando instancia de migración: %w", err)
	}

	// Ejecuta todas las migraciones pendientes
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("fallo al aplicar migraciones: %w", err)
	}

	return nil
}
