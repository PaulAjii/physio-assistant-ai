-- Row-Level Security: the point where multi-tenancy actually turns on.
--
-- ENFORCEMENT MODEL (read this first):
--   * RLS is skipped for superusers, the table OWNER, and BYPASSRLS roles.
--   * The app connects as physio_app — NOSUPERUSER, NOBYPASSRLS, NOT the owner —
--     so every policy below is enforced against it.
--   * The owner (physio) bypasses RLS. That is deliberate: migrations run as the
--     owner, and the SECURITY DEFINER functions at the bottom run with the
--     owner's rights so they can do the three lookups that MUST happen before a
--     tenant context exists (login, refresh, invite acceptance).
--
--   We therefore use ENABLE, not FORCE. FORCE would subject the owner to RLS as
--   well and break those bootstrap functions — the owner would see zero rows
--   with no context set, so no one could ever log in.
--
-- FAIL CLOSED: policies compare against current_clinic_id() (00001), which is
-- NULL when the app.current_clinic_id setting is unset. `clinic_id = NULL` is
-- NULL, so no rows match — forgetting to set the tenant context denies access
-- rather than leaking every clinic.

-- +goose Up

-- ── Tenant-isolation policies ────────────────────────────────────────────────
-- Every tenant table: enable RLS + one FOR ALL policy. USING gates what rows are
-- visible/updatable/deletable; WITH CHECK gates what rows may be written, so a
-- clinic cannot insert or move a row into someone else's tenant.

-- clinics is keyed on its own id, not a clinic_id column.
ALTER TABLE clinics ENABLE ROW LEVEL SECURITY;
CREATE POLICY clinics_tenant_isolation ON clinics
  USING (id = current_clinic_id())
  WITH CHECK (id = current_clinic_id());

ALTER TABLE users ENABLE ROW LEVEL SECURITY;
CREATE POLICY users_tenant_isolation ON users
  USING (clinic_id = current_clinic_id())
  WITH CHECK (clinic_id = current_clinic_id());

ALTER TABLE refresh_tokens ENABLE ROW LEVEL SECURITY;
CREATE POLICY refresh_tokens_tenant_isolation ON refresh_tokens
  USING (clinic_id = current_clinic_id())
  WITH CHECK (clinic_id = current_clinic_id());

ALTER TABLE invitations ENABLE ROW LEVEL SECURITY;
CREATE POLICY invitations_tenant_isolation ON invitations
  USING (clinic_id = current_clinic_id())
  WITH CHECK (clinic_id = current_clinic_id());

ALTER TABLE objective_templates ENABLE ROW LEVEL SECURITY;
CREATE POLICY objective_templates_tenant_isolation ON objective_templates
  USING (clinic_id = current_clinic_id())
  WITH CHECK (clinic_id = current_clinic_id());

ALTER TABLE consultation_sessions ENABLE ROW LEVEL SECURITY;
CREATE POLICY consultation_sessions_tenant_isolation ON consultation_sessions
  USING (clinic_id = current_clinic_id())
  WITH CHECK (clinic_id = current_clinic_id());

ALTER TABLE media_assets ENABLE ROW LEVEL SECURITY;
CREATE POLICY media_assets_tenant_isolation ON media_assets
  USING (clinic_id = current_clinic_id())
  WITH CHECK (clinic_id = current_clinic_id());

ALTER TABLE assessments ENABLE ROW LEVEL SECURITY;
CREATE POLICY assessments_tenant_isolation ON assessments
  USING (clinic_id = current_clinic_id())
  WITH CHECK (clinic_id = current_clinic_id());

ALTER TABLE transcript_segments ENABLE ROW LEVEL SECURITY;
CREATE POLICY transcript_segments_tenant_isolation ON transcript_segments
  USING (clinic_id = current_clinic_id())
  WITH CHECK (clinic_id = current_clinic_id());

-- ── Auth bootstrap (the only RLS-exempt code path) ───────────────────────────
-- These run SECURITY DEFINER (as the owner, who bypasses RLS) because they
-- answer questions asked BEFORE a tenant context can exist. Each returns only
-- the columns its caller needs — no full rows past the RLS boundary.
--
-- SET search_path pins resolution to trusted schemas so a caller cannot shadow
-- `users`/`clinics` with an object in an earlier schema (pg_temp is left out, so
-- it sorts last and cannot be used to hijack a name).

-- Login: email -> the one live user (in a live clinic) whose password we verify
-- and whose clinic_id becomes the tenant context.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION auth_find_user_for_login(p_email citext)
RETURNS TABLE (
  id            uuid,
  clinic_id     uuid,
  role          text,
  full_name     text,
  password_hash text
)
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
  SELECT u.id, u.clinic_id, u.role, u.full_name, u.password_hash
  FROM users u
  JOIN clinics c ON c.id = u.clinic_id
  WHERE u.email = p_email
    AND u.deleted_at IS NULL
    AND c.deleted_at IS NULL;
$$;
-- +goose StatementEnd

COMMENT ON FUNCTION auth_find_user_for_login(citext) IS
  'SECURITY DEFINER login lookup: minimal columns for a live user, pre-tenant-context.';

-- Refresh: token hash -> the token row. revoked_at/expires_at are returned (not
-- filtered) so the caller can detect reuse of a revoked token as a security
-- signal rather than a silent miss.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION auth_find_refresh_token(p_token_hash text)
RETURNS TABLE (
  id         uuid,
  user_id    uuid,
  clinic_id  uuid,
  expires_at timestamptz,
  revoked_at timestamptz
)
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
  SELECT rt.id, rt.user_id, rt.clinic_id, rt.expires_at, rt.revoked_at
  FROM refresh_tokens rt
  WHERE rt.token_hash = p_token_hash;
$$;
-- +goose StatementEnd

COMMENT ON FUNCTION auth_find_refresh_token(text) IS
  'SECURITY DEFINER refresh lookup: token identity + validity window, pre-tenant-context.';

-- Invite acceptance: token hash -> the invitation, so an unauthenticated invitee
-- can accept it. Validity (accepted_at/expires_at) is judged by the caller.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION auth_find_invitation(p_token_hash text)
RETURNS TABLE (
  id          uuid,
  clinic_id   uuid,
  email       citext,
  role        text,
  expires_at  timestamptz,
  accepted_at timestamptz
)
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
  SELECT i.id, i.clinic_id, i.email, i.role, i.expires_at, i.accepted_at
  FROM invitations i
  WHERE i.token_hash = p_token_hash;
$$;
-- +goose StatementEnd

COMMENT ON FUNCTION auth_find_invitation(text) IS
  'SECURITY DEFINER invite lookup: minimal columns to accept an invitation, pre-tenant-context.';

-- +goose Down

DROP FUNCTION IF EXISTS auth_find_invitation(text);
DROP FUNCTION IF EXISTS auth_find_refresh_token(text);
DROP FUNCTION IF EXISTS auth_find_user_for_login(citext);

DROP POLICY IF EXISTS transcript_segments_tenant_isolation ON transcript_segments;
ALTER TABLE transcript_segments DISABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS assessments_tenant_isolation ON assessments;
ALTER TABLE assessments DISABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS media_assets_tenant_isolation ON media_assets;
ALTER TABLE media_assets DISABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS consultation_sessions_tenant_isolation ON consultation_sessions;
ALTER TABLE consultation_sessions DISABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS objective_templates_tenant_isolation ON objective_templates;
ALTER TABLE objective_templates DISABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS invitations_tenant_isolation ON invitations;
ALTER TABLE invitations DISABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS refresh_tokens_tenant_isolation ON refresh_tokens;
ALTER TABLE refresh_tokens DISABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS users_tenant_isolation ON users;
ALTER TABLE users DISABLE ROW LEVEL SECURITY;
DROP POLICY IF EXISTS clinics_tenant_isolation ON clinics;
ALTER TABLE clinics DISABLE ROW LEVEL SECURITY;
