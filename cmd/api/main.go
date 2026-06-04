package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Julianfreak/Wallet--Engine/internal/adapters/repository"
	"github.com/Julianfreak/Wallet--Engine/internal/application"
	"github.com/Julianfreak/Wallet--Engine/internal/domain"
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

	dbUser := getEnv("DB_USER", "wallet_user")
	dbPass := getEnv("DB_PASSWORD", "wallet_password")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "wallet_db")

	// 1. CONECTAR A LA BASE DE DATOS REAL (Docker)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al abrir la conexión: %v", err)
	}
	defer db.Close()

	// Verificar que la base de datos realmente responda
	if err := db.Ping(); err != nil {
		log.Fatalf("Base de datos inaccesible: %v", err)
	}
	fmt.Printf("Conexión exitosa a PostgreSQL en %s:%s.\n", dbHost, dbPort)

	// 2. INICIALIZAR LOS ADAPTADORES REALES
	accountRepo := repository.NewPostgresAccountRepository(db)
	transactionRepo := repository.NewPostgresTransactionRepository(db)

	// Simulamos la creación/actualización de las cuentas directamente en las tablas SQL
	fmt.Println("Sembrando datos iniciales en Postgres...")
	accountRepo.Save(&domain.Account{ID: "A1", Owner: "Julian", Balance: 1000.0})
	accountRepo.Save(&domain.Account{ID: "A2", Owner: "Mercado Libre", Balance: 0.0})

	// 3. INYECCIÓN DE DEPENDENCIAS (El Switch)
	// El servicio recibe los adaptadores de Postgres sin protestar, porque cumplen el contrato.
	transferService := application.NewTransferService(accountRepo, transactionRepo)

	// 4. EJECUTAR TRANSFERENCIA REAL
	fmt.Println("Procesando transferencia de $300 de Julian a Mercado Libre...")
	err = transferService.Execute("A1", "A2", 300.0)
	if err != nil {
		log.Fatalf("La transferencia falló: %v", err)
	}
	fmt.Println("¡Transferencia procesada y guardada con éxito en los discos de Postgres!")

	// 5. CONSULTAR LOS SALDOS FINALES DIRECTO DESDE POSTGRES
	julianAcc, _ := accountRepo.GetByID("A1")
	mlAcc, _ := accountRepo.GetByID("A2")

	fmt.Printf("\n Saldos finales en la Base de Datos:\n")
	fmt.Printf("- %s: $%.2f\n", julianAcc.Owner, julianAcc.Balance)
	fmt.Printf("- %s: $%.2f\n", mlAcc.Owner, mlAcc.Balance)
}
