package http

import (
	"encoding/json"
	"net/http"

	"github.com/Julianfreak/Wallet--Engine/internal/adapters/repository"
	"github.com/Julianfreak/Wallet--Engine/internal/domain"
)

/* type contextKey string

const AccountIDKey contextKey = "user_id" */

type AccountHandler struct {
	repo *repository.PostgresAccountRepository
}

func NewAccountHandler(repo *repository.PostgresAccountRepository) *AccountHandler {
	return &AccountHandler{repo: repo}
}

// @Summary Obtiene una cuenta
// @Description Devuelve los detalles y saldo actual de una cuenta mediante su ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id query string true "ID de la cuenta"
// @Success 200 {string} string "Respuesta exitosa"
// @Router /accounts [get]
func (h *AccountHandler) GetAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "El parámetro 'id' es requerido", http.StatusBadRequest)
		return
	}

	account, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(account)
}

// HandleGetAccounts lista todas las cuentas disponibles para transferencias (excluyendo al usuario actual)
func (h *AccountHandler) HandleGetAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraer el correo del usuario autenticado desde el contexto del middleware
	userEmail, ok := r.Context().Value(AccountIDKey).(string)
	if !ok || userEmail == "" {
		http.Error(w, `{"error": "no autorizado"}`, http.StatusUnauthorized)
		return
	}

	accounts, err := h.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, `{"error": "error al obtener cuentas"}`, http.StatusInternalServerError)
		return
	}

	// Filtrar para excluir las cuentas propias del usuario en sesión
	var availableAccounts []domain.Account
	for _, acc := range accounts {
		if acc.Owner != userEmail {
			availableAccounts = append(availableAccounts, acc)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(availableAccounts)
}
