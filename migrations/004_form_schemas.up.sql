
CREATE TABLE IF NOT EXISTS form_schemas (
  id    UUID   PRIMARY KEY DEFAULT gen_random_uuid(),
  schema_id  VARCHAR(255) NOT NULL,
  version  INTEGER NOT NULL,
  org_id UUID  NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
  published_at TIMESTAMPTZ, 
  deprecated_at TIMESTAMPTZ, 

  --This stores the entire schema including field types, permissions, validation rules
  -- JSONB allows flexible schema evolution without migrations for every time form change
  migrations_form  JSONB   DEFAULT '{}',
  created_by   UUID REFERENCES users (id),

  -- UNIQUE: prevent two schemas with the same ID and version for the same org
  -- all three values must be unique
  CONSTRAINT uq_schema UNIQUE (schema_id, version, org_id)
);

-- Composite index matching the unique constraint - used for lookups
CREATE INDEX idx_form_schemas_lookup ON form_schemas (schema_id, version, org_id);



-- GIN = Generlised Inverted Index, it indexes every key and value inside
-- the JSONB document, enabling fast queries like:
-- WHERE sections @> '{"sections": [{"sectionKey": "saftey"}]}'
-- Without this index, PostgreSQL would scan every row
CREATE INDEX idx_form_schemas_gin ON form_schemas USING GIN (migrations_form);





