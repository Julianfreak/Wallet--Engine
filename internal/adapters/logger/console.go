package logger

import (
	"fmt"
	"time"
)

// ConsoleLogger es el adaptador que implementa ports.Logger escribiendo en la terminal
type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

// Info escribe el mensaje en la consola usando colores ANSI para la terminal de Ubuntu
func (l *ConsoleLogger) Info(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// \033[34m activa el color azul, \033[0m lo resetea
	fmt.Printf("🔍 \033[34m[%s] AUDITORÍA:\033[0m %s\n", timestamp, message)
}
