CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";


-- pgcrypto: Provides gen_random_uuid()
-- pg_trm: Trigram indexing for fast text search, Enables GIN/GIST indexes on text columns
