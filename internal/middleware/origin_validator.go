package middleware

import (
	"fmt"
	"net/http"

	"github.com/harshitrajsinha/medi-go/config"
)

func OriginValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var allowedOrigin string

		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}
		origin := r.Header.Get("Origin")

		allowedOriginWebsite, err := config.AllowedOrigin()
		if err != nil {
			fmt.Println(err)
			allowedOrigin = ""
		}
		allowedOrigin = allowedOriginWebsite.AllowedOrigin

		// Reject if no Origin header or unregistered origin
		if origin == "" || allowedOrigin != origin {
			http.Error(w, "Unauthorized or restricted origin", http.StatusUnauthorized)
			return
		}

		// Setting CORS headers only for allowed origins
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Pass to the next handler
		next.ServeHTTP(w, r)
	})
}
