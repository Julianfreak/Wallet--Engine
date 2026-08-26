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

func (h *UserHandler) HandleProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(AccountIDKey).(string)
	if !ok || userID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "no autorizado"})
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := h.service.GetProfile(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		_ = json.NewEncoder(w).Encode(user)

	case http.MethodPut:
		var req struct {
			Email       string `json:"email"`
			NewPassword string `json:"new_password,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "formato JSON inválido"})
			return
		}

		err := h.service.UpdateProfile(r.Context(), userID, req.Email, req.NewPassword)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]string{"message": "perfil actualizado exitosamente"})

	case http.MethodDelete:
		err := h.service.DeleteProfile(r.Context(), userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		_ = json.NewEncoder(w).Encode(map[string]string{"message": "cuenta eliminada exitosamente"})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "método no permitido"})
	}
}
