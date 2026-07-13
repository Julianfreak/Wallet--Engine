package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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
)

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	fmt.Println("--- Iniciando Billetera Digital con PostgreSQL ---")

	// 1. Cargar Configuración desde Viper
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error al cargar la configuración: %v", err)
	}

	// 2. Ejecutar Migraciones (Nueva forma)
	// Usamos la URL de conexión completa que viene de cfg.DBSource
	fmt.Println("Verificando/Ejecutando migraciones...")
	if err := repository.RunMigrations(cfg.DBSource); err != nil {
		log.Fatalf("No se pudieron aplicar las migraciones: %v", err)
	}
	fmt.Println("Estructura de base de datos lista.")

	// 3. Abrir Conexión a la DB
	db, err := sql.Open("postgres", cfg.DBSource)
	if err != nil {
		log.Fatalf("Error al abrir la conexión: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Base de datos inaccesible: %v", err)
	}
	fmt.Println("Conexión exitosa a PostgreSQL.")

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
	accountRepo.Save(ctx, &domain.Account{ID: "A1", Owner: "Julian", Balance: 1000})
	accountRepo.Save(ctx, &domain.Account{ID: "A2", Owner: "Mercado Libre", Balance: 0})

	transferService := application.NewTransferService(accountRepo, transactionRepo, txManager, consoleLogger, emailSender)
	transferHandler := httpAdapter.NewTransferHandler(transferService)

	// Rutas
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/transfers", transferHandler.HandleTransfer)
	http.HandleFunc("/accounts", accountHandler.GetAccount)

	fmt.Printf("Servidor bancario escuchando en %s...\n", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, nil); err != nil {
		log.Fatalf("Error al encender el servidor web: %v", err)
	}
}
