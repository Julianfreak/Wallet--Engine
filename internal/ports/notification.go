package ports

// NotificationSender es el puerto para enviar alertas a los usuarios
type NotificationSender interface {
	Send(recipient string, message string) error
}
