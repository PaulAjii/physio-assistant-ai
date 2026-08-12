# Configuration (backend-core)

Introduced in **Phase 0**. `backend-core` no longer hardcodes ports or the AI backend URL — all
configuration is read from environment variables via `internal/config`.

## The config package

`backend-core/internal/config/config.go` exposes a `Config` struct and `Load()` which reads env vars
with sensible defaults. It is the single place configuration enters the app and will grow (DB DSN, JWT
secret, S3 credentials) in later phases.

```go
cfg := config.Load()
// cfg.Port, cfg.AIBackendURI, cfg.UploadDir, cfg.AllowedOrigins
```

## Environment variables

| Variable | Default | Description |
| -------- | ------- | ----------- |
| `PORT` | `8080` | Port the HTTP server listens on. |
| `AI_BACKEND_URI` | `http://localhost:5000/ai/process-audio` | Full URL of the NestJS AI backend's process-audio endpoint. |
| `UPLOAD_DIR` | `uploads` | Directory for temporarily saved audio uploads. |
| `CORS_ALLOWED_ORIGINS` | `*` | Comma-separated list of allowed origins. Locked down to real origins in Phase 1. |

See [`backend-core/.env.example`](../backend-core/.env.example) for a copy-paste template.

## Port-mismatch note

Historically backend-core hardcoded `http://localhost:5000/...` while the NestJS backend defaults to
port `3000` (its `.env` sets `PORT=5000`). Both sides are now env-driven:

- **backend-core:** set `AI_BACKEND_URI` to wherever the AI backend runs.
- **backend-ai:** set `PORT` (its `.env` already uses `5000`).

Keep these consistent. In Phase 1 the docker-compose network will set both from one source of truth.

## Secrets

- `backend-ai/.env` (holds `GOOGLE_GENERATIVE_AI_API_KEY`) is **gitignored** and was never committed.
- Never commit real `.env` files — commit `.env.example` templates only.
