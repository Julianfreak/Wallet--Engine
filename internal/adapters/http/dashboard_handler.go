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
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "método no permitido, usa GET"})
		return
	}

	ctx := r.Context()
	data, err := h.service.GetDashboardData(ctx)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data)
}
