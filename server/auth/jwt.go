package auth

import (
	"log"
	"time"

	"github.com/bootkemp-dev/datacat-backend/config"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, int64, error) {
	c, err := config.NewConfig("./config.yml")
	if err != nil {
		return "", 0, err
	}

	expirationTime := time.Now().Add(1 * time.Hour).Unix()
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
			IssuedAt:  time.Now().Unix(),
			Issuer:    c.Jwt.Issuer,
		},
	}

	log.Println("jwtKey: ", c.Jwt.JwtKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.Jwt.JwtKey))
	if err != nil {
		return "", 0, err
	}

	return tokenString, expirationTime, nil
}
