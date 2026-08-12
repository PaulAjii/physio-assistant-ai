# Phase 2 — Clinic Templates & Objective Assessment

**Status:** ⬜ Not started
**Goal:** move objective-exam templates out of hardcoded frontend data into per-clinic DB records that
admins configure, and finish the objective assessment flow end-to-end.

## Templates

- [ ] Template CRUD API (clinic-scoped) backed by `objective_templates`.
- [ ] Admin UI to create/edit templates (categories, tests, input types, units, priority).
- [ ] Seed the shipped defaults (Knee, Shoulder, Lower Back, Neck) as editable starting points.
- [ ] Add the missing **Hip Pain** template (or let clinics create it).
- [ ] Frontend loads templates from API instead of `composables/templates.ts`.

## Objective assessment

- [ ] Finish the objective capture flow (binary / measurement / grading / notes) end-to-end.
- [ ] Persist findings against the assessment.
- [ ] Validation + consistent response envelope.

## Docs to produce on completion

- `docs/objective-templates.md`, `docs/objective-assessment.md`.
