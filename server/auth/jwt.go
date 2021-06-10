package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/bootkemp-dev/datacat-backend/config"
	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID       int    `json:"id"`
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

func GenerateToken(username string, id int) (string, *time.Time, error) {

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		ID:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    c.Jwt.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.Jwt.JwtKey))
	if err != nil {
		return "", nil, err
	}

	return tokenString, &expirationTime, nil
}

func isTokenValid(tokenString string) (string, int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("Token signing method is not valid: %v", token.Header["alg"])
		}

		return []byte(c.Jwt.JwtKey), nil
	})
	if err != nil {
		log.Println(err)
		return "", 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idFloat := claims["id"]
		id := int(idFloat.(float64))
		username := claims["username"]
		return username.(string), id, nil
	} else {
		return "", 0, fmt.Errorf("reading claims failed")
	}

}
