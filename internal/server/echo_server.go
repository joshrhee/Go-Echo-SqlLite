package server

import (
	"identity-coding-test/internal/config"
	"identity-coding-test/internal/handler"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	*echo.Echo

	cfg config.Config
}

func NewEchoServer(cfg config.Config) (*EchoServer, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/healthcheck", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "SERVING",
		})
	})

	e.POST("/users/sign-up", handler.SignUp(cfg.Clock(), cfg.ReadWriteDB()))

	return &EchoServer{Echo: e, cfg: cfg}, nil
}
