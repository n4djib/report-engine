-- create report_engine_app for autdit
CREATE ROLE report_engine_app WITH LOGIN PASSWORD 'audit_123';

-- Create the lookup user
CREATE ROLE pgbouncer_auth WITH LOGIN PASSWORD 'auth_pass_123';

-- Create the helper function
CREATE OR REPLACE FUNCTION public.lookup_auth(uname TEXT)
RETURNS TABLE (usename name, passwd text) AS $$
  SELECT usename, passwd FROM pg_shadow WHERE usename=$1;
$$ LANGUAGE sql SECURITY DEFINER;

-- Secure it
REVOKE ALL ON FUNCTION public.lookup_auth(TEXT) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION public.lookup_auth(TEXT) TO pgbouncer_auth;

