package ports

// Logger es el puerto (interfaz) para el sistema de auditoría y rastreo
type Logger interface {
	Info(message string)
}
