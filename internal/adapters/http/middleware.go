package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Julianfreak/Wallet--Engine/internal/infrastructure/auth"
)

type contextKey string

const AccountIDKey contextKey = "user_id"

// AuthMiddleware protege las rutas verificando el token JWT
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "autorización requerida"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": "formato de token inválido"}`, http.StatusUnauthorized)
			return
		}

		claims, err := auth.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}

		// Intentamos extraer el email; si no existe, usamos el account_id
		var identifier string
		if email, ok := claims["email"].(string); ok && email != "" {
			identifier = email
		} else if accID, ok := claims["account_id"].(string); ok && accID != "" {
			identifier = accID
		} else if accIDFloat, ok := claims["account_id"].(float64); ok {
			identifier = fmt.Sprintf("%.0f", accIDFloat)
		}

		if identifier == "" {
			http.Error(w, `{"error": "token sin identificador válido"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AccountIDKey, identifier)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
