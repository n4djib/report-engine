package main

import (
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	handlers "github.com/n4djib/report-engine/internal/api/remote"
	"github.com/n4djib/report-engine/internal/api/remote/oapi-gen"
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

	pingHandlers := handlers.RemoteHandlers{}
	// pingHandlers.RegisterHandlers(e.Group("/api"))
	oapi.RegisterHandlers(e, pingHandlers)

	// Register Swagger UI and spec endpoi nts
	spec, err := oapi.GetSwagger()
	if err != nil {
		return err
	}
	// swagger.RegisterSwagger(e.Group("/"))
	swagger.RegisterSwagger(e, spec)

	return e.Start(":" + strconv.Itoa(int(app.config.AppPort)))
}

func (app Application) openBrowser(url string) {
	if app.config.AppEnv == "production" {
		if err := utilities.OpenURL(url); err != nil {
			log.Fatal("Problem Opening the browser\n", err)
		}
	}
}
