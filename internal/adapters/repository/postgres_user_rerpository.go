package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	fmt.Printf("🔍 DB QUERY Buscando ID exacto: '%s'\n", id)
	query := `SELECT id, email, password_hash, created_at 
              FROM users 
              WHERE id = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	// 🔍 Añade esta línea para ver el valor exacto y su longitud
	fmt.Printf("🔍 DB QUERY - Buscando email exacto: '%s' (longitud: %d)\n", email, len(email))

	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("❌ No se encontró ningún registro con el email: '%s'\n", email)
			return nil, errors.New("usuario no encontrado precaución")
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET email = $1, password_hash = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.PasswordHash, user.ID)
	return err
}
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
func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
