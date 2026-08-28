package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// En producción, esta clave debe venir de variables de entorno (.env)
var jwtSecret = []byte("super_secret_wallet_key_2026")

// GenerateToken crea un JWT válido por 24 horas para un accountID específico
func GenerateToken(email string) (string, error) {
	claims := jwt.MapClaims{
		"account_id": email,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret) // Firma con jwtSecret
}

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil // Valida con el MISMO jwtSecret
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token inválido o expirado")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claims del token inválidos")
	}

	return claims, nil
}
