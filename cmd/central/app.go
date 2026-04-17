package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	handlers "github.com/n4djib/report-engine/internal/api/central"
	"github.com/n4djib/report-engine/internal/api/central/oapi-gen"
	vars "github.com/n4djib/report-engine/internal/vars/central"
	"github.com/n4djib/report-engine/pkg/swagger"
	utilities "github.com/n4djib/report-engine/pkg/utils"
)

type Application struct {
	config vars.ConfigVars
}

func NewApplication(config vars.ConfigVars) *Application {
	return &Application{
		config: config,
	}
}

func (app Application) run() error {
	e := echo.New()
	useCORSMiddleware(e)

	pingHandlers := handlers.CentralHandlers{
		Config: app.config, 
	}  
	// pingHandlers.RegisterHandlers(e.Group("/api"))
	oapi.RegisterHandlers(e, pingHandlers)

	// Register Swagger UI and spec endpoints
	spec, err := oapi.GetSwagger()
	if err != nil {
		return err
	}
	// swagger.RegisterSwagger(e.Group("/"))
	// TODO protect this API
	swagger.RegisterSwagger(e, spec)

	fmt.Println("⇨ Starting App:", app.config.AppName)
	e.HideBanner = app.config.HideBanner
	e.HidePort = app.config.HidePort

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
