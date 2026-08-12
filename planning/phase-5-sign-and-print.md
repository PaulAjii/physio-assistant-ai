# Phase 5 — Sign & Print

**Status:** ⬜ Not started
**Goal:** let an authenticated clinician finalize/sign an assessment (verifiable identity + timestamp)
and export/print the signed record.

## Signing

- [ ] Finalize action: stamp `signed_by` (authenticated clinician), `signed_at`, `signature_hash`.
- [ ] Signed assessments become immutable (edits create an amendment / new version).
- [ ] Verification view: who signed, when, integrity check via hash.

## Print / PDF

- [ ] Print-optimized assessment view (clean layout, clinic letterhead).
- [ ] Server-side PDF export.
- [ ] Include gait report + objective findings + signature block.

## Docs to produce on completion

- `docs/signing.md`, `docs/print-export.md`.
