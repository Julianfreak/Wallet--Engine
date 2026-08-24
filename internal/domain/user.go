package domain

import "time"

// User representa a un usuario registrado en el sistema con credenciales de acceso.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}
