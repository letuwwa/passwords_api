package main

import (
	"github.com/labstack/echo/v4"
	"passwords_api/handlers"
)

func main() {
	e := echo.New()
	e.POST("/token/check", handlers.TokenPost)
	e.Logger.Fatal(e.Start("localhost:5050"))
}
