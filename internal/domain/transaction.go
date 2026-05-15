package domain

import "time"

type Transaction struct {
	ID            string
	FromAccountID string
	ToAccountID   string
	Amount        float64
	CreatedAt     time.Time
}

// NewTransaction es una función constructora que asegura que una transacción
// nazca siempre con los datos mínimos obligatorios.
func NewTransaction(id, from, to string, amount float64) *Transaction {
	return &Transaction{
		ID:            id,
		FromAccountID: from,
		ToAccountID:   to,
		Amount:        amount,
		CreatedAt:     time.Now(),
	}
}
