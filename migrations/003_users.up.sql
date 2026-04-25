
CREATE TABLE IF NOT EXISTS  users (
  id    UUID   PRIMARY KEY DEFAULT gen_random_uuid(),
  org_id UUID  NOT NULL REFERENCES organisations(id) ON DELETE CASCADE,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash TEXT  NOT NULL,

  -- Role: enforces that only valid role values can be stored
  -- CHECK constraint is validated by PostgreSQL on every INSERT/UPDATE 
  -- This is a database-level guarantee, not just application-level
  role VARCHAR(50) NOT NULL
       CHECK (role IN('reporter', 'qc', 'superintendent', 'admin')),
  display_name VARCHAR(255) NOT NULL,
  is_active BOOLEAN   NOT NULL DEFAULT true,

  -- seat_token: the current seat license session seat_token
  -- NULL means the user account has never activated a seat on any machine
  -- NOT NULL means have activated and this is thier current session
  seat_token TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_login_at TIMESTAMPTZ 
);


-- Index for lookups by org (list all users in an org)
CREATE INDEX indx_users_org ON users (org_id);

-- Index login (lookup by email is the most common user query)
CREATE INDEX indx_users_email ON users (email);


-- Partial index: only count active users for seat enforcement
-- When we run COUNT(*) WHERE is_active = true, PostgreSQL uses this index
CREATE INDEX idx_users_active ON users (org_id) WHERE is_active = True;
