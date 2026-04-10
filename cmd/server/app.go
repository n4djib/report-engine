package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	handlers "github.com/n4djib/report-engine/internal/api/server"
	"github.com/n4djib/report-engine/internal/api/server/oapi-gen"
	"github.com/n4djib/report-engine/pkg/swagger"
	utilities "github.com/n4djib/report-engine/pkg/utils"
)

type Application struct {
	config ConfigVars
}

func NewApplication(config ConfigVars) *Application {
	return &Application{
		config: config,
	}
}

func (app Application) run() error {
	e := echo.New()

	useCORSMiddleware(e)

	pingHandlers := handlers.ServerHandlers{}
	// pingHandlers.RegisterHandlers(e.Group("/api"))
	oapi.RegisterHandlers(e, pingHandlers)

	// Register Swagger UI and spec endpoints
	spec, err := oapi.GetSwagger()
	if err != nil {
		return err
	}
	// swagger.RegisterSwagger(e.Group("/"))
	swagger.RegisterSwagger(e, spec)

	// Hide Banner
	e.HideBanner = true

	return e.Start(":" + strconv.Itoa(int(app.config.AppPort)))
}

func (app Application) openBrowser(url string) {
	if app.config.AppEnv == "production" {
		if err := utilities.OpenURL(url); err != nil {
			log.Fatal("Problem Opening the browser\n", err)
		}
	}
}

func useCORSMiddleware(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		// AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
}
