package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type PingServer struct{}

func (h PingServer) RegisterHandlers(e *echo.Group) {
	e.GET("/ping", h.pong)
}

func (h PingServer) pong(ctx *echo.Context) error {
	resp := struct {
		Message string
	}{
		Message: "Pong from server!",
	}
	return ctx.JSON(http.StatusAccepted, resp)
}
