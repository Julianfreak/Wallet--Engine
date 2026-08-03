package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/Julianfreak/Wallet--Engine/docs"
	"github.com/Julianfreak/Wallet--Engine/internal/adapters/api"
	httpAdapter "github.com/Julianfreak/Wallet--Engine/internal/adapters/http"
	"github.com/Julianfreak/Wallet--Engine/internal/adapters/logger"
	"github.com/Julianfreak/Wallet--Engine/internal/adapters/notification"
	"github.com/Julianfreak/Wallet--Engine/internal/adapters/repository"
	"github.com/Julianfreak/Wallet--Engine/internal/application"
	"github.com/Julianfreak/Wallet--Engine/internal/config"
	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Wallet-Engine API
// @version 1.0
// @description Motor de billetera digital con transacciones atómicas.
// @host localhost:8082
// @BasePath /
func main() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(jsonHandler))

	slog.Info("--- Iniciando Billetera Digital con PostgreSQL ---")

	// 1. Cargar Configuración desde Viper
	cfg, err := config.LoadConfig(".")
	if err != nil {
		slog.Error("error al cargar la configuración", "error", err)
		os.Exit(1)
	}

	// 2. Ejecutar Migraciones (Nueva forma)
	// Usamos la URL de conexión completa que viene de cfg.DBSource
	slog.Info("Verificando/Ejecutando migraciones...")
	if err := repository.RunMigrations(cfg.DBSource); err != nil {
		slog.Error("error al ejecutar migraciones", "error", err)
		os.Exit(1)
	}
	slog.Info("Estructura de base de datos lista.")

	// 3. Abrir Conexión a la DB
	db, err := sql.Open("postgres", cfg.DBSource)
	if err != nil {
		slog.Error("error al abrir la conexión", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("error al cerrar la conexión a la base de datos", "error", err)
		}
	}()

	if err := db.Ping(); err != nil {
		slog.Error("error al pingear la base de datos", "error", err)
		os.Exit(1)
	}
	slog.Info("Conexión exitosa a PostgreSQL.")

	// Inicialización de adaptadores y servicios
	ctx := context.Background()
	txManager := repository.NewPostgresTxManager(db)
	accountRepo := repository.NewPostgresAccountRepository(db)
	accountHandler := api.NewAccountHandler(accountRepo)
	transactionRepo := repository.NewPostgresTransactionRepository(db)
	consoleLogger := logger.NewConsoleLogger()
	emailSender := notification.NewEmailSender()

	// Sembrado de datos iniciales (Solo para pruebas iniciales)
	// Nota: Esto podría fallar si la cuenta ya existe, considera validarlo
	_ = accountRepo.Save(ctx, &domain.Account{ID: "A1", Owner: "Julian", Balance: 1000})
	_ = accountRepo.Save(ctx, &domain.Account{ID: "A2", Owner: "Mercado Libre", Balance: 0})

	transferService := application.NewTransferService(accountRepo, transactionRepo, txManager, consoleLogger, emailSender)
	transferHandler := httpAdapter.NewTransferHandler(transferService)

	// Rutas
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/transfers", transferHandler.HandleTransfer)
	http.HandleFunc("/accounts", accountHandler.GetAccount)

	// NUEVA RUTA: La página web de Swagger
	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))

	slog.Info("Servidor bancario escuchando en %s...", cfg.ServerAddress)
	// En cmd/api/main.go
	slog.Info("Servidor de Wallet-Engine iniciado con éxito en http://%s \n", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, nil); err != nil {
		slog.Error("error al encender el servidor web", "error", err)
		os.Exit(1)
	}
}
