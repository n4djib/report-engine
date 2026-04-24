
CREATE TABLE NOT EXISTS reports (
  id   UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
  report_id  VARCHAR(255) NOT NULL,
  org_id    UUID   NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
  schema_id  VARCHAR(255) NOT NULL,

  -- schema_version is FROZEN when the report is submitted
  -- this report will alwayse render using the version it was submitted with 
  schema_version INTEGER  NOT NULL,
  project_id VARCHAR(255) NOT NULL,


  -- status: the current state machine state
  -- CHECK constraint enforces only valid values
  -- Changed by state machine transitions only - never by direct SQL UPDATE
  status   VARCHAR(50) NOT NULL DEFAULT 'draft'
           CHECK (status IN(
              'draft', 'in_review', 'controlled', 'change_request', 'archived'
          )),

  created_by UUID    NOT NULL REFERENCES users(id),
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  submitted_at  TIMESTAMPTZ,   -- set when transitioning draft -> in_review
  controlled_at  TIMESTAMPTZ,   -- set when transitioning in_review -> controlled


  -- data: all filled-in field values as JSONB
  -- Structure mirrors the form schema sections
  data    JSONB   NOT NULL DEFAULT '{}',


  -- locks: metadata about which sections are currently advisory-locked
  -- example: {"safety": {"userID": "...", "acquiredAt": "...", "expiresAt": "..."}}
  -- NOTE: The actual lock is a PostgreSQL advisory lock
  -- this field just stores metadata for display in the UI
  locks    JSONB   NOT NULL DEFAULT '{}',

  -- conflicts: unresolved concurrent edit forks.
  -- {"safety.incidents.0.severity": {"authorA": {...}, "authorB": {...}}}
  -- Non-empty means Submit() is blocked.
  conflicts      JSONB        NOT NULL DEFAULT '{}',

  -- change_requests: history of change request records.
  -- An array of {requestedBy, reason, status, reviewedBy, note, timestamps}
  change_requests JSONB       NOT NULL DEFAULT '[]',

  -- Foreign key to form_schemas ensures we cannot create a report
  -- referencing a schema that does not exist
  CONSTRAINT fk_schema FOREIGN KEY (schema_id, schema_version, org_id)
      REFERENCES form_schemas (schema_id, version, org_id)
);

CREATE UNIQUE INDEX idx_reports_uq ON reports (report_id, org_id);
CREATE INDEX idx_reports_org_status ON reports (org_id, status);
CREATE INDEX idx_reports_project ON reports (project_id);
CREATE INDEX idx_reports_created_at ON reports (created_at DESC);
-- GIN index on data JSONB for fast field queries
CREATE INDEX idx_reports_data ON reports USING GIN (data);

-- Trigger: automatically update updated_at on every UPDATE
-- Without this, you would have to remember to set updated_at in every
-- Go function that modifies a report.
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER reports_updated_at
    BEFORE UPDATE ON reports
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

