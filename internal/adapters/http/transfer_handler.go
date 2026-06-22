package http

import (
	"encoding/json"
	"net/http"

	"github.com/Julianfreak/Wallet--Engine/internal/application"
)

// TransferRequest define la estructura del JSON que esperamos recibir del cliente
type TransferRequest struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
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
