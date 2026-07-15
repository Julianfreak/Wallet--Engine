# Wallet Engine - Go Backend Senior Path , exercise
# Wallet-Engine 

Un motor de billetera digital robusto y escalable construido con Go. Diseñado para manejar transacciones financieras de forma segura, garantizando la atomicidad de las operaciones mediante PostgreSQL.

## Stack Tecnológico

*   **Lenguaje:** Go (Golang)
*   **Base de Datos:** PostgreSQL
*   **Herramientas de Desarrollo:** Air (Live reloading)
*   **Arquitectura:** Clean Architecture / Hexagonal (Separación por capas)

## Arquitectura del Proyecto

El proyecto sigue principios de diseño modular para garantizar que la lógica de negocio esté aislada de las dependencias externas (como la base de datos o el protocolo HTTP).

*   `cmd/api/`: Punto de entrada de la aplicación.
*   `internal/adapters/api/`: Controladores (Handlers) que manejan las peticiones HTTP.
*   `internal/domain/`: (Próximamente) Entidades y reglas de negocio.
*   `internal/repository/`: Interacción directa con la base de datos PostgreSQL.

## ⚙️ Requisitos Previos

Asegúrate de tener instalado en tu sistema local:
*   [Go](https://golang.org/doc/install) (v1.20 o superior)
*   [PostgreSQL](https://www.postgresql.org/download/)
*   [Air](https://github.com/air-verse/air) (Para recarga automática en desarrollo)

## Instalación y Uso

1. **Clonar el repositorio:**
   ```bash
   git clone [https://github.com/TU_USUARIO/Wallet--Engine.git](https://github.com/TU_USUARIO/Wallet--Engine.git)
   cd Wallet--Engine