package main

import (
	"fmt"
	"log"

	"github.com/Julianfreak/Wallet--Engine/internal/adapters/repository"
	"github.com/Julianfreak/Wallet--Engine/internal/application"
	"github.com/Julianfreak/Wallet--Engine/internal/domain"
)

func main() {
	// 1. INICIALIZAR ADAPTADORES
	accountRepo := repository.NewAccountMemoryRepo()
	transactionRepo := repository.NewTransactionMemoryRepo()

	// Simulamos que ya existen dos cuentas en nuestra "base de datos"
	accountRepo.Save(&domain.Account{ID: "A1", Owner: "Julian", Balance: 1000.0})
	accountRepo.Save(&domain.Account{ID: "A2", Owner: "Mercado Libre", Balance: 0.0})

	// 2. INYECCIÓN DE DEPENDENCIAS
	// Pasamos los adaptadores al constructor del servicio.
	transferService := application.NewTransferService(accountRepo, transactionRepo)

	fmt.Println("--- Iniciando simulación de Billetera Digital ---")

	// 3. EJECUTAR CASO DE USO Y GESTIONAR ERRORES
	// Intentamos transferir $300 de Julian a Mercado Libre
	err := transferService.Execute("A1", "A2", 300.0)

	// En Go, los errores son valores.
	if err != nil {
		log.Fatalf("La transferencia falló: %v", err)
	}

	fmt.Println("¡Transferencia exitosa!")

	// 4. VERIFICAR RESULTADOS FINALES
	julianAcc, _ := accountRepo.GetByID("A1")
	mlAcc, _ := accountRepo.GetByID("A2")

	fmt.Printf("Saldo de %s: $%.2f\n", julianAcc.Owner, julianAcc.Balance)
	fmt.Printf("Saldo de %s: $%.2f\n", mlAcc.Owner, mlAcc.Balance)
}
