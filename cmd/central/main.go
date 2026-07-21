package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/n4djib/report-engine/pkg/config"
	// "github.com/davecgh/go-spew/spew"
	vars "github.com/n4djib/report-engine/internal/vars/central"
)

func main() {
	// Note: .env.local is not copied to docker
	// .env is copied only in dev
	configFiles := []string{"./cmd/central/env/.env", "./cmd/central/env/.env.local"}
	// we will load the env files in the docker compose file and in the local development environment
	// configFiles := []string{}

	cfg := vars.ConfigVars{}
	err := config.LoadConfigFromFiles(&cfg, configFiles)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	// log the config struct
	// spew.Dump(cfg)

	app, err := NewApplication(cfg)

	// TODO we should open to the frontend app not the api
	// add env var ponting to frontend url and open that instead
	// open browser to APP url
	go app.openBrowser(cfg.AppUrl + ":" + fmt.Sprint(cfg.AppPort) + "/api/ping")

	// root context canceled on Ctrl+C or SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.run(ctx); err != nil {
		slog.Error("Central app failed to run", "error", err)
		os.Exit(1)
	}

	fmt.Println("Central Server Ended!")
}
