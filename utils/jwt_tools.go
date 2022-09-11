package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
)

var JWTSalt = os.Getenv("JWTSalt")

func GenerateJWT(username string, passwordHash string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["password_hash"] = passwordHash
	return token.SignedString([]byte(JWTSalt))
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing token")
		}
		return []byte(JWTSalt), nil
	})
}
