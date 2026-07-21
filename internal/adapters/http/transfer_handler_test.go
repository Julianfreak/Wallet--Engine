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
	fakeAccountRepo := &testutils.FakeAccountRepository{
		Accounts: map[string]*domain.Account{
			"A1":           {ID: "A1", Owner: "Julian", Balance: 1000.0},
			"A2":           {ID: "A2", Owner: "Mercado Libre", Balance: 0.0},
			"cuenta_pobre": {ID: "cuenta_pobre", Balance: 50.0}, // Añadimos esta para probar fondos insuficientes
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

	// 2. Definimos la tabla con TODOS los caminos tristes
	casos := []struct {
		nombre         string
		metodo         string
		cuerpoRequest  string
		codigoEsperado int
	}{
		{
			nombre:         "Rechazar monto negativo",
			metodo:         http.MethodPost,
			cuerpoRequest:  `{"from_account_id": "A1", "to_account_id": "A2", "amount": -10.0}`,
			codigoEsperado: http.StatusBadRequest,
		},
		{
			nombre:         "Rechazar autotransferencia",
			metodo:         http.MethodPost,
			cuerpoRequest:  `{"from_account_id": "A1", "to_account_id": "A1", "amount": 50.0}`,
			codigoEsperado: http.StatusBadRequest,
		},
		{
			nombre:         "Falla por método incorrecto (GET)",
			metodo:         http.MethodGet,
			cuerpoRequest:  `{"from_account_id": "A1", "to_account_id": "A2", "amount": 50.0}`,
			codigoEsperado: http.StatusMethodNotAllowed,
		},
		{
			nombre:         "Falla por JSON inválido",
			metodo:         http.MethodPost,
			cuerpoRequest:  `{"from_account_id": "A1", mal_formado}`,
			codigoEsperado: http.StatusBadRequest,
		},
		{
			nombre:         "Falla por cuentas vacías",
			metodo:         http.MethodPost,
			cuerpoRequest:  `{"from_account_id": "", "to_account_id": "", "amount": 50.0}`,
			codigoEsperado: http.StatusBadRequest,
		},
		{
			nombre:         "Falla en el servicio (Fondos insuficientes)",
			metodo:         http.MethodPost,
			cuerpoRequest:  `{"from_account_id": "cuenta_pobre", "to_account_id": "A2", "amount": 5000.0}`,
			codigoEsperado: http.StatusBadRequest,
		},
	}

	// 3. Iteramos sobre la tabla y ejecutamos cada caso
	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			payload := []byte(caso.cuerpoRequest)
			req := httptest.NewRequest(caso.metodo, "/transfers", bytes.NewBuffer(payload))
			rec := httptest.NewRecorder()

			handler.HandleTransfer(rec, req)

			if rec.Code != caso.codigoEsperado {
				t.Errorf("Se esperaba %d, se obtuvo: %d en el caso '%s'", caso.codigoEsperado, rec.Code, caso.nombre)
			}
		})
	}
}
