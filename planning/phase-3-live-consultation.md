# Phase 3 — Live Consultation (Real-time Audio + Voice-driven Objective)

**Status:** ⬜ Not started
**Goal:** replace the fire-and-forget upload + one-shot SSE model with a live session over WebSockets:
streaming transcription and objective findings filled by voice as the clinician speaks.

## Transport

- [ ] WebSocket hub in backend-core; session lifecycle (`live` → `processing` → `complete`).
- [ ] Authenticated WS connections (JWT); reconnect handling.

## Real-time audio

- [ ] Browser mic capture + chunked streaming (MediaRecorder / AudioWorklet).
- [ ] Gemini Live integration in NestJS; stream partial transcripts back.
- [ ] Persist `transcript_segments` + audio to storage.

## Voice-driven objective + subjective

- [ ] Real-time structured extraction: transcript → subjective fields.
- [ ] Map spoken findings onto template fields (`{test, value, unit, result}`); capture ad-hoc findings.
- [ ] Manual entry fallback always available.
- [ ] Live UI: transcript pane + live-filling subjective/objective panels.

## Docs to produce on completion

- `docs/live-session.md`, `docs/realtime-transcription.md`, `docs/voice-objective-extraction.md`.
