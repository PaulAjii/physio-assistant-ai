# Phase 0 — Safety & Hygiene

**Status:** ✅ Done
**Goal:** de-risk the repo and fix the cheap, high-impact bugs before feature work begins.

## Secrets & repo hygiene

- [x] Confirm `backend-ai/.env` is **not** tracked and **not** in git history (verified: never committed).
- [x] Confirm `main.exe` is **not** tracked and **not** in git history (verified: only on disk, ignored).
- [x] Harden the **root `.gitignore`** (env files, build artifacts, OS files, uploads, binaries).
- [x] Add a proper **`backend-core/.gitignore`** (Go binaries, `.env`, `uploads/`).
- [x] Add **`backend-core/.env.example`** documenting required env vars.

> ⚠️ Follow-up for the maintainer: the Google API key in `backend-ai/.env` was never committed, so no
> history scrub is needed. Rotate it only if it was ever shared/exposed elsewhere.

## Configuration (kills the port-mismatch bug)

- [x] Add `backend-core/internal/config` package (`Config` struct + `Load()` from env with defaults).
- [x] `PORT` env-driven (default `8080`).
- [x] `AI_BACKEND_URI` env-driven (default `http://localhost:5000/ai/process-audio`) — no more hardcoded URL.
- [x] `UPLOAD_DIR` env-driven (default `uploads`).
- [x] `CORS_ALLOWED_ORIGINS` env-driven (default `*`; real lockdown lands in Phase 1).
- [x] Inject config into `ConsultationHandler` instead of hardcoding.

## Bug fixes

- [x] `AssessmentStore.GetAssessment` — fix double `RLock()` (was a latent deadlock) → `RUnlock()`.
- [x] Pain `intensity` contract bug — Go `int` → `*int` (nullable) to match the AI's `number().nullable()`
      schema so a `null` no longer breaks `json.Unmarshal`.
- [x] Frontend: `PainProfile.intensity` typed `number | null`; `painColor()` and badge labels made null-safe.
- [x] Fix malformed `frontend/eslint.config.mjs` (invalid object literal passed to `withNuxt()`).

## Verification

- [x] `go build ./...` passes in `backend-core`.
- [ ] Frontend `pnpm lint` + `pnpm typecheck` — **run in CI** (node_modules not installed locally).

## Docs produced

- [`../docs/configuration.md`](../docs/configuration.md) — env vars & config approach for backend-core.

## Notes / deferred to later phases

- Full CORS lockdown to known origins → Phase 1 (needs the deploy origins + auth).
- Replacing in-memory stores with Postgres → Phase 1.
- Uploaded files cleanup / move to object storage → Phase 1.
