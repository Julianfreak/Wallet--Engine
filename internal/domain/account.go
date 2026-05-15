package domain

import (
	"errors"
	"time"
)

// Errores de negocio específicos que entenderá cualquier parte de nuestra app
var ErrInsufficientFunds = errors.New("fondos insuficientes en la cuenta")
var ErrInvalidAmount = errors.New("el monto debe ser mayor a cero")

type Account struct {
	ID        string
	Owner     string
	Balance   float64
	CreatedAt time.Time
}

// Deposit añade saldo a la cuenta validando que el monto sea lógico
func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.Balance += amount
	return nil
}

// Withdraw retira saldo validando que existan fondos suficientes
func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if amount > a.Balance {
		return ErrInsufficientFunds
	}
	a.Balance -= amount
	return nil
}
