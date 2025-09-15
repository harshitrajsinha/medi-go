package auth

import (
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type CustomClaims struct {
	Email  string    `json:"email"`
	UserID uuid.UUID `json:"userid"`
	jwt.StandardClaims
}

func generateKey() []byte {
	jwtKeyString := os.Getenv("JWT_KEY")
	return []byte(jwtKeyString)
}

func CheckPassowrd(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
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

func VerifyToken(authHeader string) (*Claims, error) {

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return generateKey(), nil
	})

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}
