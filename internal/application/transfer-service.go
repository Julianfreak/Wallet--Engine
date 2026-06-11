package application

import (
	"context"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/Julianfreak/Wallet--Engine/internal/ports"
	"github.com/google/uuid"
)

type TransferService struct {
	accountRepo     ports.AccountRepository
	transactionRepo ports.TransactionRepository
	txManager       ports.TxManager
}

func NewTransferService(ar ports.AccountRepository, tr ports.TransactionRepository, tm ports.TxManager) *TransferService {
	return &TransferService{
		accountRepo:     ar,
		transactionRepo: tr,
		txManager:       tm,
	}
}

func (s *TransferService) Execute(ctx context.Context, fromID, toID string, amount float64) error {
	// Envolvemos todo el proceso en una transacción atómica externa
	return s.txManager.WithTransaction(ctx, func(txCtx context.Context) error {

		// 1. Buscar cuentas usando el contexto transaccional (txCtx)
		fromAccount, err := s.accountRepo.GetByID(txCtx, fromID)
		if err != nil {
			return err
		}

		toAccount, err := s.accountRepo.GetByID(txCtx, toID)
		if err != nil {
			return err
		}

		// 2. Lógica matemática del Dominio
		if err := fromAccount.Withdraw(amount); err != nil {
			return err
		}

		if err := toAccount.Deposit(amount); err != nil {
			return err
		}

		// 3. Generar identificador único de auditoría
		txID := uuid.NewString()
		tx := domain.NewTransaction(txID, fromID, toID, amount)

		// 4. Persistir los cambios en lote
		if err := s.accountRepo.Save(txCtx, fromAccount); err != nil {
			return err
		}
		if err := s.accountRepo.Save(txCtx, toAccount); err != nil {
			return err
		}
		if err := s.transactionRepo.Save(txCtx, tx); err != nil {
			return err
		}

		return nil // Si llegamos aquí sin errores, el TxManager hace COMMIT automáticamente
	})
}
