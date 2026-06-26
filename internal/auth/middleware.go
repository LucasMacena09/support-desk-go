package auth

import (
	"context"
	"net/http"
	"strings"

	"support-desk-go/internal/httperr"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Middleware(jwtManager *JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httperr.Unauthorized(w, "token de autenticação não enviado")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			httperr.Unauthorized(w, "formato de token inválido, use: Bearer <token>")
			return
		}

		claims, err := jwtManager.Validate(parts[1])
		if err != nil {
			httperr.Unauthorized(w, "token inválido ou expirado")
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			httperr.Unauthorized(w, "token inválido")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}