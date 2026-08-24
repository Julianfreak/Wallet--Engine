package http

import (
	"encoding/json"
	"net/http"

	"github.com/Julianfreak/Wallet--Engine/internal/application"
)

type UserHandler struct {
	service *application.UserService
}

func NewUserHandler(service *application.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// HandleProfile enruta según el método HTTP (GET para consultar, PUT para actualizar)
func (h *UserHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	// Extraer el account_id inyectado por el AuthMiddleware
	accountID, ok := r.Context().Value(AccountIDKey).(string)
	if !ok || accountID == "" {
		http.Error(w, `{"error": "no autorizado"}`, http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := h.service.GetProfile(r.Context(), accountID)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user)

	case http.MethodPut:
		var req struct {
			Email string `json:"email" validate:"required,email"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "formato JSON inválido"}`, http.StatusBadRequest)
			return
		}

		err := h.service.UpdateProfile(r.Context(), accountID, req.Email)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "perfil actualizado exitosamente"})

	default:
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
	}
}
