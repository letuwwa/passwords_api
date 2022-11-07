package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"passwords_api/handlers"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} - ${uri} - ${status} - ${time_rfc3339_nano}\n",
	}))

	e.POST("/token", handlers.TokenPost)
	e.POST("/token/check", handlers.TokenCheck)
	e.Logger.Fatal(e.Start("localhost:5050"))
}
