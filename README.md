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

go test -v ./...

[![Billetera Digital CI](https://github.com/Julianfreak/Wallet--Engine/actions/workflows/ci.yml/badge.svg)](https://github.com/Julianfreak/Wallet--Engine/actions/workflows/ci.yml)

### Ejecución de Pruebas

Para ejecutar las pruebas unitarias e de integración (que requieren el contenedor de base de datos activo):

1. Levanta el servicio de base de datos con Docker Compose:
   ```bash
   docker-compose up -d wallet-db

##  Calidad de Código y Análisis Estático

   El proyecto utiliza `golangci-lint` para asegurar las buenas prácticas y evitar errores comunes. 

   Para ejecutar el análisis estático localmente mediante Docker:

   docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run   

##  Auditoría de Seguridad

El proyecto utiliza `govulncheck` (la herramienta oficial de Go) para escanear vulnerabilidades conocidas tanto en el código fuente como en las dependencias de terceros.

   Para ejecutar la auditoría de seguridad localmente:

   
   go install golang.org/x/vuln/cmd/govulncheck@latest
   govulncheck ./...


## Observabilidad y Logging

Wallet-Engine implementa **Logging Estructurado en formato JSON** utilizando el paquete nativo de la biblioteca estándar de Go (`log/slog`). 

*   **Formato JSON:** Todos los eventos del sistema se emiten estructurados con pares clave-valor, facilitando su indexación y auditoría en plataformas de análisis de registros (como ELK, Datadog o Grafana Loki).
*   **Niveles de Severidad:** Manejo estricto de niveles (`INFO`, `ERROR`) optimizados para entornos de producción en contenedores Docker.

## 🔐 Autenticación y Seguridad (JWT)

El motor ahora cuenta con un sistema de autenticación seguro basado en **JSON Web Tokens (JWT)** y contraseñas cifradas con **Bcrypt**.

### Endpoints Disponibles

*   **Registrar Usuario**
    *   `POST /register`
    *   *Body JSON:*
        ```json
        {
          "id": "U1",
          "email": "usuario@correo.com",
          "password": "tu_password_segura"
        }
        ```

*   **Iniciar Sesión (Login)**
    *   `POST /login`
    *   *Body JSON:*
        ```json
        {
          "email": "usuario@correo.com",
          "password": "tu_password_segura"
        }
        ```
    *   *Respuesta:* Retorna un token JWT válido por 24 horas que debe ser enviado en las peticiones protegidas mediante el header `Authorization: Bearer <token>`.