package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"passwordHash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
