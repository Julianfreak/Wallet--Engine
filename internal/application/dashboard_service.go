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

func (s *DashboardService) GetDashboardData(ctx context.Context, userEmail string) (DashboardData, error) {

	accounts, err := s.accountRepo.FindByOwner(ctx, userEmail)
	if err != nil {
		return DashboardData{}, err
	}

	// 3. Obtener solo las transacciones correspondientes a las cuentas de este usuario
	var userTransactions []domain.Transaction
	for _, acc := range accounts {
		txs, err := s.transactionRepo.GetByAccountID(ctx, acc.ID)
		if err == nil {
			userTransactions = append(userTransactions, txs...)
		}
	}

	return DashboardData{
		Accounts:     accounts,
		Transactions: userTransactions,
	}, nil
}
