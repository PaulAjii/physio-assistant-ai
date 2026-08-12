# Phase 1 — Foundation (DB · Auth · Multi-tenancy · Storage)

**Status:** ⬜ Not started
**Goal:** make the current app persistent, authenticated, and multi-tenant — with a data model that
already anticipates live sessions, gait, and signing, so later phases don't force a rebuild.

## Infrastructure

- [ ] `docker-compose.yml` at repo root: Postgres + MinIO (+ backend services wired via env).
- [ ] Migration tooling (goose or golang-migrate) + `make migrate` targets.
- [ ] `.env.example` for every service; document the full env matrix.

## Database schema (forward-compatible)

- [ ] `clinics` (tenant root).
- [ ] `users` (clinic_id, email, password_hash, name, role, license_no).
- [ ] `refresh_tokens`.
- [ ] `objective_templates` (clinic_id, complaint, categories JSONB).
- [ ] `consultation_sessions` (clinic_id, clinician_id, complaint, status).
- [ ] `media_assets` (session_id, type, s3_key, mime, duration).
- [ ] `assessments` (clinic_id, session_id, clinician_id, subjective/objective/gait/ai_draft JSONB,
      status, signed_by, signed_at, signature_hash).
- [ ] `transcript_segments` (session_id, speaker, text, ts).

## Persistence layer

- [ ] Repository interfaces + Postgres implementations.
- [ ] Replace in-memory `Store` (jobs) and `Assessments` stores.
- [ ] Keep the existing upload → AI → assessment flow working, now persisted.

## Auth (custom JWT in Go)

- [ ] Register clinic + admin; login; refresh; logout.
- [ ] Argon2/bcrypt password hashing.
- [ ] JWT access + refresh, rotation, middleware.
- [ ] RBAC (admin vs clinician) + tenant scoping on every query.
- [ ] Tighten CORS to known origins.

## Storage (S3-compatible)

- [ ] Storage client interface; MinIO (local) / R2 (prod) implementations.
- [ ] Presigned upload/download; move audio off local disk.

## Docs to produce on completion

- `docs/database.md`, `docs/auth.md`, `docs/storage.md`, `docs/multi-tenancy.md`.
