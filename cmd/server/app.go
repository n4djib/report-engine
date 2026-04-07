package main

import (
	"context"
	"log"
	"strconv"

	"github.com/labstack/echo/v5"
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

	pingHandlers := handlers.PingServer{}
	pingHandlers.RegisterHandlers(e.Group("/api"))

	// Register Swagger UI and spec endpoints
	spec, err := oapi.GetSwagger()
	if err != nil {
		return err
	}
	// swagger.RegisterSwagger(e.Group("/"))
	swagger.RegisterSwagger(e, spec)

	// Configure the http server
	sc := echo.StartConfig{
		Address:    ":" + strconv.Itoa(int(app.config.AppPort)),
		HideBanner: app.config.HideBanner, // This replaces e.HideBanner = true
		HidePort:   app.config.HidePort,   // Optional: also hides the "listening on..." message
	}

	// // we need to handle shutdown gracefully, we need to listen to the interrupt signal and shutdown the server gracefully
	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // start shutdown process on signal
	// defer cancel()
	ctx := context.Background()

	return sc.Start(ctx, e)
}

func (app Application) openBrowser(url string) {
	if app.config.AppEnv == "production" {
		if err := utilities.OpenURL(url); err != nil {
			log.Fatal("Problem Opening the browser\n", err)
		}
	}
}
