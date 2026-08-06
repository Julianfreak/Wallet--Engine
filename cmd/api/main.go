package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	userRepo := repository.NewPostgresUserRepository(db)
	jwtSecret := "super_secret_wallet_key_2026"
	authService := application.NewAuthService(userRepo, jwtSecret)
	authHandler := httpAdapter.NewAuthHandler(authService)
	consoleLogger := logger.NewConsoleLogger()
	emailSender := notification.NewEmailSender()

	// Sembrado de datos iniciales (Solo para pruebas iniciales)
	_ = accountRepo.Save(ctx, &domain.Account{ID: "A1", Owner: "Julian", Balance: 1000})
	_ = accountRepo.Save(ctx, &domain.Account{ID: "A2", Owner: "Mercado Libre", Balance: 0})

	transferService := application.NewTransferService(accountRepo, transactionRepo, txManager, consoleLogger, emailSender)
	transferHandler := httpAdapter.NewTransferHandler(transferService)

	// Rutas
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/transfers", transferHandler.HandleTransfer)
	http.HandleFunc("/accounts", accountHandler.GetAccount)

	http.HandleFunc("/register", enableCORS(authHandler.HandleRegister))
	http.HandleFunc("/login", enableCORS(authHandler.HandleLogin))

	// NUEVA RUTA: La página web de Swagger
	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"),
	))

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: nil, // Utiliza el DefaultServeMux donde registraste los endpoints
	}

	// Canal para capturar señales del sistema operativo (Ctrl+C o Docker stop)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Arrancamos el servidor en una Goroutine paralela para que no bloquee el hilo principal
	go func() {
		slog.Info("servidor de Wallet-Engine iniciado con éxito", "address", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error crítico en el servidor web", "error", err)
			os.Exit(1)
		}
	}()

	// El hilo principal se detiene aquí hasta recibir una señal de apagado
	sig := <-sigChan
	slog.Info("señal de apagado recibida, iniciando Graceful Shutdown...", "signal", sig.String())

	// Creamos un contexto con los 7 segundos de gracia solicitados
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	// 1. Apagamos el servidor HTTP de forma ordenada (espera peticiones en vuelo)
	if err := server.Shutdown(ctxShutdown); err != nil {
		slog.Error("el servidor se detuvo de forma forzosa por timeout o error", "error", err)
	} else {
		slog.Info("servidor HTTP detenido ordenadamente")
	}

	// 2. Cerramos limpiamente la conexión a la base de datos
	if err := db.Close(); err != nil {
		slog.Error("error al cerrar la conexión a la base de datos", "error", err)
	} else {
		slog.Info("conexión a la base de datos cerrada correctamente")
	}

	slog.Info("aplicación finalizada de forma limpia")
}

// Función middleware para habilitar CORS
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Permitir peticiones desde tu frontend de React
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Si es una petición de preflight (OPTIONS), respondemos inmediatamente con OK
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
