package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// 1. Contador para el total de transferencias procesadas
	TransfersTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "wallet_transfers_processed_total",
			Help: "El numero total de transferencias procesadas por el motor",
		},
		[]string{"status"}, // Nos permitira filtrar por "success" o "error"
	)

	// 2. Histograma para medir cuanto tarda en milisegundos la base de datos y la logica
	TransferDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "wallet_transfer_duration_seconds",
			Help:    "Tiempo de ejecucion de la transferencia en segundos",
			Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 3.0}, // Rangos de tiempo
		},
	)
)
