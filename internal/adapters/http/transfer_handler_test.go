package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Julianfreak/Wallet--Engine/internal/application"
	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/Julianfreak/Wallet--Engine/internal/testutils"
)

func TestTransferHandler_ValidationShield(t *testing.T) {
	// 1. Arrange: Construimos el servicio con los Fakes que ya demostraron ser estables
	// (Asegúrate de que tus Fakes sean accesibles desde este paquete o muévelos a un paquete de pruebas común)
	fakeAccountRepo := &testutils.FakeAccountRepository{
		Accounts: map[string]*domain.Account{
			"A1": {ID: "A1", Owner: "Julian", Balance: 1000.0},
			"A2": {ID: "A2", Owner: "Mercado Libre", Balance: 0.0},
		},
	}

	fakeTxRepo := &testutils.FakeTransactionRepository{}
	fakeTxManager := &testutils.FakeTxManager{}
	fakeLogger := &testutils.FakeLogger{}
	fakeNotifier := &testutils.FakeNotificationSender{
		Done: make(chan bool, 1),
	}
	service := application.NewTransferService(fakeAccountRepo, fakeTxRepo, fakeTxManager, fakeLogger, fakeNotifier)
	handler := NewTransferHandler(service)

	t.Run("Rechazar monto negativo", func(t *testing.T) {
		payload := []byte(`{"from_account_id": "A1", "to_account_id": "A2", "amount": -10.0}`)
		req := httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(payload))
		rec := httptest.NewRecorder()

		handler.HandleTransfer(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Se esperaba 400 Bad Request, se obtuvo: %d", rec.Code)
		}
	})

	t.Run("Rechazar autotransferencia", func(t *testing.T) {
		payload := []byte(`{"from_account_id": "A1", "to_account_id": "A1", "amount": 50.0}`)
		req := httptest.NewRequest(http.MethodPost, "/transfers", bytes.NewBuffer(payload))
		rec := httptest.NewRecorder()

		handler.HandleTransfer(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Se esperaba 400 Bad Request (autotransferencia bloqueada), se obtuvo: %d", rec.Code)
		}
	})
}
