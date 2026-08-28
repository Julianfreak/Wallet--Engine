package http

import (
	"encoding/json"
	"net/http"

	"github.com/Julianfreak/Wallet--Engine/internal/application"
)

type DashboardHandler struct {
	service *application.DashboardService
}

func NewDashboardHandler(service *application.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// HandleDashboard procesa la petición GET /dashboard
func (h *DashboardHandler) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	// 1. Extraer el correo del usuario autenticado desde el contexto del middleware
	userEmail, ok := r.Context().Value(AccountIDKey).(string)
	if !ok || userEmail == "" {
		http.Error(w, `{"error": "no autorizado"}`, http.StatusUnauthorized)
		return
	}

	// 2. Pasar el userEmail como segundo argumento al servicio
	data, err := h.service.GetDashboardData(r.Context(), userEmail)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
