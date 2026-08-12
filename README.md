# Physio Assistant AI

An AI-powered assistant for physiotherapists. A clinician records (or, soon, live-streams) a patient
consultation; the system transcribes and structures it into a clinical SOAP note, lets the clinician
correct it, add an objective examination (soon voice-driven) and a gait analysis, then sign and export
the finished assessment. Built for clinics in Nigeria — it understands **English, Yoruba, and Pidgin**.

> 🚧 Active development. See the [roadmap](#roadmap) below for what's built and what's next.

## Architecture

Three-tier system (a fourth service, Python vision, arrives in Phase 4):

```
Frontend (Nuxt 4)  ──WS + HTTPS──►  backend-core (Go)  ──►  backend-ai (NestJS + Gemini)
 mic + webcam                        identity · Postgres      Gemini Live + note generation
                                     JWT · S3 · WS hub
                                            │
                                            └──────────────►  backend-vision (Python)  [Phase 4]
                                                              MediaPipe pose + gait math
```

The frontend only talks to **backend-core**, which owns identity, persistence, and orchestration, and
brokers to the AI (and, later, vision) services.

## Services

| Directory | Stack | Role |
| --------- | ----- | ---- |
| [`frontend/`](frontend) | Nuxt 4 · Vue 3 · Nuxt UI 4 · Tailwind 4 | Clinician-facing web app |
| [`backend-core/`](backend-core) | Go · GoFiber v3 · (Postgres, JWT, S3 — Phase 1) | Identity, orchestration, persistence |
| [`backend-ai/`](backend-ai) | NestJS 11 · Mastra · Zod · Google Gemini 2.5 Flash | Audio → structured SOAP note |
| `backend-vision/` | Python · FastAPI · MediaPipe _(Phase 4)_ | Gait analysis from video |

## Technologies

- **Frontend:** Nuxt 4, Vue 3, Nuxt UI 4, Tailwind CSS 4, TypeScript, pnpm, ESLint.
- **Core backend:** Go, GoFiber v3, google/uuid; Postgres + JWT + S3-compatible storage (Phase 1).
- **AI backend:** NestJS 11, Mastra (`@mastra/core`), Zod, Google Generative AI (Gemini 2.5 Flash),
  Gemini Live for real-time transcription (Phase 3).
- **Vision (Phase 4):** Python, MediaPipe/OpenPose pose estimation.
- **Infra:** Docker Compose, Postgres, MinIO (local) / Cloudflare R2 (prod) — Phase 1.

## Roadmap

Work is tracked phase by phase. Each phase is a **feature checklist** file under [`planning/`](planning);
tick items with `- [x]` as they land and update the file's **Status** line. When a phase completes,
matching per-feature documentation is written under [`docs/`](docs).

**Status legend:** `- [ ]` not started · `- [x]` done · `- [~]` in progress.

| Phase | Title | Status | Checklist |
| ----- | ----- | ------ | --------- |
| 0 | Safety & hygiene | ✅ Done | [phase-0](planning/phase-0-safety-and-hygiene.md) |
| 1 | Foundation (DB · auth · multi-tenancy · storage) | ⬜ Not started | [phase-1](planning/phase-1-foundation.md) |
| 2 | Clinic templates & objective assessment | ⬜ Not started | [phase-2](planning/phase-2-templates-and-objective.md) |
| 3 | Live consultation (real-time audio + voice-driven objective) | ⬜ Not started | [phase-3](planning/phase-3-live-consultation.md) |
| 4 | Vision / gait analysis (Python service) | ⬜ Not started | [phase-4](planning/phase-4-vision-gait.md) |
| 5 | Sign & print | ⬜ Not started | [phase-5](planning/phase-5-sign-and-print.md) |
| 6 | UI overhaul | ⬜ Not started | [phase-6](planning/phase-6-ui-overhaul.md) |

### Key decisions locked in

- **Database:** local **Postgres** in Docker (self-managed; we own schema, migrations, tenant isolation).
- **Auth:** **custom JWT in Go** (backend-core is the identity authority). Signing = authenticated
  clinician id + timestamp stamped on finalize.
- **Vision/gait:** dedicated **Python ML service** (MediaPipe pose + gait metrics).
- **Real-time audio:** **Gemini Live** (via NestJS) for streaming transcription — accuracy across
  English / Yoruba / Pidgin without self-hosting a GPU. Python service stays the ML learning ground.
- **Live objective:** clinic templates are a *scaffold*; voice-driven capture fills matching fields
  as the clinician speaks, ad-hoc findings are captured too, and manual entry is always available.
- **Media storage:** S3-compatible abstraction — **MinIO** locally, **Cloudflare R2** in production.

## Documentation

Per-feature reference lives in [`docs/`](docs/README.md) and grows each phase.

- [Configuration](docs/configuration.md) — env vars & config loading (backend-core).

## Getting started

Each service currently runs independently (unified docker-compose lands in Phase 1):

```bash
# AI backend (set GOOGLE_GENERATIVE_AI_API_KEY and PORT in backend-ai/.env)
cd backend-ai && pnpm install && pnpm run start:dev

# Core backend (env vars optional; see backend-core/.env.example)
cd backend-core && go run ./cmd/api

# Frontend (set NUXT_PUBLIC_API_BASE_URI if not using the default)
cd frontend && pnpm install && pnpm dev
```

## License

See [LICENSE](LICENSE).
