package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Julianfreak/Wallet--Engine/internal/application"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// TransferRequest define la estructura del JSON que esperamos recibir del cliente
type TransferRequest struct {
	FromAccountID string  `json:"from_account_id" validate:"required"`
	ToAccountID   string  `json:"to_account_id" validate:"required,nefield=FromAccountID"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

// TransferResponse define la estructura de la respuesta en caso de éxito
type TransferResponse struct {
	Message string `json:"message"`
}

// ErrorResponse define la estructura estándar para reportar fallos en formato JSON
type ErrorResponse struct {
	Error string `json:"error"`
}

// TransferHandler es nuestro adaptador primario HTTP
type TransferHandler struct {
	service *application.TransferService
}

func NewTransferHandler(service *application.TransferService) *TransferHandler {
	return &TransferHandler{service: service}
}

// HandleTransfer procesa la petición POST /transfers
// @Summary Realiza una transferencia
// @Description Transfiere fondos de una cuenta origen a una cuenta destino garantizando la atomicidad de la operación en base de datos.
// @Tags transfers
// @Accept json
// @Produce json
// @Param request body string true "JSON con origin_id, destination_id y amount"
// @Success 200 {string} string "Transferencia exitosa"
// @Failure 400 {string} string "Datos inválidos o fondos insuficientes"
// @Failure 500 {string} string "Error interno del servidor"
// @Router /transfers [post]
func (h *TransferHandler) HandleTransactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreateTransfer(w, r)
	case http.MethodGet:
		h.handleGetTransactions(w, r)
	default:
		h.respondWithError(w, http.StatusMethodNotAllowed, "método no permitido, usa GET o POST")
	}
}

// handleCreateTransfer procesa la lógica interna para el POST /transactions
func (h *TransferHandler) handleCreateTransfer(w http.ResponseWriter, r *http.Request) {

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	if err := validate.Struct(req); err != nil {
		errMsg := fmt.Sprintf("Datos inválidos: %v", err)
		h.respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	ctx := r.Context()
	cmd := application.TransferCommand{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	userEmail, ok := ctx.Value(AccountIDKey).(string)
	if !ok || userEmail == "" {
		h.respondWithError(w, http.StatusUnauthorized, "no autorizado")
		return
	}

	err := h.service.Execute(ctx, userEmail, cmd)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(TransferResponse{
		Message: "Transferencia procesada de forma atómica y auditable con éxito",
	})
}

// handleGetTransactions procesa la lógica interna para el GET /transactions
// En internal/adapters/http/transfer_handler.go
func (h *TransferHandler) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userEmail, ok := ctx.Value(AccountIDKey).(string)
	if !ok || userEmail == "" {
		h.respondWithError(w, http.StatusUnauthorized, "no autorizado")
		return
	}

	// Pasamos el correo del usuario autenticado al servicio
	transactions, err := h.service.GetTransactions(ctx, userEmail)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "error al obtener el historial de transacciones")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(transactions)
}

// Función auxiliar para responder errores en formato JSON limpio
func (h *TransferHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
