package models

type Credentials struct {
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
