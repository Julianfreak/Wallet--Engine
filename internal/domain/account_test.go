package domain

import (
	"testing"
)

func TestAccount_Withdraw(t *testing.T) {
	// Escenario 1: Retiro exitoso con fondos suficientes
	acc := &Account{ID: "A1", Owner: "Julian", Balance: 500.0}
	err := acc.Withdraw(200.0)

	if err != nil {
		t.Fatalf("se esperaba un retiro exitoso, pero dio error: %v", err)
	}
	if acc.Balance != 300.0 {
		t.Errorf("saldo esperado: 300.0, saldo obtenido: %.2f", acc.Balance)
	}

	// Escenario 2: Intento de retiro con fondos insuficientes
	err = acc.Withdraw(1000.0)
	if err == nil {
		t.Error("se esperaba un error por fondos insuficientes, pero la operación pasó")
	}
	if err.Error() != "fondos insuficientes en la cuenta" {
		t.Errorf("mensaje de error inesperado: %v", err)
	}
}

func TestAccount_Deposit(t *testing.T) {
	acc := &Account{ID: "A2", Owner: "Mercado Libre", Balance: 100.0}
	acc.Deposit(150.0)

	if acc.Balance != 250.0 {
		t.Errorf("saldo esperado tras depósito: 250.0, obtenido: %.2f", acc.Balance)
	}
}
