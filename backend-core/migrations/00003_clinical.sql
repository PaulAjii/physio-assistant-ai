-- Clinical tables: the consultation workflow the scribe produces.
--
-- All tenant-scoped (clinic_id) and soft-deleted (deleted_at), same rules as
-- 00002. RLS is still deferred to 00004. JSONB is used where the shape is
-- clinic-defined or evolving (template fields, assessment sections) so we don't
-- migrate the schema every time a clinic tweaks a form.

-- +goose Up

-- objective_templates: clinic-configurable scaffolds for the objective exam.
-- `body` is an ordered array of field definitions, e.g.
--   [{"key":"knee_flexion","label":"Knee flexion","type":"number","unit":"deg"},
--    {"key":"gait","label":"Gait","type":"select","options":["antalgic","normal"]}]
-- The clinician speaks findings; matching keys get filled, the rest entered by
-- hand. Shape is owned by the clinic, hence jsonb rather than columns.
CREATE TABLE objective_templates (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  clinic_id  uuid NOT NULL REFERENCES clinics(id),
  name       text NOT NULL,
  body       jsonb NOT NULL DEFAULT '[]'::jsonb,
  is_active  boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE INDEX objective_templates_clinic_id_idx ON objective_templates (clinic_id);
-- No two live templates with the same name inside one clinic.
CREATE UNIQUE INDEX objective_templates_name_uniq ON objective_templates (clinic_id, name)
  WHERE deleted_at IS NULL;

CREATE TRIGGER objective_templates_set_updated_at
  BEFORE UPDATE ON objective_templates
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- consultation_sessions: one physiotherapy visit.
--
-- SIGNING: on finalize we stamp signed_by (the authenticated clinician) and
-- signed_at. The CHECK makes an unsigned finalized session impossible at the
-- database level — so "who signed this" is always answerable and always tied to
-- a real user, which is the whole point of the auth requirement.
CREATE TABLE consultation_sessions (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  clinic_id    uuid NOT NULL REFERENCES clinics(id),
  clinician_id uuid NOT NULL REFERENCES users(id),
  patient_ref  text,                    -- free-text name / MRN for now
  status       text NOT NULL DEFAULT 'draft'
                 CHECK (status IN ('draft', 'recording', 'finalized')),
  language     text,                    -- 'en' | 'yo' | 'pcm' | ... (open-ended)
  started_at   timestamptz,
  finalized_at timestamptz,
  signed_by    uuid REFERENCES users(id),
  signed_at    timestamptz,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  deleted_at   timestamptz,
  CONSTRAINT consultation_sessions_finalized_is_signed
    CHECK (status <> 'finalized'
           OR (signed_by IS NOT NULL AND signed_at IS NOT NULL AND finalized_at IS NOT NULL))
);

CREATE INDEX consultation_sessions_clinic_id_idx ON consultation_sessions (clinic_id);
CREATE INDEX consultation_sessions_clinician_id_idx ON consultation_sessions (clinician_id);

CREATE TRIGGER consultation_sessions_set_updated_at
  BEFORE UPDATE ON consultation_sessions
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- media_assets: audio/video recordings. The blob lives in MinIO/R2; Postgres
-- only holds the object key we presign on demand. Nothing binary in the DB.
CREATE TABLE media_assets (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  clinic_id    uuid NOT NULL REFERENCES clinics(id),
  session_id   uuid NOT NULL REFERENCES consultation_sessions(id),
  kind         text NOT NULL CHECK (kind IN ('audio', 'video')),
  storage_key  text NOT NULL,
  content_type text,
  bytes        bigint,
  status       text NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'uploaded', 'failed')),
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now(),
  deleted_at   timestamptz
);

CREATE UNIQUE INDEX media_assets_storage_key_uniq ON media_assets (storage_key);
CREATE INDEX media_assets_session_id_idx ON media_assets (session_id);
CREATE INDEX media_assets_clinic_id_idx ON media_assets (clinic_id);

CREATE TRIGGER media_assets_set_updated_at
  BEFORE UPDATE ON media_assets
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- assessments: the structured note (replaces the in-memory store). One live row
-- per session; `version` bumps on each edit so pre-finalize changes are visible.
-- subjective/objective are jsonb sections; pain_intensity is nullable to match
-- the "not reported" case (models.Intensity *int) and bounded to the 0-10 scale.
CREATE TABLE assessments (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  clinic_id      uuid NOT NULL REFERENCES clinics(id),
  session_id     uuid NOT NULL REFERENCES consultation_sessions(id),
  subjective     jsonb NOT NULL DEFAULT '{}'::jsonb,
  objective      jsonb NOT NULL DEFAULT '{}'::jsonb,
  ai_summary     text,
  pain_intensity int CHECK (pain_intensity IS NULL OR pain_intensity BETWEEN 0 AND 10),
  version        int NOT NULL DEFAULT 1,
  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now(),
  deleted_at     timestamptz
);

-- One live assessment per session.
CREATE UNIQUE INDEX assessments_session_uniq ON assessments (session_id)
  WHERE deleted_at IS NULL;
CREATE INDEX assessments_clinic_id_idx ON assessments (clinic_id);

CREATE TRIGGER assessments_set_updated_at
  BEFORE UPDATE ON assessments
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- transcript_segments: output of live speech-to-text. Append-only event log;
-- interim results are stored with is_final=false and superseded by final ones.
-- No soft-delete: these are immutable records that belong to their session.
CREATE TABLE transcript_segments (
  id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  clinic_id  uuid NOT NULL REFERENCES clinics(id),
  session_id uuid NOT NULL REFERENCES consultation_sessions(id),
  speaker    text,                       -- 'clinician' | 'patient' | NULL
  content    text NOT NULL,
  start_ms   int,
  end_ms     int,
  is_final   boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX transcript_segments_session_id_idx ON transcript_segments (session_id);
CREATE INDEX transcript_segments_clinic_id_idx ON transcript_segments (clinic_id);

-- +goose Down
-- Reverse dependency order. Triggers drop automatically with their tables.
DROP TABLE IF EXISTS transcript_segments;
DROP TABLE IF EXISTS assessments;
DROP TABLE IF EXISTS media_assets;
DROP TABLE IF EXISTS consultation_sessions;
DROP TABLE IF EXISTS objective_templates;
