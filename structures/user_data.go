package structures

import "github.com/golang-jwt/jwt"

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (r UserData) FromToken(t *jwt.Token) UserData {
	claims := t.Claims.(jwt.MapClaims)
	userData := UserData{
		Username: claims["username"].(string),
		Password: claims["password"].(string),
		Database: claims["database"].(string),
	}
	return userData
}
