package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/bootkemp-dev/datacat-backend/config"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var c config.Config

func init() {
	config, err := config.NewConfig("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	c = *config
}

func GenerateToken(username string) (string, int64, error) {

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

func isTokenValid(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("Token signing method is not valid: %v", token.Header["alg"])
		}

		return []byte(c.Jwt.JwtKey), nil
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"]
		return username.(string), nil
	} else {
		return "", fmt.Errorf("reading claims failed")
	}

}
