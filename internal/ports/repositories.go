package ports

import "github.com/Julianfreak/wallet-engine/internal/domain"

// AccountRepository es el puerto de salida para el almacenamiento de cuentas.
// Cualquier base de datos (Postgres, MySQL, Mock) deberá implementar estas funciones.
type AccountRepository interface {
	GetByID(id string) (*domain.Account, error)
	Save(account *domain.Account) error
}

// TransactionRepository es el puerto de salida para el historial de transacciones.
type TransactionRepository interface {
	Save(transaction *domain.Transaction) error
}
