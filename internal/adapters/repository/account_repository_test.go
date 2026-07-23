package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	_ "github.com/lib/pq" // Driver de Postgres
	"github.com/stretchr/testify/assert"
)

// setupTestDB prepara la conexión a Docker y limpia la base de datos
func setupTestDB(t *testing.T) *sql.DB {
	// Usamos los datos extraídos de tu docker-compose.yml
	dsn := "postgres://wallet_user:wallet_password@localhost:5432/wallet_db?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	assert.NoError(t, err, "No se pudo conectar a PostgreSQL en Docker")

	err = db.Ping()
	assert.NoError(t, err, "El contenedor de DB está arriba pero la base de datos no responde")

	// Limpiamos la tabla para evitar choques con pruebas anteriores
	_, err = db.Exec("TRUNCATE TABLE accounts CASCADE")
	assert.NoError(t, err, "No se pudo limpiar la tabla 'accounts'")

	return db
}

func TestPostgresAccountRepository_SaveAndGet(t *testing.T) {
	// 1. Arrange: Conectamos a la BD y creamos tu repositorio real
	db := setupTestDB(t)
	defer db.Close() // Se asegura de soltar la conexión al terminar

	repo := NewPostgresAccountRepository(db)
	ctx := context.Background()

	// Creamos una cuenta de prueba (ajusta los campos si tu struct es distinto)
	cuentaEsperada := &domain.Account{
		ID:      "UUID-TEST-001",
		Owner:   "Julian",
		Balance: 1500,
	}

	// 2. Act: Guardamos la cuenta y la volvemos a buscar
	errSave := repo.Save(ctx, cuentaEsperada)
	assert.NoError(t, errSave, "Falló al hacer el INSERT en la base de datos")

	cuentaRecuperada, errGet := repo.GetByID(ctx, cuentaEsperada.ID)

	// 3. Assert: Verificamos que los datos guardados coincidan exactamente
	assert.NoError(t, errGet, "Falló al hacer el SELECT en la base de datos")
	assert.NotNil(t, cuentaRecuperada)

	assert.Equal(t, cuentaEsperada.ID, cuentaRecuperada.ID)
	assert.Equal(t, cuentaEsperada.Owner, cuentaRecuperada.Owner)
	assert.Equal(t, cuentaEsperada.Balance, cuentaRecuperada.Balance)
}
