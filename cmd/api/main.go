package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Julianfreak/Wallet--Engine/internal/adapters/repository"
	"github.com/Julianfreak/Wallet--Engine/internal/application"
	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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
	fmt.Printf("🔌 Conexión exitosa a PostgreSQL en %s:%s.\n", dbHost, dbPort)

	// Crear el contexto inicial de la aplicación
	ctx := context.Background()

	// INICIALIZAR EL CONTROLADOR DE TRANSACCIONES Y ADAPTADORES
	txManager := repository.NewPostgresTxManager(db)
	accountRepo := repository.NewPostgresAccountRepository(db)
	transactionRepo := repository.NewPostgresTransactionRepository(db)

	fmt.Println("Sembrando datos iniciales en Postgres...")
	accountRepo.Save(ctx, &domain.Account{ID: "A1", Owner: "Julian", Balance: 1000.0})
	accountRepo.Save(ctx, &domain.Account{ID: "A2", Owner: "Mercado Libre", Balance: 0.0})

	// INYECCIÓN DE DEPENDENCIAS CON EL TX_MANAGER INCLUIDO
	transferService := application.NewTransferService(accountRepo, transactionRepo, txManager)

	fmt.Println("Procesando transferencia de $300 de Julian a Mercado Libre con protección ACID...")
	err = transferService.Execute(ctx, "A1", "A2", 300.0)
	if err != nil {
		log.Fatalf("La transferencia falló: %v", err)
	}
	fmt.Println("¡Transferencia procesada y blindada con éxito!")

	// CONSULTAR LOS SALDOS FINALES
	julianAcc, _ := accountRepo.GetByID(ctx, "A1")
	mlAcc, _ := accountRepo.GetByID(ctx, "A2")

	fmt.Printf("\n Saldos finales en la Base de Datos:\n")
	fmt.Printf("- %s: $%.2f\n", julianAcc.Owner, julianAcc.Balance)
	fmt.Printf("- %s: $%.2f\n", mlAcc.Owner, mlAcc.Balance)
}
