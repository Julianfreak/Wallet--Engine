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

##  Requisitos Previos

Asegúrate de tener instalado en tu sistema local:
*   [Go](https://golang.org/doc/install) (v1.20 o superior)
*   [PostgreSQL](https://www.postgresql.org/download/)
*   [Air](https://github.com/air-verse/air) (Para recarga automática en desarrollo)

## Instalación y Uso

1. **Clonar el repositorio:**
   ```bash
   git clone [https://github.com/Julianfreak/Wallet--Engine.git](https://github.com/Julianfreak/Wallet--Engine.git)
   cd Wallet--Engine

## Testing e Integración Continua (CI)

Este proyecto prioriza la fiabilidad mediante pruebas unitarias exhaustivas en la capa de adaptadores (Handlers) y lógica de negocio (Servicios), utilizando Mocks para aislar el comportamiento de la base de datos.

### Cobertura de Código (Code Coverage)
El proyecto utiliza pruebas basadas en tablas (Table-Driven Tests) y dobles de prueba (Fakes) para simular la infraestructura. Actualmente, la cobertura en la capa de adaptadores HTTP y la lógica de negocio supera el 75%.

### Ejecución Local
Para ejecutar la suite de pruebas completa en tu entorno local con salida detallada:

```bash
go test -v ./...

[![Billetera Digital CI](https://github.com/Julianfreak/Wallet--Engine/actions/workflows/ci.yml/badge.svg)](https://github.com/Julianfreak/Wallet--Engine/actions/workflows/ci.yml)

### Ejecución de Pruebas

Para ejecutar las pruebas unitarias e de integración (que requieren el contenedor de base de datos activo):

1. Levanta el servicio de base de datos con Docker Compose:
   ```bash
   docker-compose up -d wallet-db