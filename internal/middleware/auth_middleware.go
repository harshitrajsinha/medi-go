package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/harshitrajsinha/medi-go/internal/auth"
)

// Middleware to validate authentication token from request

type Key string

const contextKey Key = "email"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response(w, http.StatusUnauthorized, "Authorization header required", "Authorization header missing")
			return
		}

		claims, err := auth.VerifyToken(authHeader)

		if err != nil {
			response(w, http.StatusUnauthorized, "Invalid token", "Invalid token")
			return
		}
		ctx := context.WithValue(r.Context(), contextKey, claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))

	})

}

func response(w http.ResponseWriter, code int, message string, logMessage string) {

	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}{
		Code: code, Message: message,
	})
	fmt.Println(logMessage)
}
