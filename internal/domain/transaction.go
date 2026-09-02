package domain

import "time"

type Transaction struct {
	ID            string    `json:"id"`
	FromAccountID string    `json:"from_account_id"`
	ToAccountID   string    `json:"to_account_id"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
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
