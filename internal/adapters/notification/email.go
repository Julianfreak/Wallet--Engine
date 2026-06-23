package notification

import (
	"fmt"
	"time"
)

// EmailSender es el adaptador que simula el envío de correos real
type EmailSender struct{}

func NewEmailSender() *EmailSender {
	return &EmailSender{}
}

// Send simula una conexión de red lenta de 3 segundos hacia un servidor de correos
func (s *EmailSender) Send(recipient string, message string) error {
	// ⏳ SIMULACIÓN DE LATENCIA: Aquí es donde la app normalmente se quedaría congelada
	time.Sleep(3 * time.Second)

	// \033[33m activa el color amarillo para diferenciarlo de la auditoría azul
	fmt.Printf("📧 \033[33m[EMAIL ENVIADO]\033[0m Notificación exitosa para %s: %s\n", recipient, message)
	return nil
}
