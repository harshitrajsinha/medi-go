package middleware

import (
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type allowedOrigins struct {
	AllowedOrigin string `envconfig:"ALLOWED_ORIGIN"`
}

func OriginValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}

		var originCfg allowedOrigins

		err := envconfig.Process("", &originCfg) // load env from program's environment to declared struct
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		origin := r.Header.Get("Origin")

		// Reject if no Origin header or unregistered origin
		if originCfg.AllowedOrigin != origin {
			http.Error(w, "Unauthorized or restricted origin", http.StatusUnauthorized)
			return
		}

		// Set CORS headers only for allowed origins
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
