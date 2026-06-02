package repository

import (
	"database/sql"
	"errors"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	_ "github.com/lib/pq" // El guion bajo indica que se importa solo por su efecto secundario (registrar el driver)
)

// PostgresAccountRepository implementa ports.AccountRepository usando SQL real
type PostgresAccountRepository struct {
	db *sql.DB
}

func NewPostgresAccountRepository(db *sql.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

// GetByID busca la cuenta en la tabla de PostgreSQL
func (r *PostgresAccountRepository) GetByID(id string) (*domain.Account, error) {
	query := `SELECT id, owner, balance FROM accounts WHERE id = $1`

	var acc domain.Account
	// Ejecutamos la consulta y escaneamos los valores directamente en la estructura
	err := r.db.QueryRow(query, id).Scan(&acc.ID, &acc.Owner, &acc.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("cuenta no encontrada en la base de datos")
		}
		return nil, err
	}

	return &acc, nil
}

// Save inserta una nueva cuenta o actualiza su saldo si ya existe (Upsert)
func (r *PostgresAccountRepository) Save(acc *domain.Account) error {
	query := `
		INSERT INTO accounts (id, owner, balance) 
		VALUES ($1, $2, $3)
		ON CONFLICT (id) 
		DO UPDATE SET balance = EXCLUDED.balance`

	_, err := r.db.Exec(query, acc.ID, acc.Owner, acc.Balance)
	return err
}

// PostgresTransactionRepository implementa ports.TransactionRepository
type PostgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{db: db}
}

// Save inserta el registro histórico de la transferencia
func (r *PostgresTransactionRepository) Save(tx *domain.Transaction) error {
	query := `INSERT INTO transactions (id, from_account_id, to_account_id, amount) VALUES ($1, $2, $3, $4)`

	_, err := r.db.Exec(query, tx.ID, tx.FromAccountID, tx.ToAccountID, tx.Amount)
	return err
}
