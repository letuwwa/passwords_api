package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"passwords_api/handlers"
)

func main() {
	log.SetFlags(log.Lshortfile)
	e := echo.New()

	// Recover Middleware
	e.Use(middleware.Recover())

	// Logging Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} - ${uri} - ${status} - ${time_rfc3339_nano}\n",
	}))

	// CORS Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "root view\n")
	})
	e.POST("/token", handlers.TokenPost)
	e.POST("/token/check", handlers.TokenCheck)

	e.Logger.Fatal(e.Start("localhost:5050"))
}
