package testutils

import (
	"context"
	"errors"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
)

// ==========================================
// 1. FAKE TX MANAGER
// ==========================================
type FakeTxManager struct{}

func (f *FakeTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx) // Ejecuta la función de inmediato sin SQL real
}

// ==========================================
// 2. FAKE ACCOUNT REPOSITORY
// ==========================================
type FakeAccountRepository struct {
	Accounts map[string]*domain.Account // "A" Mayúscula para ser público
}

func (r *FakeAccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	acc, exists := r.Accounts[id]
	if !exists {
		return nil, errors.New("cuenta no encontrada")
	}
	// Devolvemos copia para evitar mutaciones directas en el mapa
	return &domain.Account{ID: acc.ID, Owner: acc.Owner, Balance: acc.Balance}, nil
}

func (r *FakeAccountRepository) Save(ctx context.Context, acc *domain.Account) error {
	if r.Accounts == nil {
		r.Accounts = make(map[string]*domain.Account)
	}
	r.Accounts[acc.ID] = acc
	return nil
}

// ==========================================
// 3. FAKE TRANSACTION REPOSITORY
// ==========================================
type FakeTransactionRepository struct {
	SavedTransactions []*domain.Transaction // "S" Mayúscula para ser público
}

func (r *FakeTransactionRepository) Save(ctx context.Context, tx *domain.Transaction) error {
	r.SavedTransactions = append(r.SavedTransactions, tx)
	return nil
}

// ==========================================
// 4. FAKE LOGGER
// ==========================================
type FakeLogger struct {
	LastMessage string
}

func (f *FakeLogger) Info(message string) {
	f.LastMessage = message
}

// ==========================================
// 5. FAKE NOTIFICATION SENDER
// ==========================================
type FakeNotificationSender struct {
	Called bool
	Done   chan bool
}

func (f *FakeNotificationSender) Send(recipient string, message string) error {
	f.Called = true
	if f.Done != nil {
		f.Done <- true // Evita bloquear si el canal no está inicializado
	}
	return nil
}
