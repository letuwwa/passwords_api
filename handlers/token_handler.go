package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type RequestToken struct {
	TokenValue string `json:"token"`
}

func TokenPost(c echo.Context) error {
	requestToken := new(RequestToken)
	if err := c.Bind(requestToken); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusAccepted, requestToken.TokenValue)
}
