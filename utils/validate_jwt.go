package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"os"
)

var JWTSalt = os.Getenv("secret_key")

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error in parsing token")
		}
		return []byte(JWTSalt), nil
	})
}
