package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	httpAdapter "github.com/Julianfreak/Wallet--Engine/internal/adapters/http"
	"github.com/Julianfreak/Wallet--Engine/internal/adapters/logger"
	"github.com/Julianfreak/Wallet--Engine/internal/adapters/notification"
	"github.com/Julianfreak/Wallet--Engine/internal/adapters/repository"
	"github.com/Julianfreak/Wallet--Engine/internal/application"
	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()
	if err != nil {
		fmt.Println("ℹNo se encontró el archivo .env, usando variables del sistema")
	} else {
		fmt.Println("Variables de configuración cargadas desde el archivo .env con éxito.")
	}

	dbUser := getEnv("DB_USER", "wallet_user")
	dbPass := getEnv("DB_PASSWORD", "wallet_password")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "wallet_db")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al abrir la conexión: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Base de datos inaccesible: %v", err)
	}
	fmt.Printf("Conexión exitosa a PostgreSQL en %s:%s.\n", dbHost, dbPort)

	// Migraciones autoimaticas
	fmt.Println("Ejecutando migraciones automáticas...")
	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("No se pudieron aplicar las migraciones: %v", err)
	}
	fmt.Println("Estructura de base de datos verificada/creada.")

	// Crear el contexto inicial de la aplicación
	ctx := context.Background()

	// INICIALIZAR EL CONTROLADOR DE TRANSACCIONES Y ADAPTADORES
	txManager := repository.NewPostgresTxManager(db)
	accountRepo := repository.NewPostgresAccountRepository(db)
	transactionRepo := repository.NewPostgresTransactionRepository(db)
	consoleLogger := logger.NewConsoleLogger()
	emailSender := notification.NewEmailSender()

	fmt.Println("Sembrando datos iniciales en Postgres...")
	accountRepo.Save(ctx, &domain.Account{ID: "A1", Owner: "Julian", Balance: 1000.0})
	accountRepo.Save(ctx, &domain.Account{ID: "A2", Owner: "Mercado Libre", Balance: 0.0})

	// INYECCIÓN DE DEPENDENCIAS CON EL TX_MANAGER INCLUIDO
	transferService := application.NewTransferService(accountRepo, transactionRepo, txManager, consoleLogger, emailSender)

	transferHandler := httpAdapter.NewTransferHandler(transferService)

	// Registramos la ruta HTTP y asociamos su función manejadora
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/transfers", transferHandler.HandleTransfer)

	port := ":8082"
	fmt.Printf("Servidor bancario escuchando en el puerto %s...\n", port)
	fmt.Println("Presiona Ctrl+C para apagar el servidor")

	// Encendemos el servidor. Este método es bloqueante; mantendrá la app viva.
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error al encender el servidor web: %v", err)
	}
}
