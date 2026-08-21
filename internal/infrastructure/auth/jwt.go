package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// En producción, esta clave debe venir de variables de entorno (.env)
var SecretKey = []byte("super_secreto_wallet_engine_2026")

// GenerateToken crea un JWT válido por 24 horas para un accountID específico
func GenerateToken(accountID string) (string, error) {
	claims := jwt.MapClaims{
		"account_id": accountID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// ValidateToken verifica la firma y retorna los datos (claims) si es válido
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token inválido o expirado")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claims inválidos")
	}

	return claims, nil
}
