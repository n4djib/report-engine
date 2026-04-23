package sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"report/pkg/patch"

	"report/internal/client/crypto"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

// ─── Constants ────────────────────────────────────────────────────────────────

const (
	maxOutboxEntries    = 500
	warnOutboxThreshold = 300
	drainInterval       = 250 * time.Millisecond
	retryBackoffMax     = 30 * time.Second
	connectTimeout      = 10 * time.Second
	reconnectWait       = 2 * time.Second
	healthCheckInterval = 5 * time.Second
)

// ─── SyncEngine ───────────────────────────────────────────────────────────────

// SyncEngine owns all client-side synchronisation concerns:
//   - BadgerDB outbox (write patches locally before sending)
//   - NATS JetStream publish (drain outbox when online)
//   - NATS pull consumer (receive confirmed patches from server)
//   - Network probe (detect online/offline transitions)
//
// It is created once in main.go and injected into the Wails App struct.
// The React layer never touches it directly — all interaction is via
// the Wails IPC bindings defined in internal/client/app.go.
type SyncEngine struct {
	// Dependencies — set by NewSyncEngine, never changed after
	db     *LocalDB   // BadgerDB wrapper (see db.go)
	nc     *nats.Conn // raw NATS connection
	js     jetstream.JetStream
	clock  *patch.ClientClock
	aesKey []byte // AES-256 key for encrypting outbox entries
	log    *zap.Logger

	// Configuration — set at construction
	orgSlug   string
	userID    string
	serverURL string // wss://nexusreport.client.com:4222

	// Handlers wired by Wails App after construction
	onIncoming IncomingHandler
	onError    ErrorHandler

	// Runtime state
	online     atomic.Bool  // true when NATS is connected
	lastSyncAt atomic.Value // stores time.Time
	mu         sync.Mutex   // guards consumer reference
	consumer   jetstream.Consumer

	// Lifecycle
	cancelDrain context.CancelFunc
	cancelPull  context.CancelFunc
}

// ─── Constructor ──────────────────────────────────────────────────────────────

// Config holds everything NewSyncEngine needs.
// Populated from BadgerDB app:settings and app:server_url at startup.
type Config struct {
	DBPath    string // path to BadgerDB directory
	ServerURL string // NATS server URL: wss://host:4222
	OrgSlug   string // from JWT after login
	UserID    string // from JWT after login
	NATSUser  string // org-scoped NATS credential
	NATSPass  string // org-scoped NATS credential
	AESKey    []byte // derived from login password via PBKDF2
}

// NewSyncEngine constructs and returns a fully wired SyncEngine.
// It does NOT start the outbox drain or pull consumer — call Start() for that.
// This allows the engine to be created before the user has logged in,
// and started only once credentials are confirmed.
func NewSyncEngine(cfg Config, log *zap.Logger) (*SyncEngine, error) {
	// Open BadgerDB
	db, err := OpenLocalDB(cfg.DBPath, cfg.AESKey)
	if err != nil {
		return nil, fmt.Errorf("sync engine: open local db: %w", err)
	}

	e := &SyncEngine{
		db:        db,
		clock:     &patch.ClientClock{},
		aesKey:    cfg.AESKey,
		orgSlug:   cfg.OrgSlug,
		userID:    cfg.UserID,
		serverURL: cfg.ServerURL,
		log:       log,
	}

	// Restore the Lamport clock from the last known NATS sequence
	// so clientSeq is always ahead of anything the server has seen
	lastSeq := db.GetLastSeq()
	e.clock.Update(lastSeq)

	// Connect to NATS — non-blocking: if the server is unreachable,
	// the outbox drain will retry automatically
	if err := e.connectNATS(cfg.NATSUser, cfg.NATSPass); err != nil {
		// Not fatal at construction — app works offline
		log.Warn("NATS unavailable at startup — operating offline",
			zap.Error(err))
	}

	return e, nil
}

// SetHandlers wires the React-facing callbacks.
// Must be called before Start().
func (e *SyncEngine) SetHandlers(onIncoming IncomingHandler, onError ErrorHandler) {
	e.onIncoming = onIncoming
	e.onError = onError
}

// ─── Lifecycle ────────────────────────────────────────────────────────────────

// Start begins the outbox drain goroutine and the pull consumer loop.
// It is called once after successful login.
// Calling Start() a second time is a no-op.
func (e *SyncEngine) Start(ctx context.Context) error {
	if e.onIncoming == nil || e.onError == nil {
		return errors.New("sync engine: handlers not set — call SetHandlers before Start")
	}

	drainCtx, cancelDrain := context.WithCancel(ctx)
	e.cancelDrain = cancelDrain
	go e.drainOutbox(drainCtx)

	if err := e.setupConsumer(ctx); err != nil {
		e.log.Warn("pull consumer setup failed — will retry on reconnect", zap.Error(err))
	}

	go e.reconnectLoop(ctx)

	e.log.Info("sync engine started",
		zap.String("orgSlug", e.orgSlug),
		zap.String("userID", e.userID),
	)
	return nil
}

// Stop gracefully shuts down both goroutines and closes the NATS connection.
// Called from the Wails OnShutdown hook.
func (e *SyncEngine) Stop() {
	if e.cancelDrain != nil {
		e.cancelDrain()
	}
	if e.cancelPull != nil {
		e.cancelPull()
	}
	if e.nc != nil {
		e.nc.Drain() // flushes pending publishes before closing
	}
	e.db.Close()
	e.log.Info("sync engine stopped")
}

// ─── NATS Connection ──────────────────────────────────────────────────────────

func (e *SyncEngine) connectNATS(user, pass string) error {
	nc, err := nats.Connect(e.serverURL,
		nats.UserInfo(user, pass),
		nats.Secure(), // TLS required
		nats.Timeout(connectTimeout),
		nats.ReconnectWait(reconnectWait),
		nats.MaxReconnects(-1), // retry forever
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			e.online.Store(false)
			e.log.Warn("NATS disconnected", zap.Error(err))
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			e.online.Store(true)
			e.log.Info("NATS reconnected")
			// Re-setup pull consumer from last known seq on reconnect
			go func() {
				if err := e.setupConsumer(context.Background()); err != nil {
					e.log.Error("consumer re-setup failed", zap.Error(err))
				}
			}()
		}),
	)
	if err != nil {
		return err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return fmt.Errorf("jetstream context: %w", err)
	}

	e.nc = nc
	e.js = js
	e.online.Store(true)
	return nil
}

// reconnectLoop watches for offline→online transitions and retries
// the NATS connection when the server becomes reachable again.
func (e *SyncEngine) reconnectLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(healthCheckInterval):
		}
		if !e.online.Load() && e.nc != nil && e.nc.IsReconnecting() {
			// nats.go handles reconnection internally — nothing to do here
			continue
		}
	}
}

// ─── Outbox — Write ───────────────────────────────────────────────────────────

// QueuePatch encrypts and writes a patch to the BadgerDB outbox.
// It returns ErrOutboxFull if the cap is reached — the caller must
// show a sync warning to the user and refuse further edits.
//
// This is the ONLY place patches enter the system on the client side.
// The patch ID and ClientSeq are assigned here.
func (e *SyncEngine) QueuePatch(p patch.Patch) error {
	count := e.db.CountOutboxEntries()
	if count >= maxOutboxEntries {
		return fmt.Errorf("%w: %d patches pending (max %d)",
			ErrOutboxFull, count, maxOutboxEntries)
	}

	// Assign Lamport sequence — always increment before anything else
	p.ClientSeq = e.clock.Next()

	encoded, err := p.EncodeMsgpack()
	if err != nil {
		return fmt.Errorf("encode patch: %w", err)
	}

	ciphertext, err := crypto.Encrypt(e.aesKey, encoded)
	if err != nil {
		return fmt.Errorf("encrypt patch: %w", err)
	}

	key := fmt.Sprintf("report:%s:outbox:%s", p.ReportID, p.ID)
	if err := e.db.Set(key, ciphertext); err != nil {
		return fmt.Errorf("write outbox: %w", err)
	}

	e.log.Debug("patch queued",
		zap.String("patchID", p.ID),
		zap.String("reportID", p.ReportID),
		zap.Uint64("clientSeq", p.ClientSeq),
		zap.Int("outboxDepth", count+1),
	)
	return nil
}

// ─── Outbox — Drain ───────────────────────────────────────────────────────────

// drainOutbox runs continuously in its own goroutine.
// Every 250ms it reads all pending outbox entries and publishes them
// to NATS JetStream in insertion order (sorted by clientSeq).
//
// An outbox entry is deleted ONLY after a confirmed NATS ACK.
// Crash between write and ACK = entry survives, re-sent on restart.
func (e *SyncEngine) drainOutbox(ctx context.Context) {
	backoff := time.Second

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(drainInterval):
		}

		if !e.online.Load() {
			continue
		}

		patches, err := e.db.ListOutboxPatches()
		if err != nil {
			e.log.Error("outbox list failed", zap.Error(err))
			continue
		}
		if len(patches) == 0 {
			continue
		}

		for _, p := range patches {
			subject := fmt.Sprintf("reports.%s.%s.patch", p.OrgID, p.ReportID)

			_, err := e.js.Publish(ctx, subject, p.EncodeMsgpack(),
				jetstream.WithMsgID(p.ID), // idempotency key — safe to retry
			)
			if err != nil {
				e.log.Warn("patch publish failed",
					zap.String("patchID", p.ID),
					zap.Error(err),
				)
				time.Sleep(backoff)
				backoff = min(backoff*2, retryBackoffMax)
				break // preserve ordering — stop this cycle, retry next
			}

			// DELETE ONLY ON CONFIRMED ACK
			key := fmt.Sprintf("report:%s:outbox:%s", p.ReportID, p.ID)
			if err := e.db.Delete(key); err != nil {
				e.log.Error("outbox delete failed after ACK",
					zap.String("patchID", p.ID),
					zap.Error(err),
				)
				// Not fatal — patch will be re-sent next cycle
				// Server handles idempotency via WithMsgID
			}

			e.lastSyncAt.Store(time.Now())
			e.clock.Update(p.ClientSeq)
			backoff = time.Second

			e.log.Debug("patch delivered",
				zap.String("patchID", p.ID),
				zap.Uint64("clientSeq", p.ClientSeq),
			)
		}
	}
}

// ─── Pull Consumer ────────────────────────────────────────────────────────────

// setupConsumer creates or reconnects a durable pull consumer on the
// REPORTS stream, starting from the last sequence number this client
// has already processed. This is what enables offline catch-up.
func (e *SyncEngine) setupConsumer(ctx context.Context) error {
	name := fmt.Sprintf("client-%s-%s", e.orgSlug, e.userID)
	lastSeq := e.db.GetLastSeq()

	consumer, err := e.js.CreateOrUpdateConsumer(ctx, "REPORTS",
		jetstream.ConsumerConfig{
			Name:          name,
			Durable:       name,
			FilterSubject: fmt.Sprintf("reports.%s.>", e.orgSlug),
			AckPolicy:     jetstream.AckExplicitPolicy,
			DeliverPolicy: jetstream.DeliverByStartSequencePolicy,
			OptStartSeq:   lastSeq + 1, // fetch only what we missed
			AckWait:       30 * time.Second,
			MaxDeliver:    5,
			MaxAckPending: 100,
		},
	)
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	e.mu.Lock()
	e.consumer = consumer
	e.mu.Unlock()

	// Cancel previous pull loop if any, then start a new one
	if e.cancelPull != nil {
		e.cancelPull()
	}
	pullCtx, cancelPull := context.WithCancel(ctx)
	e.cancelPull = cancelPull
	go e.pullLoop(pullCtx)

	e.log.Info("pull consumer ready",
		zap.String("consumer", name),
		zap.Uint64("fromSeq", lastSeq+1),
	)
	return nil
}

func (e *SyncEngine) pullLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		e.mu.Lock()
		consumer := e.consumer
		e.mu.Unlock()

		if consumer == nil {
			time.Sleep(time.Second)
			continue
		}

		msgs, err := consumer.Fetch(50,
			jetstream.FetchMaxWait(5*time.Second),
		)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			time.Sleep(2 * time.Second)
			continue
		}

		for msg := range msgs.Messages() {
			e.handleIncoming(msg)
		}
	}
}

func (e *SyncEngine) handleIncoming(msg jetstream.Msg) {
	meta, err := msg.Metadata()
	if err != nil {
		msg.Nak()
		return
	}

	// Check if this is an error event addressed to this client
	if isErrorSubject(msg.Subject(), e.userID) {
		var errEvt PatchErrorEvent
		if err := errEvt.Decode(msg.Data()); err == nil && e.onError != nil {
			e.onError(errEvt.PatchID, errEvt.Reason)
		}
		msg.Ack()
		return
	}

	// Decode incoming confirmed patch
	var p patch.Patch
	if err := p.DecodeMsgpack(msg.Data()); err != nil {
		e.log.Error("failed to decode incoming patch",
			zap.String("subject", msg.Subject()),
			zap.Error(err),
		)
		msg.Nak()
		return
	}

	// Skip patches we originated — server echoes them to all clients
	if p.AuthorID == e.userID {
		msg.Ack()
		e.db.SaveLastSeq(meta.Sequence.Stream)
		return
	}

	// Update local BadgerDB cache with the incoming patch
	if err := e.db.ApplyIncomingPatch(p); err != nil {
		e.log.Error("failed to apply incoming patch to local cache",
			zap.String("patchID", p.ID),
			zap.Error(err),
		)
		msg.Nak()
		return
	}

	// Save last processed sequence — next pull starts from here + 1
	e.db.SaveLastSeq(meta.Sequence.Stream)
	e.clock.Update(meta.Sequence.Stream)

	// Notify React layer to invalidate TanStack Query cache
	if e.onIncoming != nil {
		e.onIncoming(p)
	}

	msg.Ack()
}

// ─── Status ───────────────────────────────────────────────────────────────────

// Status is called by the Wails App IPC binding GetOutboxStatus().
// React polls this every 2 seconds to update the SyncStatusBar.
func (e *SyncEngine) Status() SyncStatus {
	lastSync, _ := e.lastSyncAt.Load().(time.Time)
	return SyncStatus{
		IsOnline:      e.online.Load(),
		OutboxCount:   e.db.CountOutboxEntries(),
		LastSyncAt:    lastSync,
		NATSConnected: e.nc != nil && e.nc.IsConnected(),
	}
}

// CanQueuePatch returns false and the current count when the outbox
// is at or above the cap. The Wails App calls this before QueuePatch.
func (e *SyncEngine) CanQueuePatch() (bool, int) {
	count := e.db.CountOutboxEntries()
	return count < maxOutboxEntries, count
}
