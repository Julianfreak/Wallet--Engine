package application

import (
	"context"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/Julianfreak/Wallet--Engine/internal/ports"
)

type DashboardData struct {
	Accounts     []domain.Account     `json:"accounts"`
	Transactions []domain.Transaction `json:"transactions"`
}

type DashboardService struct {
	accountRepo     ports.AccountRepository
	transactionRepo ports.TransactionRepository
}

func NewDashboardService(accountRepo ports.AccountRepository, transactionRepo ports.TransactionRepository) *DashboardService {
	return &DashboardService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *DashboardService) GetDashboardData(ctx context.Context) (DashboardData, error) {
	// 1. Obtener transacciones usando el método GetAll que acabamos de crear
	transactions, err := s.transactionRepo.GetAll(ctx)
	if err != nil {
		return DashboardData{}, err
	}

	// 2. Obtener cuentas (si tienes un método para listar cuentas, úsalo aquí.
	// De lo contrario, puedes consultar las cuentas por defecto A1 y A2 o agregarlas a la respuesta).
	// Ejemplo básico consultando las cuentas principales:
	acc1, _ := s.accountRepo.GetByID(ctx, "A1")
	acc2, _ := s.accountRepo.GetByID(ctx, "A2")

	var accounts []domain.Account
	if acc1 != nil {
		accounts = append(accounts, *acc1)
	}
	if acc2 != nil {
		accounts = append(accounts, *acc2)
	}

	return DashboardData{
		Accounts:     accounts,
		Transactions: transactions,
	}, nil
}
