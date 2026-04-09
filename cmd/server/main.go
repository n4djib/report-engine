package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/n4djib/report-engine/pkg/config"
	// "github.com/davecgh/go-spew/spew"
)

func main() {
	configFiles := []string{"./cmd/server/.env", "./cmd/server/.env.local"}

	cfg := ConfigVars{}
	err := config.LoadConfigFromFiles(&cfg, configFiles)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	// log the config struct
	// spew.Dump(cfg)

	app := NewApplication(cfg)

	// TODO we should open to the frontend app not the api
	// open browser to APP url
	go app.openBrowser(cfg.AppUrl + ":" + fmt.Sprint(cfg.AppPort) + "/api/ping")

	if err := app.run(); err != nil {
		slog.Error("server failed to start", "error", err)
		os.Exit(1)
	}

	fmt.Println("Server Ended!")
}
