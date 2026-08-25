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
	// El contexto ahora trae el email del usuario
	email, ok := r.Context().Value(AccountIDKey).(string)
	if !ok || email == "" {
		http.Error(w, `{"error": "no autorizado"}`, http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		user, err := h.service.GetProfileByEmail(r.Context(), email)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(user)

	case http.MethodPut:
		var req struct {
			Email       string `json:"email"`
			NewPassword string `json:"new_password,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "formato JSON inválido"}`, http.StatusBadRequest)
			return
		}

		// Si el usuario no envió un nuevo correo, conservamos el actual
		targetEmail := req.Email
		if targetEmail == "" {
			targetEmail = email
		}

		err := h.service.UpdateProfileByEmail(r.Context(), email, targetEmail, req.NewPassword)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "perfil actualizado exitosamente"})

	case http.MethodDelete:
		err := h.service.DeleteProfileByEmail(r.Context(), email)
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "cuenta de usuario eliminada exitosamente"})

	default:
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
	}
}
