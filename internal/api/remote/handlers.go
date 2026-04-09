package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/n4djib/report-engine/internal/api/remote/oapi-gen"
)

type RemoteHandlers struct{}

func (h RemoteHandlers) PingPong(ctx echo.Context) error {
	resp := oapi.SharedModelsPingResponse{
		Message: "Pong from remote!",
	}
	return ctx.JSON(http.StatusAccepted, resp)
}
