package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/harshitrajsinha/medi-go/models"
	"github.com/harshitrajsinha/medi-go/store"
	"github.com/joho/godotenv"
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

func GenerateToken(email string) (string, error) {

	expiration := time.Now().Add(30 * time.Minute) // Expiration set as 30 minute

	claims := &jwt.StandardClaims{
		ExpiresAt: expiration.Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Load JWT key
	_ = godotenv.Load()
	jwtKeyString := os.Getenv("JWT_KEY")

	signedToken, err := token.SignedString([]byte(jwtKeyString))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (l *LoginRoutes) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code: http.StatusBadRequest, Message: "Invalid Request body for authorization",
		})
		log.Println("Invalid Request body for authorization")
		return
	}

	hashedPassword, err := l.service.GetLoginInfo(ctx, &credentials)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code: http.StatusInternalServerError, Message: "Error occured during authentication",
		})
		log.Println("Error occured during authentication")
		return
	}

	// Verifying a password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code: http.StatusBadRequest, Message: "Incorrect email or password for authorization",
		})
		log.Println("Incorrect username or password for authorization")
		return
	}

	tokenString, err := GenerateToken(credentials.Email)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code: http.StatusInternalServerError, Message: "Failed to generate token for authorization",
		})
		log.Println("Failed to generate token for authorization")
		return
	}
	response := make([]map[string]string, 0)
	response = append(response, map[string]string{"token": tokenString})
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Code    int                 `json:"code"`
		Message string              `json:"message"`
		Data    []map[string]string `json:"data"`
	}{
		Code: http.StatusCreated, Message: "Authorization token generated successfully. Valid for next 30mins", Data: response,
	})
	log.Println("Authorization token generated successfully")
}
