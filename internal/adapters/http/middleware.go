package http

import (
	"context"
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

		// Extraemos el correo con seguridad absoluta (es un string garantizado)
		email, ok := claims["email"].(string)
		if !ok || email == "" {
			http.Error(w, `{"error": "token sin correo válido"}`, http.StatusUnauthorized)
			return
		}

		// Inyectamos el email en el contexto
		ctx := context.WithValue(r.Context(), AccountIDKey, email)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
