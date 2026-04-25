
CREATE TABLE IF NOT EXISTS change_requests (
  id    UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
  report_id  VARCHAR(255) NOT NULL,
  org_id  UUID NOT NULL,
  requested_by   UUID NOT NULL REFERENCES users(id),
  reason    TEXT NOT NULL,
  status    VARCHAR(50) NOT NULL DEFAULT 'pending'
            CHECK (status IN('pending', 'approved', 'denied')),
  reviewed_by   UUID  REFERENCES users (id),
  review_note   TEXT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  resolved_at TIMESTAMPTZ
);
CREATE INDEX idx_cr_report  ON change_requests (report_id);
CREATE INDEX idx_cr_pending ON change_requests (org_id, status) WHERE status = 'pending';
