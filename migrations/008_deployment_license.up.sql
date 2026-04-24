
-- Only one row. Inserted at first startup with the signed license blob
-- Validated at startup and on daily heartbeat aginst licensing.report-engine.io (our licensing server)

CREATE TABLE IF NOT EXISTS deployment_license (
  id      UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
  license_id   VARCHAR(255)  NOT NULL UNIQUE,
  deployment_name   VARCHAR(255) NOT NULL,
  max_seats   INTEGER    NOT NULL,
  licensed_until TIMESTAMPTZ,
  last_validated TIMESTAMPTZ,
  is_valid BOOLEAN  NOT NULL DEFAULT false,
  raw_license    TEXT   NOT NULL -- full Ed25519-signed license blob from our licensing server
);
