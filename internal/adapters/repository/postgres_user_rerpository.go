package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, email, password_hash) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.PasswordHash)
	return err
}

// GetAll recupera todas las transacciones almacenadas en la base de datos
func (r *PostgresTransactionRepository) GetAll(ctx context.Context) ([]domain.Transaction, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, created_at FROM transactions ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var transactions []domain.Transaction
	for rows.Next() {
		var tx domain.Transaction
		if err := rows.Scan(&tx.ID, &tx.FromAccountID, &tx.ToAccountID, &tx.Amount, &tx.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	if transactions == nil {
		return []domain.Transaction{}, nil
	}

	return transactions, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}
