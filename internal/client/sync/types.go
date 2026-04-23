package sync

import (
	"time"

	"report/pkg/patch"
)

// SyncStatus is returned to the ReactUI via Wails IPC
// Tells the user whether they are online, how many patches are pending
// and when the last successful sync happened.
type SyncStatus struct {
	IsOnline      bool      `json:"isOnline"`
	OutboxCount   int       `json:"outboxCount"`
	LastSyncAt    time.Time `json:"lastSyncAt"`
	NATSConnected bool      `json:"natsConnected"`
}

// IncomingHandler is called by the pull loop when a confirmed patch
// arrives from the server (broadcast to all clients in the org)
// The Wails App wires this to invalidate the TanStack Query cache
type IncomingHandler func(p patch.Patch)

// ErrorHandler is called when a patch is rejected by the server
// The Wails App wires this to trigger a TanStack Query rollback
type ErrorHandler func(patchID string, reason string)
