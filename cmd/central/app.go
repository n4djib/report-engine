package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"

	handlers "github.com/n4djib/report-engine/internal/api/central"
	"github.com/n4djib/report-engine/internal/api/central/oapi-gen"
	frontendembed "github.com/n4djib/report-engine/internal/embed/central"
	vars "github.com/n4djib/report-engine/internal/vars/central"
	"github.com/n4djib/report-engine/pkg/swagger"
	utilities "github.com/n4djib/report-engine/pkg/utils"
)

type Application struct {
	config vars.ConfigVars
	logger *slog.Logger
	// db     *sql.DB
	nats   *nats.Conn
	echo   *echo.Echo
}

func NewApplication(config vars.ConfigVars) (*Application, error) {
	app := &Application{
		config: config,
	}

	// init logger
	app.initLogger()

	// // init database
	// if err := app.initDB(); err != nil {
	// 	return nil, err
	// }

	// init NATS
	if err := app.initNATS(); err != nil {
		return nil, err
	}

	// init HTTP server
	if err := app.initHTTP(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *Application) initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	app.logger = slog.New(handler)

	app.logger = app.logger.With(
		"service", "central",
		"env", app.config.AppEnv,
	)
}

func (app *Application) initHTTP() error {
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

	fmt.Println("⇨ Starting", app.config.AppName,"App:", app.config.AppName)
	e.HideBanner = app.config.HideBanner
	e.HidePort = app.config.HidePort

	app.echo = e
	return nil
}

func (app *Application) initNATS() error {
	// nc, err := nats.Connect(app.config.NatsURL)
	// if err != nil {
	// 	return err
	// }

	// app.nats = nc
	return nil
}

func (app *Application) run(ctx context.Context) error {
	errCh := make(chan error, 1)

	// Start HTTP server
	go func() {
		app.logger.Info("starting HTTP server", "port", app.config.AppPort)

		if err := app.echo.Start(":" + strconv.Itoa(int(app.config.AppPort))); err != nil &&
			err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Start other background workers if needed
	// e.g. NATS consumers already running

	select {
	case <-ctx.Done():
		app.logger.Info("shutdown signal received")
	case err := <-errCh:
		return err
	}

	// Begin graceful shutdown
	return app.Shutdown(context.Background())
}

func (app *Application) Shutdown(parent context.Context) error {
	app.logger.Info("starting graceful shutdown")

	// global timeout
	ctx, cancel := context.WithTimeout(parent, 10*time.Second)
	defer cancel()

	// --- HTTP ---
	if app.echo != nil {
		app.logger.Info("shutting down HTTP server")
		if err := app.echo.Shutdown(ctx); err != nil {
			app.logger.Error("HTTP shutdown failed", "error", err)
		}
	}

	// --- NATS ---
	if app.nats != nil {
		app.logger.Info("draining NATS connection")
		if err := app.nats.Drain(); err != nil {
			app.logger.Error("NATS drain failed", "error", err)
		}
		app.nats.Close()
	}

	// // --- DB ---
	// if app.db != nil {
	// 	app.logger.Info("closing database")
	// 	if err := app.db.Close(); err != nil {
	// 		app.logger.Error("DB close failed", "error", err)
	// 	}
	// }

	app.logger.Info("shutdown complete")
	return nil
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
