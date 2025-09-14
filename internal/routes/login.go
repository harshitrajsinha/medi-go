package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/harshitrajsinha/medi-go/internal/models"
	"github.com/harshitrajsinha/medi-go/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type LoginRoutes struct {
	service *store.Store
}

func NewLoginRoutes(service *store.Store) *LoginRoutes {
	return &LoginRoutes{
		service: service,
	}
}

type CustomClaims struct {
	Email  string    `json:"email"`
	UserID uuid.UUID `json:"userid"`
	jwt.StandardClaims
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func GenerateToken(email string, userId uuid.UUID) (string, error) {

	expiration := time.Now().Add(30 * time.Minute) // Expiration set as 30 minute

	claims := &CustomClaims{
		Email:  email,
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
			IssuedAt:  time.Now().Unix(),
			Subject:   email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Load JWT key
	jwtKeyString := os.Getenv("JWT_KEY")

	signedToken, err := token.SignedString([]byte(jwtKeyString))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// POST: Return auth token based on credentials
func (l *LoginRoutes) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid Request body for authentication"})
		log.Println("Invalid Request body for authentication")
		return
	}

	credentials.Email = strings.TrimSpace(credentials.Email)
	credentials.Password = strings.TrimSpace(credentials.Password)
	credentials.Role = strings.TrimSpace(credentials.Role)

	if credentials.Email == "" || credentials.Password == "" || credentials.Role == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid Request body for authentication - email, password, role are mandatory fields"})
		log.Println("Invalid Request body for authentication, either of email, password, role are missing")
		return
	}

	loginResponse, err := l.service.GetLoginInfo(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error occured during authentication"})
		panic(err)
	}

	// Verifying a password
	err = bcrypt.CompareHashAndPassword([]byte(loginResponse.HashedPassword), []byte(credentials.Password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Incorrect email or password for authentication"})
		log.Println("Incorrect username or password for authentication")
		return
	}

	// Generate JWT token for authentication
	tokenString, err := GenerateToken(credentials.Email, loginResponse.UserID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Failed to generate token for authentication"})
		log.Println("Failed to generate token for authentication")
		panic(err)
	}
	resp := map[string]string{"token": tokenString}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Code: http.StatusCreated, Message: "Authentication token generated successfully. Valid for next 30mins", Data: resp})

	log.Println("Authentication token generated successfully")
}
