package main

import (
	"fmt"
	"log/slog"
	"os"

	// "github.com/davecgh/go-spew/spew"

	vars "github.com/n4djib/report-engine/internal/vars/remote"
	"github.com/n4djib/report-engine/pkg/config"
)

func main() {
	configFiles := []string{"./cmd/remote/env/.env", "./cmd/remote/env/.env.local"}
	cfg := vars.ConfigVars{}
	err := config.LoadConfigFromFiles(&cfg, configFiles)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	// log the config struct
	// spew.Dump(cfg)

	app := NewApplication(cfg)

	// TODO we should open to the frontend app not the api
	// add env var ponting to frontend url and open that instead
	// open browser to APP url
	go app.openBrowser(cfg.AppUrl + ":" + fmt.Sprint(cfg.AppPort) + "/api/ping")

	if err := app.run(); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}

	fmt.Println("Hello from remote!")
}
