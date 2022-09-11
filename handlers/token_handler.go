package handlers

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"net/http"
	"passwords_api/utils"
)

type RequestToken struct {
	TokenValue string `json:"token"`
}

type UserData struct {
	UserName     string `json:"user_name"`
	PasswordHash string `json:"password_hash"`
}

func TokenPost(c echo.Context) error {
	requestToken := new(RequestToken)
	if err := c.Bind(requestToken); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	token, _ := utils.ValidateJWT(requestToken.TokenValue)
	claims := token.Claims.(jwt.MapClaims)
	userData := UserData{
		UserName:     claims["user_name"].(string),
		PasswordHash: claims["password_hash"].(string),
	}
	return c.JSON(http.StatusOK, userData)
}
