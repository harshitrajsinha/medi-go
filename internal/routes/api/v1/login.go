package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"

	"github.com/harshitrajsinha/medi-go/internal/auth"
	"github.com/harshitrajsinha/medi-go/internal/models"
)

// POST: Return auth token based on credentials
func (l *APIRoutes) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid Request body for authentication"})
		log.Println("Invalid Request body for authentication")
		return
	}

	credentials.Email = strings.TrimSpace(credentials.Email)
	credentials.Password = strings.TrimSpace(credentials.Password)
	credentials.Role = strings.TrimSpace(credentials.Role)

	if credentials.Email == "" || credentials.Password == "" || credentials.Role == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid Request body for authentication - email, password, role are mandatory fields"})
		log.Println("Invalid Request body for authentication, either of email, password, role are missing")
		return
	}

	// Regex test for email
	pattern := `^[a-zA-Z0-9.]+@medi\.go$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(credentials.Email) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid email format"})
		log.Println("Invalid email format")
		return
	}

	loginResponse, err := l.service.GetLoginInfo(&credentials)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured during authentication"})
		panic(err)
	}

	// Verifying a password
	err = auth.CheckPassowrd(loginResponse.HashedPassword, credentials.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Incorrect email or password for authentication"})
		log.Println("Incorrect username or password for authentication")
		return
	}

	// Generate JWT token for authentication
	tokenString, err := auth.GenerateToken(credentials.Email, loginResponse.UserID)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Failed to generate token for authentication"})
		log.Println("Failed to generate token for authentication")
		panic(err)
	}
	resp := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Code: http.StatusCreated, Message: "Authentication token generated successfully. Valid for next 30mins", Data: resp})

	log.Println("Authentication token generated successfully")
}
