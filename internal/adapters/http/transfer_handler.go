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
func (h *TransferHandler) HandleTransfer(w http.ResponseWriter, r *http.Request) {
	// 1. Validar que estrictamente sea un método POST
	if r.Method != http.MethodPost {
		h.respondWithError(w, http.StatusMethodNotAllowed, "método no permitido, usa POST")
		return
	}

	// 2. Decodificar el cuerpo JSON de la petición
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "formato JSON inválido")
		return
	}

	if err := validate.Struct(req); err != nil {
		// Si rompe las reglas, devolvemos un HTTP 400 (Bad Request) y abortamos.
		// El código NO llega a tocar la lógica de negocio ni la base de datos.
		errMsg := fmt.Sprintf("Datos inválidos: %v", err)
		// Usamos tu método personalizado para mantener la consistencia
		h.respondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	// 3. Validación básica de entrada (Defensa básica)
	if req.FromAccountID == "" || req.ToAccountID == "" {
		h.respondWithError(w, http.StatusBadRequest, "los campos 'from_account_id' y 'to_account_id' son obligatorios")
		return
	}
	if req.Amount <= 0 {
		h.respondWithError(w, http.StatusBadRequest, "el monto a transferir debe ser mayor a cero")
		return
	}

	// 4. Invocar el núcleo de tu negocio (El caso de uso)
	ctx := r.Context()
	err := h.service.Execute(ctx, req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		// Si el servicio falla (por ejemplo, fondos insuficientes), respondemos con un 400 o 500 según corresponda
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 5. Responder con éxito si todo salió perfecto
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(TransferResponse{
		Message: "Transferencia procesada de forma atómica y auditable con éxito",
	})
}

// Función auxiliar para responder errores en formato JSON limpio
func (h *TransferHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
