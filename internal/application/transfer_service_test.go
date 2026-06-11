package application

import (
	"context"
	"errors"
	"testing"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
)

// --- DOBLES DE PRUEBA (FAKES) ---

type FakeTxManager struct{}

func (f *FakeTxManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx) // El simulador ejecuta la función inmediatamente sin abrir SQL real
}

type FakeAccountRepository struct {
	accounts map[string]*domain.Account
}

func (r *FakeAccountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	acc, exists := r.accounts[id]
	if !exists {
		return nil, errors.New("cuenta no encontrada")
	}
	// Devolvemos una copia para evitar mutaciones directas indeseadas en el mapa
	return &domain.Account{ID: acc.ID, Owner: acc.Owner, Balance: acc.Balance}, nil
}

func (r *FakeAccountRepository) Save(ctx context.Context, acc *domain.Account) error {
	r.accounts[acc.ID] = acc
	return nil
}

type FakeTransactionRepository struct {
	savedTransactions []*domain.Transaction
}

func (r *FakeTransactionRepository) Save(ctx context.Context, tx *domain.Transaction) error {
	r.savedTransactions = append(r.savedTransactions, tx)
	return nil
}

// --- PRUEBAS DEL CASO DE USO ---

func TestTransferService_Execute_Success(t *testing.T) {
	// Configuración del entorno simulado
	accountRepo := &FakeAccountRepository{
		accounts: map[string]*domain.Account{
			"A1": {ID: "A1", Owner: "Julian", Balance: 1000.0},
			"A2": {ID: "A2", Owner: "Mercado Libre", Balance: 0.0},
		},
	}
	transactionRepo := &FakeTransactionRepository{}
	txManager := &FakeTxManager{}

	service := NewTransferService(accountRepo, transactionRepo, txManager)
	ctx := context.Background()

	// Ejecución
	err := service.Execute(ctx, "A1", "A2", 400.0)

	// Verificaciones (Assertions)
	if err != nil {
		t.Fatalf("se esperaba una transferencia exitosa, pero falló: %v", err)
	}

	// Verificar saldos finales en el repositorio simulado
	fromAcc, _ := accountRepo.GetByID(ctx, "A1")
	toAcc, _ := accountRepo.GetByID(ctx, "A2")

	if fromAcc.Balance != 600.0 {
		t.Errorf("origen esperado: 600.0, obtenido: %.2f", fromAcc.Balance)
	}
	if toAcc.Balance != 400.0 {
		t.Errorf("destino esperado: 400.0, obtenido: %.2f", toAcc.Balance)
	}

	// Verificar que se haya registrado la auditoría de la transacción
	if len(transactionRepo.savedTransactions) != 1 {
		t.Errorf("se esperaba 1 registro de transacción guardado, se obtuvieron: %d", len(transactionRepo.savedTransactions))
	}
}

func TestTransferService_Execute_InsufficientFunds(t *testing.T) {
	accountRepo := &FakeAccountRepository{
		accounts: map[string]*domain.Account{
			"A1": {ID: "A1", Owner: "Julian", Balance: 50.0},
			"A2": {ID: "A2", Owner: "Mercado Libre", Balance: 0.0},
		},
	}
	transactionRepo := &FakeTransactionRepository{}
	txManager := &FakeTxManager{}

	service := NewTransferService(accountRepo, transactionRepo, txManager)
	ctx := context.Background()

	// Intentar transferir más de lo que se tiene
	err := service.Execute(ctx, "A1", "A2", 200.0)

	if err == nil {
		t.Error("se esperaba un fallo por fondos insuficientes, pero el servicio retornó éxito")
	}

	// Verificar que los saldos se mantuvieron intactos (no hubo cambios accidentales)
	fromAcc, _ := accountRepo.GetByID(ctx, "A1")
	if fromAcc.Balance != 50.0 {
		t.Errorf("el saldo de origen cambió de forma insegura a: %.2f", fromAcc.Balance)
	}
}
