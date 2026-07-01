# ==========================================
# ETAPA 1: Construcción (Builder)
# ==========================================
FROM golang:1.23.4-alpine AS builder

# Configuramos el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copiamos los archivos de dependencias y las descargamos (aprovecha la caché de Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copiamos el resto del código fuente
COPY . .

# Compilamos el binario nativo de Linux (CGO_ENABLED=0 lo hace ultra portable)
RUN CGO_ENABLED=0 GOOS=linux go build -o wallet-api ./cmd/api/main.go

# ==========================================
# ETAPA 2: Producción (Runner)
# ==========================================
FROM alpine:latest

# Añadimos certificados de seguridad por si el API necesita hacer peticiones HTTPS al exterior
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiamos SOLO el binario terminado desde la etapa 'builder'
COPY --from=builder /app/wallet-api .

# Exponemos el puerto de la aplicación
EXPOSE 8082

# Comando para arrancar el servidor
CMD ["./wallet-api"]