package middleware

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const TokenKey contextKey = "token"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		ctx := context.WithValue(r.Context(), TokenKey, tokenString)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
