package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/n4djib/report-engine/internal/api/central/oapi-gen"
	vars "github.com/n4djib/report-engine/internal/vars/central"
)

type CentralHandlers struct{
	Config vars.ConfigVars
}

func (h CentralHandlers) PingPong(ctx echo.Context) error {
	resp := oapi.SharedModelsPingResponse{
		Message: "Pong from " + h.Config.AppName + "!",
	}
	return ctx.JSON(http.StatusAccepted, resp)
}
