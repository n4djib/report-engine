--
-- The audit trail records every state-changing event immutably.
-- Key design decisions:
--
-- 1. BIGSERIAL (not UUID): Audit entries are sequential.
--    The hash chain depends on ordering — we need a monotonic ID.
--    BIGSERIAL gives us up to 9,223,372,036,854,775,807 entries (plenty).
--
-- 2. prev_hash + signature: Together these implement the hash chain.
--    Each entry signs the previous entry's signature.
--    This makes it impossible to delete an entry without breaking the chain.
--
-- 3. Row-Level Security: The application role (report-engine_app) can INSERT
--    but cannot UPDATE or DELETE. PostgreSQL enforces this.

CREATE TABLE IF NOT EXISTS audit_trail (
    id          BIGSERIAL    PRIMARY KEY,
    timestamp   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    report_id   VARCHAR(255) NOT NULL,
    org_id      UUID         NOT NULL,

    event_type  VARCHAR(50)  NOT NULL
                CHECK (event_type IN (
                    'PATCH', 'STATE_TRANSITION', 'COMMENT',
                    'APPROVAL', 'DENIAL', 'LOCK_ACQUIRE', 'LOCK_RELEASE',
                    'CONFLICT_DETECTED', 'CONFLICT_RESOLVED',
                    'SCHEMA_PUBLISHED', 'USER_CREATED', 'USER_DEACTIVATED',
                    'SEAT_ACTIVATED', 'SEAT_REVOKED',
                    'DEPLOYMENT_LICENSE_VALIDATED', 'DEPLOYMENT_LICENSE_REVOKED'
                )),

    -- actor: who performed this action
    -- {"userID": "...", "role": "reporter", "clientIP": "...", "sessionID": "..."}
    actor       JSONB        NOT NULL,

    -- before_val / after_val: the field value before and after the change
    before_val  JSONB,
    after_val   JSONB,

    -- delta: the RFC 6902 JSON Patch operation that was applied
    -- {"op": "replace", "path": "/safety/dailyHeadcount", "value": 42}
    delta       JSONB,

    nats_seq    BIGINT,       -- NATS JetStream sequence number for correlation

    -- prev_hash: SHA-256 of the previous entry's signature.
    -- For the first entry of a report: "REPORT1"
    -- This creates the hash chain — each entry references the previous one.
    prev_hash   TEXT         NOT NULL,

    -- signature: HMAC-SHA256 of (prev_hash|timestamp|report_id|actor_user_id|delta_json)
    -- Computed using a secret key stored in AUDIT_HMAC_SECRET_FILE.
    -- If any entry is modified, its signature no longer matches.
    signature   TEXT         NOT NULL
);

CREATE INDEX idx_audit_report_ts ON audit_trail (report_id, timestamp DESC);
CREATE INDEX idx_audit_org_ts ON audit_trail (org_id, timestamp DESC);

-- Enable Row-Level Security on this table
ALTER TABLE audit_trail ENABLE ROW LEVEL SECURITY;

-- Allow INSERT for the report-engine_app role (the application role)
-- No UPDATE or DELETE policy is defined, so those operations are blocked.
-- IMPORTANT: The report-engine_app role must exist before this migration runs.
-- Create it with: CREATE ROLE report-engine_app LOGIN PASSWORD '...';
-- Then grant permissions: GRANT INSERT, SELECT ON audit_trail TO report-engine_app;
CREATE POLICY audit_insert_only ON audit_trail
    FOR INSERT TO report_engine_app WITH CHECK (true);

-- Allow SELECT (reading the audit trail) for report-engine_app
CREATE POLICY audit_select ON audit_trail
    FOR SELECT TO report_engine_app USING (true);
