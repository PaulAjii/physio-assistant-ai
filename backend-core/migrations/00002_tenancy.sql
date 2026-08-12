-- Tenancy and identity: the multi-tenant spine.
--
-- Every tenant-scoped table carries clinic_id and (in 00004) an RLS policy that
-- filters on current_clinic_id(). This migration only creates the tables,
-- indexes, foreign keys and updated_at triggers — RLS is enabled in one place
-- (00004) so the whole isolation model can be reviewed together.
--
-- SOFT DELETE: nothing is ever hard-deleted. A row is "live" while deleted_at
-- IS NULL and "deleted" once it is set. Foreign keys therefore use the default
-- (NO ACTION), which also guards against accidental hard deletes of a clinic
-- that still has records. deleted_at filtering lives in the query layer, not in
-- RLS, so that the UPDATE which sets deleted_at is never blocked by its own
-- policy.

-- +goose Up

-- clinics: the tenant root. Everything else references it.
CREATE TABLE clinics (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name       text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE TRIGGER clinics_set_updated_at
  BEFORE UPDATE ON clinics
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- users: clinic admins and clinicians.
CREATE TABLE users (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  clinic_id     uuid NOT NULL REFERENCES clinics(id),
  email         citext NOT NULL,
  password_hash text NOT NULL,
  full_name     text NOT NULL,
  role          text NOT NULL CHECK (role IN ('admin', 'clinician')),
  license_no    text,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now(),
  deleted_at    timestamptz
);

-- Email is unique GLOBALLY, because login happens before the clinic is known
-- (email -> user -> clinic_id -> tenant context). The index is partial so a
-- soft-deleted user's address can be re-used by a new account.
CREATE UNIQUE INDEX users_email_live_uniq ON users (email) WHERE deleted_at IS NULL;
CREATE INDEX users_clinic_id_idx ON users (clinic_id);

CREATE TRIGGER users_set_updated_at
  BEFORE UPDATE ON users
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- refresh_tokens: server-side half of JWT refresh rotation. Only the SHA-256
-- hash of the token is stored, never the token itself, so a database leak
-- cannot be replayed as a live session. clinic_id is denormalised from the
-- user so RLS can isolate these rows without a join.
CREATE TABLE refresh_tokens (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    uuid NOT NULL REFERENCES users(id),
  clinic_id  uuid NOT NULL REFERENCES clinics(id),
  token_hash text NOT NULL,
  expires_at timestamptz NOT NULL,
  revoked_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX refresh_tokens_token_hash_uniq ON refresh_tokens (token_hash);
CREATE INDEX refresh_tokens_user_id_idx ON refresh_tokens (user_id);

-- invitations: an admin invites someone into their clinic. The invitee accepts
-- with the token (delivered out-of-band for now — email is not wired up yet),
-- sets a password, and that creates their users row. Token stored as a hash,
-- same reasoning as refresh_tokens.
CREATE TABLE invitations (
  id          uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  clinic_id   uuid NOT NULL REFERENCES clinics(id),
  email       citext NOT NULL,
  role        text NOT NULL CHECK (role IN ('admin', 'clinician')),
  token_hash  text NOT NULL,
  invited_by  uuid REFERENCES users(id),
  expires_at  timestamptz NOT NULL,
  accepted_at timestamptz,
  created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX invitations_token_hash_uniq ON invitations (token_hash);
CREATE INDEX invitations_clinic_id_idx ON invitations (clinic_id);
-- At most one pending (not-yet-accepted) invitation per email per clinic.
CREATE UNIQUE INDEX invitations_pending_uniq ON invitations (clinic_id, email)
  WHERE accepted_at IS NULL;

-- +goose Down
-- Reverse dependency order. Triggers drop automatically with their tables.
DROP TABLE IF EXISTS invitations;
DROP TABLE IF EXISTS refresh_tokens;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS clinics;
