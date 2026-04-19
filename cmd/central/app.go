package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	// "github.com/labstack/echo/v4/middleware"
	handlers "github.com/n4djib/report-engine/internal/api/central"
	"github.com/n4djib/report-engine/internal/api/central/oapi-gen"
	frontendembed "github.com/n4djib/report-engine/internal/embed/central"
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
	// how to server swagger through the frontend app
	swagger.RegisterSwagger(e, spec)

	// register react static pages build from react tanstack router
	frontendembed.RegisterHandlers(e)
	
	// middlewares
	// app.useCORSMiddleware(e)

	fmt.Println("⇨ Starting Central App:", app.config.AppName)
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

// func (app Application) useCORSMiddleware(e *echo.Echo) {
// 	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
// 		// AllowOrigins: []string{"*"},
// 		// AllowOrigins: []string{"http://localhost:3000"},
// 		AllowOrigins: []string{app.config.AllowOrigins},
// 		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
// 		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
// 	}))
// }
