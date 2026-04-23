package client

import (
	"context"
	"fmt"

	"report/internal/client/sync"
	"report/pkg/patch"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/zap"
)

type App struct {
	sync     *sync.SyncEngine
	license  *license.SeatActivator
	security *security.Monitor
	log      *zap.Logger
}

func NewApp(cfg sync.Config, log *zap.Logger) (*App, error) {
	engine, err := sync.NewSyncEngine(cfg, log)
	if err != nil {
		return nil, err
	}

	app := &App{sync: engine, log: log}

	// Wire React-facing handlers into the engine
	engine.SetHandlers(
		// onIncoming: tell React a confirmed patch arrived from another user
		func(p patch.Patch) {
			runtime.EventsEmit(app.ctx,
				"patch:incoming",
				map[string]any{
					"reportID":   p.ReportID,
					"sectionKey": p.SectionKey,
					"fieldPath":  p.FieldPath,
					"value":      p.Value,
					"authorID":   p.AuthorID,
				},
			)
		},
		// onError: tell React a patch was rejected — triggers rollback
		func(patchID, reason string) {
			runtime.EventsEmit(app.ctx,
				"patch:error",
				map[string]any{"patchID": patchID, "reason": reason},
			)
		},
	)

	return app, nil
}

// ── Wails IPC bindings ────────────────────────────────────────────────────────

func (a *App) ApplyPatch(p patch.Patch) error {
	ok, count := a.sync.CanQueuePatch()
	if !ok {
		return fmt.Errorf("outbox full (%d entries) — connect to server before continuing", count)
	}
	return a.sync.QueuePatch(p)
}

func (a *App) GetOutboxStatus() enginesync.SyncStatus {
	return a.sync.Status()
}

// ── Wails lifecycle hooks ─────────────────────────────────────────────────────

func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
	if err := a.sync.Start(ctx); err != nil {
		a.log.Error("sync engine failed to start", zap.Error(err))
	}
	a.security.Start(ctx, a.license)
}

func (a *App) OnShutdown(_ context.Context) {
	a.sync.Stop()
}
