# Phase 4 — Vision / Gait Analysis (Python Service)

**Status:** ⬜ Not started
**Goal:** add a Python ML service that analyses live video of a patient's gait via pose estimation and
returns quantitative metrics attached to the assessment. (Also the primary Python/ML learning vehicle.)

## Service

- [ ] New `backend-vision` Python service (FastAPI + MediaPipe/OpenPose).
- [ ] Add to docker-compose; env + healthcheck.

## Capture & streaming

- [ ] Webcam capture in the frontend.
- [ ] Stream frames to the vision service (WebSocket frames for MVP; evaluate WebRTC later for latency).
- [ ] Live skeleton overlay in the UI.

## Gait metrics

- [ ] Pose keypoint extraction per frame.
- [ ] Metrics: cadence, stride length/symmetry, joint ROM angles, stance/swing timing.
- [ ] Aggregate into a gait report; store to `assessments.gait_analysis` and video to S3.

## Docs to produce on completion

- `docs/vision-service.md`, `docs/gait-analysis.md`.
