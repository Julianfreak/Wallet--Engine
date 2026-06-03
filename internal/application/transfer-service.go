package application

import (
	"errors"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/Julianfreak/Wallet--Engine/internal/ports"
	"github.com/google/uuid"
)

// Errores específicos de la aplicación
var ErrAccountNotFound = errors.New("cuenta no encontrada")

// TransferService es nuestra estructura de servicio.
// Fíjate que "inyectamos" las interfaces (los puertos), NO bases de datos reales.
type TransferService struct {
	accountRepo     ports.AccountRepository
	transactionRepo ports.TransactionRepository
}

// NewTransferService es el constructor. Quien quiera hacer una transferencia,
// debe pasarnos "algo" que cumpla con los contratos de los repositorios.
func NewTransferService(ar ports.AccountRepository, tr ports.TransactionRepository) *TransferService {
	return &TransferService{
		accountRepo:     ar,
		transactionRepo: tr,
	}
}

// Execute orquesta todo el proceso de transferir dinero de una cuenta a otra.
func (s *TransferService) Execute(fromID, toID string, amount float64) error {
	// 1. Buscar las cuentas usando el puerto (no sabemos si vienen de Postgres o Memoria)
	fromAccount, err := s.accountRepo.GetByID(fromID)
	if err != nil {
		return err
	}

	toAccount, err := s.accountRepo.GetByID(toID)
	if err != nil {
		return err
	}
	// 2. Usar el Dominio puro para la lógica matemática y validaciones
	if err := fromAccount.Withdraw(amount); err != nil {
		return err // Si no hay fondos, el dominio lanza el error y cortamos aquí
	}

	if err := toAccount.Deposit(amount); err != nil {
		return err
	}
	// 3. Crear el registro histórico de la transacción
	txID := uuid.NewString()
	tx := domain.NewTransaction(txID, fromID, toID, amount)
	// 4. Guardar los nuevos estados usando los puertos
	if err := s.accountRepo.Save(fromAccount); err != nil {
		return err
	}
	if err := s.accountRepo.Save(toAccount); err != nil {
		return err
	}
	if err := s.transactionRepo.Save(tx); err != nil {
		return err
	}

	return nil
}
