package http

import (
	"encoding/json"
	"net/http"

	"github.com/Julianfreak/Wallet--Engine/internal/application"
)

type AuthHandler struct {
	authService *application.AuthService
}

func NewAuthHandler(authService *application.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Registrar un nuevo usuario
// @Description Crea un nuevo usuario y su cuenta asociada en la plataforma.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Datos del registro"
// @Success 201 {object} map[string]string "Usuario registrado exitosamente"
// @Failure 400 {object} map[string]string "Error en la solicitud"
// @Router /register [post]
func (h *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	err := h.authService.Register(r.Context(), req.ID, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"message": "usuario registrado exitosamente"})
}

// @Summary Iniciar sesión
// @Description Autentica a un usuario y devuelve un token JWT.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Credenciales de acceso"
// @Success 200 {object} map[string]string "Token JWT"
// @Failure 401 {object} map[string]string "Credenciales inválidas"
// @Router /login [post]
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "cuerpo de solicitud inválido", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"token": token})
}
