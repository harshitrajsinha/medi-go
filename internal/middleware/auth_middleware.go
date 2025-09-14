package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func generateKey() []byte {

	jwtKeyString := os.Getenv("JWT_KEY")
	return []byte(jwtKeyString)
}

// Middleware to validate authentication token from request

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{
				Code: http.StatusUnauthorized, Message: "Authorization header required",
			})
			log.Println("Authorization header missing")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return generateKey(), nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{
				Code: http.StatusUnauthorized, Message: "Invalid token",
			})
			log.Println("Invalid token")
			return
		}
		ctx := context.WithValue(r.Context(), "email", claims.Email)
		next.ServeHTTP(w, r.WithContext(ctx))

	})

}
