-- Shared building blocks used by every later migration.
--
-- Run as the OWNER role (physio), so everything created here is owned by the
-- owner and the runtime role (physio_app) only gets the privileges granted by
-- ALTER DEFAULT PRIVILEGES in infra/postgres/init/01-create-app-role.sh.

-- +goose Up

-- Case-insensitive text, used for email so that Ada@clinic.ng and
-- ada@clinic.ng cannot both be registered.
CREATE EXTENSION IF NOT EXISTS citext;

-- Keeps updated_at honest without the application having to remember.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END;
$$;
-- +goose StatementEnd

COMMENT ON FUNCTION set_updated_at() IS
  'Trigger function: stamps NEW.updated_at with now() on UPDATE.';

-- ─────────────────────────────────────────────────────────────────────────────
-- Tenant context for Row-Level Security.
--
-- backend-core sets the current clinic per transaction:
--     SELECT set_config('app.current_clinic_id', $1, true);
-- The third argument (is_local = true) scopes the setting to the surrounding
-- transaction, which is essential with a connection pool: a plain SET would
-- persist on the pooled connection and leak the previous request's tenant into
-- the next one.
--
-- Every RLS policy compares clinic_id against this function. When the setting
-- is absent it returns NULL, and `clinic_id = NULL` matches no rows — so the
-- policies fail CLOSED. Forgetting to set the context denies access rather
-- than exposing every clinic.
-- ─────────────────────────────────────────────────────────────────────────────
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION current_clinic_id()
RETURNS uuid
LANGUAGE sql
STABLE
AS $$
  -- missing_ok = true: return NULL instead of raising when unset.
  SELECT NULLIF(current_setting('app.current_clinic_id', true), '')::uuid;
$$;
-- +goose StatementEnd

COMMENT ON FUNCTION current_clinic_id() IS
  'Clinic id for the current transaction, from the app.current_clinic_id setting. NULL when unset, which makes RLS policies fail closed.';

-- +goose Down

DROP FUNCTION IF EXISTS current_clinic_id();
DROP FUNCTION IF EXISTS set_updated_at();
-- citext is intentionally left installed: dropping an extension other objects
-- may depend on is riskier than leaving it in place.
