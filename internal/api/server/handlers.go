package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/n4djib/report-engine/internal/api/server/oapi-gen"
)

type ServerHandlers struct{}

func (h ServerHandlers) PingPong(ctx echo.Context) error {
	resp := oapi.SharedModelsPingResponse{
		Message: "Pong from server!",
	}
	return ctx.JSON(http.StatusAccepted, resp)
}
