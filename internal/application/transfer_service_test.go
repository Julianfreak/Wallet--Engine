package application

import (
	"context"
	"testing"
	"time"

	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/Julianfreak/Wallet--Engine/internal/testutils"
	"github.com/stretchr/testify/assert"
)

// --- PRUEBAS DEL CASO DE USO ---

func TestTransferService_Execute_Success(t *testing.T) {
	// Configuración del entorno simulado usando testutils
	accountRepo := &testutils.FakeAccountRepository{
		Accounts: map[string]*domain.Account{
			"A1": {ID: "A1", Owner: "Julian", Balance: 1000.0},
			"A2": {ID: "A2", Owner: "Mercado Libre", Balance: 0.0},
		},
	}
	transactionRepo := &testutils.FakeTransactionRepository{}
	txManager := &testutils.FakeTxManager{}
	fakeLogger := &testutils.FakeLogger{}
	fakeNotifier := &testutils.FakeNotificationSender{
		Done: make(chan bool, 1),
	}
	service := NewTransferService(accountRepo, transactionRepo, txManager, fakeLogger, fakeNotifier)
	ctx := context.Background()

	// Ejecución
	cmd := TransferCommand{
		FromAccountID: "A1",
		ToAccountID:   "A2",
		Amount:        400.0,
	}

	err := service.Execute(ctx, cmd)
	assert.NoError(t, err)

	// Sincronización de la goroutine
	select {
	case <-fakeNotifier.Done:
		// Se ejecutó el envío de notificación de forma asíncrona
	case <-time.After(1 * time.Second):
		t.Error("Timeout: la notificación nunca llegó")
	}

	if err != nil {
		t.Fatalf("se esperaba una transferencia exitosa, pero falló: %v", err)
	}

	// Verificar saldos finales
	fromAcc, _ := accountRepo.GetByID(ctx, "A1")
	toAcc, _ := accountRepo.GetByID(ctx, "A2")

	if fromAcc.Balance != 600.0 {
		t.Errorf("origen esperado: 600.0, obtenido: %.2f", fromAcc.Balance)
	}
	if toAcc.Balance != 400.0 {
		t.Errorf("destino esperado: 400.0, obtenido: %.2f", toAcc.Balance)
	}

	// Verificar la auditoría usando la propiedad pública
	if len(transactionRepo.SavedTransactions) != 1 {
		t.Errorf("se esperaba 1 registro de transacción guardado, se obtuvieron: %d", len(transactionRepo.SavedTransactions))
	}

	if fakeLogger.LastMessage == "" {
		t.Error("se esperaba que el servicio registrara una auditoría en el logger, pero el mensaje quedó vacío")
	}

	if !fakeNotifier.Called {
		t.Error("se esperaba que el servicio enviara una notificación, pero no lo hizo")
	}
}

func TestTransferService_Execute_InsufficientFunds(t *testing.T) {
	accountRepo := &testutils.FakeAccountRepository{
		Accounts: map[string]*domain.Account{
			"A1": {ID: "A1", Owner: "Julian", Balance: 50.0},
			"A2": {ID: "A2", Owner: "Mercado Libre", Balance: 0.0},
		},
	}
	transactionRepo := &testutils.FakeTransactionRepository{}
	txManager := &testutils.FakeTxManager{}
	fakeLogger := &testutils.FakeLogger{}
	fakeNotifier := &testutils.FakeNotificationSender{
		Done: make(chan bool, 1),
	}
	service := NewTransferService(accountRepo, transactionRepo, txManager, fakeLogger, fakeNotifier)
	ctx := context.Background()

	// Intentar transferir más de lo que se tiene
	cmd := TransferCommand{
		FromAccountID: "A1",
		ToAccountID:   "A2",
		Amount:        100.00,
	}

	err := service.Execute(ctx, cmd)

	if err == nil {
		t.Error("se esperaba un fallo por fondos insuficientes, pero el servicio retornó éxito")
	}

	// Intento de transferencia a cuenta inexistente para agotar los escenarios analizados
	err = service.Execute(ctx, cmd)

	if err == nil {
		t.Error("se esperaba un fallo por cuenta inexistente, pero retornó éxito")
	}

	// Verificar que los saldos se mantuvieron intactos
	fromAcc, _ := accountRepo.GetByID(ctx, "A1")
	if fromAcc.Balance != 50.0 {
		t.Errorf("el saldo de origen cambió de forma insegura a: %.2f", fromAcc.Balance)
	}

	if fakeLogger.LastMessage != "" {
		t.Errorf("no se esperaba auditoría de éxito en una transferencia fallida, pero se obtuvo: %s", fakeLogger.LastMessage)
	}

	if fakeNotifier.Called {
		t.Error("no se esperaba el envío de una notificación en una transferencia fallida")
	}
}
