package ports

import (
	"context"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
)

// TxManager define el contrato para manejar la atomicidad de la base de datos
type TxManager interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// AccountRepository es el puerto de salida para el almacenamiento de cuentas.
// Cualquier base de datos (Postgres, MySQL, Mock) deberá implementar estas funciones.
type AccountRepository interface {
	GetByID(ctx context.Context, id string) (*domain.Account, error)
	Save(ctx context.Context, acc *domain.Account) error
}

// TransactionRepository es el puerto de salida para el historial de transacciones.
type TransactionRepository interface {
	Save(ctx context.Context, tx *domain.Transaction) error
}
