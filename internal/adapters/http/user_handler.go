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

// HandleProfile gestiona el perfil del usuario (GET, PUT, DELETE)
// @Summary Obtener perfil del usuario
// @Description Devuelve el perfil del usuario autenticado.
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.User
// @Failure 401 {object} map[string]string "No autorizado"
// @Router /profile [get]
//
// @Summary Actualizar perfil del usuario
// @Description Actualiza el email o contraseña del usuario autenticado.
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]string true "Nuevos datos (email, new_password)"
// @Success 200 {object} map[string]string "Perfil actualizado exitosamente"
// @Failure 400 {object} map[string]string "Formato inválido"
// @Router /profile [put]
//
// @Summary Eliminar cuenta de usuario
// @Description Elimina de forma permanente el perfil del usuario autenticado.
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "Cuenta eliminada exitosamente"
// @Failure 401 {object} map[string]string "No autorizado"
// @Router /profile [delete]
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
