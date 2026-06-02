package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Julianfreak/Wallet--Engine/internal/adapters/repository"
	"github.com/Julianfreak/Wallet--Engine/internal/application"
	"github.com/Julianfreak/Wallet--Engine/internal/domain"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("--- Iniciando Billetera Digital con PostgreSQL ---")

	// 1. CONECTAR A LA BASE DE DATOS REAL (Docker)
	// Usamos las mismas credenciales que definimos en el docker-compose.yml
	connStr := "postgres://wallet_user:wallet_password@localhost:5432/wallet_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error al abrir la conexión: %v", err)
	}
	defer db.Close() // Nos aseguramos de cerrar la conexión al apagar la app

	// Verificar que la base de datos realmente responda
	if err := db.Ping(); err != nil {
		log.Fatalf("Base de datos inaccesible: %v", err)
	}
	fmt.Println("Conexión exitosa a PostgreSQL en Docker.")

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
