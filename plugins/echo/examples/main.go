package main

import (
	"net/http"

	souin_echo "github.com/darkweak/souin/plugins/echo"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Use the Souin default configuration
	s := souin_echo.New(souin_echo.DevDefaultConfiguration)
	e.Use(s.Process)

	// Handler
	e.GET("/*", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":80"))
}
