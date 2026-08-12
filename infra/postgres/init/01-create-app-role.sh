#!/bin/sh
# ─────────────────────────────────────────────────────────────────────────────
# Creates the least-privilege runtime role that backend-core connects as.
#
# Runs automatically, once, when Postgres initialises an empty data volume.
# To re-run it you must destroy the volume: `make reset`.
#
# WHY A SEPARATE ROLE:
# Postgres silently bypasses Row-Level Security for superusers and for a
# table's owner. If the application connected as ${POSTGRES_USER} (the owner
# that runs migrations), every RLS policy would still be defined but would
# enforce nothing — the most dangerous outcome, because the schema looks
# protected while cross-clinic reads succeed. So:
#
#   ${POSTGRES_USER}  → owns the schema, runs migrations, bypasses RLS
#   ${APP_DB_USER}    → runtime role, owns nothing, RLS is enforced
#
# The app role deliberately gets DML (select/insert/update/delete) but no DDL:
# any table it created would be owned by it, and therefore exempt from RLS.
# ─────────────────────────────────────────────────────────────────────────────
set -eu

echo "init: creating runtime role '${APP_DB_USER}' in database '${POSTGRES_DB}'"

psql -v ON_ERROR_STOP=1 \
     --username "${POSTGRES_USER}" \
     --dbname "${POSTGRES_DB}" \
     -v app_user="${APP_DB_USER}" \
     -v app_password="${APP_DB_PASSWORD}" \
     -v owner="${POSTGRES_USER}" \
     -v db_name="${POSTGRES_DB}" <<'EOSQL'
\set quoted_app_user :"app_user"
\set quoted_app_password :'app_password'

-- CREATE ROLE has no IF NOT EXISTS, so guard it.
SELECT format(
  'CREATE ROLE %I LOGIN PASSWORD %L NOSUPERUSER NOCREATEDB NOCREATEROLE NOBYPASSRLS',
  :'app_user', :'app_password'
) AS stmt
WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = :'app_user')
\gexec

-- Idempotent: keep the password in sync with the environment on re-init.
SELECT format('ALTER ROLE %I PASSWORD %L', :'app_user', :'app_password') \gexec

GRANT CONNECT ON DATABASE :"db_name" TO :"app_user";
GRANT USAGE ON SCHEMA public TO :"app_user";

-- Objects that already exist (none on a first run, but keeps this re-runnable).
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO :"app_user";
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO :"app_user";

-- Objects that migrations will create later, as the owner.
ALTER DEFAULT PRIVILEGES FOR ROLE :"owner" IN SCHEMA public
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO :"app_user";
ALTER DEFAULT PRIVILEGES FOR ROLE :"owner" IN SCHEMA public
  GRANT USAGE, SELECT ON SEQUENCES TO :"app_user";

-- Belt and braces: nobody gets to create objects in public by default.
REVOKE CREATE ON SCHEMA public FROM PUBLIC;
EOSQL

echo "init: runtime role '${APP_DB_USER}' ready (NOSUPERUSER, NOBYPASSRLS, no DDL)"
