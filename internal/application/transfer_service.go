package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/Julianfreak/Wallet--Engine/internal/infrastructure/metrics"
	"github.com/Julianfreak/Wallet--Engine/internal/ports"
	"github.com/google/uuid"
)

var (
	ErrInsufficientFunds = errors.New("fondos insuficientes para realizar la transferencia")
	ErrSameAccount       = errors.New("no se puede transferir a la misma cuenta")
	ErrAccountNotFound   = errors.New("cuenta no encontrada")
)

type TransferService struct {
	accountRepo     ports.AccountRepository
	transactionRepo ports.TransactionRepository
	txManager       ports.TxManager
	logger          ports.Logger
	notifier        ports.NotificationSender
}

func NewTransferService(
	ar ports.AccountRepository,
	tr ports.TransactionRepository,
	tm ports.TxManager,
	log ports.Logger,
	notif ports.NotificationSender,
) *TransferService {
	return &TransferService{
		accountRepo:     ar,
		transactionRepo: tr,
		txManager:       tm,
		logger:          log,
		notifier:        notif,
	}
}

// En internal/application/transfer_service.go
func (s *TransferService) GetTransactions(ctx context.Context, ownerEmail string) ([]domain.Transaction, error) {
	// 1. Buscar las cuentas que pertenecen a este usuario
	accounts, err := s.accountRepo.FindByOwner(ctx, ownerEmail)
	if err != nil {
		return nil, err
	}

	var allTransactions []domain.Transaction
	seenTxIDs := make(map[string]bool)

	// 2. Extraer las transacciones de cada cuenta del usuario (evitando duplicados)
	for _, acc := range accounts {
		txs, err := s.transactionRepo.GetByAccountID(ctx, acc.ID)
		if err != nil {
			continue
		}
		for _, tx := range txs {
			if !seenTxIDs[tx.ID] {
				seenTxIDs[tx.ID] = true
				allTransactions = append(allTransactions, tx)
			}
		}
	}

	if allTransactions == nil {
		return []domain.Transaction{}, nil
	}

	return allTransactions, nil
}

type TransferCommand struct {
	FromAccountID string
	ToAccountID   string
	Amount        float64
}

// Añade ownerEmail como parámetro para validar la titularidad
func (s *TransferService) Execute(ctx context.Context, ownerEmail string, cmd TransferCommand) error {
	startTime := time.Now()
	defer func() {
		metrics.TransferDuration.Observe(time.Since(startTime).Seconds())
	}()

	if cmd.FromAccountID == cmd.ToAccountID {
		return ErrSameAccount
	}

	if cmd.Amount <= 0 {
		return errors.New("el monto de la transferencia debe ser mayor a cero")
	}

	err := s.txManager.WithTransaction(ctx, func(ctxTx context.Context) error {
		// A. Obtener cuenta de origen
		fromAccount, err := s.accountRepo.GetByID(ctxTx, cmd.FromAccountID)
		if err != nil {
			return fmt.Errorf("error al obtener cuenta de origen: %w", err)
		}
		if fromAccount == nil {
			return fmt.Errorf("%w: %s", ErrAccountNotFound, cmd.FromAccountID)
		}

		// VALIDACIÓN DE SEGURIDAD: Verificar que la cuenta de origen pertenezca al usuario autenticado
		if fromAccount.Owner != ownerEmail {
			return errors.New("no autorizado: la cuenta de origen no pertenece al usuario autenticado")
		}

		// B. Validar fondos suficientes
		if fromAccount.Balance < cmd.Amount {
			return ErrInsufficientFunds
		}

		// C. Obtener cuenta de destino
		toAccount, err := s.accountRepo.GetByID(ctxTx, cmd.ToAccountID)
		if err != nil {
			return fmt.Errorf("error al obtener cuenta de destino: %w", err)
		}
		if toAccount == nil {
			return fmt.Errorf("%w: %s", ErrAccountNotFound, cmd.ToAccountID)
		}

		// D. Modificar saldos en memoria
		fromAccount.Balance -= cmd.Amount
		toAccount.Balance += cmd.Amount

		// E. Guardar cambios en el repositorio de cuentas
		if err := s.accountRepo.Save(ctxTx, fromAccount); err != nil {
			return fmt.Errorf("error al actualizar cuenta origen: %w", err)
		}
		if err := s.accountRepo.Save(ctxTx, toAccount); err != nil {
			return fmt.Errorf("error al actualizar cuenta destino: %w", err)
		}

		// F. Registrar la transacción histórica con ID único y fecha
		txRecord := &domain.Transaction{
			ID:            uuid.New().String(),
			FromAccountID: cmd.FromAccountID,
			ToAccountID:   cmd.ToAccountID,
			Amount:        cmd.Amount,
			CreatedAt:     time.Now(),
		}
		if err := s.transactionRepo.Save(ctxTx, txRecord); err != nil {
			return fmt.Errorf("error al guardar el registro de transacción: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("fallo en la transferencia: %w", err)
	}

	s.logger.Info(fmt.Sprintf("Transferencia exitosa de %.2f desde %s hacia %s", cmd.Amount, cmd.FromAccountID, cmd.ToAccountID))

	go func() {
		_ = s.notifier.Send(cmd.ToAccountID, "Has recibido una transferencia exitosa")
	}()

	return nil
}
