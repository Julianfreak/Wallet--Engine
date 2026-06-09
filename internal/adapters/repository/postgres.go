package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	_ "github.com/lib/pq"
)

// Clave privada para evitar colisiones dentro del contexto
type txKey struct{}

// PostgresTxManager implementa ports.TxManager
type PostgresTxManager struct {
	db *sql.DB
}

func NewPostgresTxManager(db *sql.DB) *PostgresTxManager {
	return &PostgresTxManager{db: db}
}

// WithTransaction ejecuta una función de forma atómica (BEGIN, COMMIT/ROLLBACK)
func (m *PostgresTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Manejo de pánicos para evitar bloqueos en la base de datos
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// Inyectamos la transacción de Postgres dentro del contexto
	txCtx := context.WithValue(ctx, txKey{}, tx)

	// Ejecutamos la lógica de negocio
	if err := fn(txCtx); err != nil {
		_ = tx.Rollback() // Si algo falla dentro del flujo, revertimos TODO
		return err
	}

	return tx.Commit() // Si todo sale bien, guardamos permanentemente en el disco
}

// Interfaz interna para unificar consultas sobre sql.DB o sql.Tx
type dbExecutor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type PostgresAccountRepository struct {
	db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

// getExecutor extrae la transacción del contexto si existe, de lo contrario usa la DB global
func (r *PostgresAccountRepository) getExecutor(ctx context.Context) dbExecutor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return r.db
}

func (r *PostgresAccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	query := `SELECT id, owner, balance FROM accounts WHERE id = $1`
	var acc domain.Account

	err := r.getExecutor(ctx).QueryRowContext(ctx, query, id).Scan(&acc.ID, &acc.Owner, &acc.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("cuenta no encontrada en la base de datos")
		}
		return nil, err
	}
	return &acc, nil
}

func (r *PostgresAccountRepository) Save(ctx context.Context, acc *domain.Account) error {
	query := `
		INSERT INTO accounts (id, owner, balance) 
		VALUES ($1, $2, $3)
		ON CONFLICT (id) 
		DO UPDATE SET balance = EXCLUDED.balance`

	_, err := r.getExecutor(ctx).ExecContext(ctx, query, acc.ID, acc.Owner, acc.Balance)
	return err
}

type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

func (r *PostgresTransactionRepository) getExecutor(ctx context.Context) dbExecutor {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return r.db
}

func (r *PostgresTransactionRepository) Save(ctx context.Context, tx *domain.Transaction) error {
	query := `INSERT INTO transactions (id, from_account_id, to_account_id, amount) VALUES ($1, $2, $3, $4)`
	_, err := r.getExecutor(ctx).ExecContext(ctx, query, tx.ID, tx.FromAccountID, tx.ToAccountID, tx.Amount)
	return err
}
