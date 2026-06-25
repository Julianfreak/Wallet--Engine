package application

import (
	"context"
	"fmt"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/Julianfreak/Wallet--Engine/internal/ports"
	"github.com/google/uuid"
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
	notif ports.NotificationSender, // <-- Quinto argumento
) *TransferService {
	return &TransferService{
		accountRepo:     ar,
		transactionRepo: tr,
		txManager:       tm,
		logger:          log,
		notifier:        notif, // <-- Asignamos el notificador
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

		msg := fmt.Sprintf("Transferencia exitosa de %s a %s por un monto de $%.2f. (ID Auditoría: %s)", fromID, toID, amount, txID)
		s.logger.Info(msg)

		emailMsg := fmt.Sprintf("¡Hola! Recibiste una transferencia de $%.2f de la cuenta %s.", amount, fromID)
		_ = s.notifier.Send("usuario_mercado_libre@email.com", emailMsg)

		go func() {
			_ = s.notifier.Send("usuario_mercado_libre@email.com", emailMsg)
		}()

		return nil // Si llegamos aquí sin errores, el TxManager hace COMMIT automáticamente
	})
}
