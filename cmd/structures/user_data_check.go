package structures

import "github.com/golang-jwt/jwt"

type UserDataCheck struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func (r UserDataCheck) FromToken(t *jwt.Token) UserDataCheck {
	claims := t.Claims.(jwt.MapClaims)
	userDataCheck := UserDataCheck{
		Username: claims["username"].(string),
		Password: claims["password"].(string),
		Database: claims["database"].(string),
	}
	return userDataCheck
}
