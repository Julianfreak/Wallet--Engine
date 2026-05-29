package repository

import (
	"errors"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
)

// AccountMemoryRepo implementa ports.AccountRepository usando la memoria RAM
type AccountMemoryRepo struct {
	data map[string]*domain.Account
}

func NewAccountMemoryRepo() *AccountMemoryRepo {
	return &AccountMemoryRepo{
		data: make(map[string]*domain.Account),
	}
}

func (r *AccountMemoryRepo) GetByID(id string) (*domain.Account, error) {
	acc, exists := r.data[id]
	if !exists {
		return nil, errors.New("cuenta no encontrada en la base de datos")
	}
	return acc, nil
}

func (r *AccountMemoryRepo) Save(acc *domain.Account) error {
	r.data[acc.ID] = acc
	return nil
}

// TransactionMemoryRepo implementa ports.TransactionRepository
type TransactionMemoryRepo struct {
	data map[string]*domain.Transaction
}

func NewTransactionMemoryRepo() *TransactionMemoryRepo {
	return &TransactionMemoryRepo{
		data: make(map[string]*domain.Transaction),
	}
}

func (r *TransactionMemoryRepo) Save(tx *domain.Transaction) error {
	r.data[tx.ID] = tx
	return nil
}
