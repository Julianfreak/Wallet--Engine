package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/Julianfreak/Wallet--Engine/internal/infrastructure/metrics"
	"github.com/Julianfreak/Wallet--Engine/internal/ports"
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

type TransferCommand struct {
	FromAccountID string
	ToAccountID   string
	Amount        float64 // En centavos
}

func (s *TransferService) Execute(ctx context.Context, cmd TransferCommand) error {
	startTime := time.Now()
	// Al final de la funcion, calculamos la duracion y la registramos en el histograma
	defer func() {
		metrics.TransferDuration.Observe(time.Since(startTime).Seconds())
	}()
	if cmd.FromAccountID == cmd.ToAccountID {
		return ErrSameAccount
	}

	// 2. Regla de negocio: Monto válido
	if cmd.Amount <= 0 {
		return errors.New("el monto de la transferencia debe ser mayor a cero")
	}

	// 3. Ejecutar dentro de una transacción atómica usando el TxManager
	err := s.txManager.WithTransaction(ctx, func(ctxTx context.Context) error {
		// A. Obtener cuenta de origen
		fromAccount, err := s.accountRepo.GetByID(ctxTx, cmd.FromAccountID)
		if err != nil {
			return fmt.Errorf("error al obtener cuenta de origen: %w", err)
		}
		if fromAccount == nil {
			return fmt.Errorf("%w: %s", ErrAccountNotFound, cmd.FromAccountID)
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
		// (Nota: Asegúrate de tener un método Update o usar Save según tu diseño de repositorio)
		if err := s.accountRepo.Save(ctxTx, fromAccount); err != nil {
			return fmt.Errorf("error al actualizar cuenta origen: %w", err)
		}
		if err := s.accountRepo.Save(ctxTx, toAccount); err != nil {
			return fmt.Errorf("error al actualizar cuenta destino: %w", err)
		}

		// F. Registrar la transacción histórica
		txRecord := &domain.Transaction{
			FromAccountID: cmd.FromAccountID,
			ToAccountID:   cmd.ToAccountID,
			Amount:        cmd.Amount,
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

	// 2. Enviar notificación asíncrona
	go func() {
		_ = s.notifier.Send(cmd.ToAccountID, "Has recibido una transferencia exitosa")
	}()

	return nil
}
