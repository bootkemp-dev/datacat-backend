package models

type PasswordResetRequest struct {
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}
