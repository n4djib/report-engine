
-- Organisation is one report-engine installation
-- Slug field is what we use in NATS subject names:
-- reports.{slug}.{reportID}.patch
-- It must be immutable after creation - changing it will break all NATS subjects.

CREATE TABLE IF NOT EXISTS organisations (
  id   UUID   PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(100) NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  settings JSONB NOT NULL DEFAULT '{}'

);

-- Index on slug for fast lookups by NATS subject prefix
CREATE INDEX idx_organisations_slug ON organisations(slug);

