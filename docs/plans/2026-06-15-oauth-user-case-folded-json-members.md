---
title: OAuth User Case-Folded JSON Members
type: security
status: in_progress
date: 2026-06-15
execution: code
---

# OAuth User Case-Folded JSON Members

## Problem

The bounded OAuth user-profile scanner rejects byte-identical duplicate object
members, but Go's typed JSON decoder also matches struct fields using Unicode
case folding. Distinct names such as `id` and `ID` can therefore pass the scan
and still target the same identity field with last-value-wins behavior.

## Approach

- Fold every object member name with semantics equivalent to
  `encoding/json` field matching before recording it in the per-object set.
- Reject case-variant ambiguity at any object nesting level before typed
  decoding, caching, or cookie publication.
- Preserve exact duplicate rejection, unique unknown fields, body/origin/media
  boundaries, and valid OAuth identity behavior.
- Add focused ASCII and Unicode case-fold regressions plus mutation-sensitive
  source, guidance, and completed-plan contracts.

## Files

- `auth.go`
- `auth_test.go`
- `scripts/check-baseline.sh`
- `README.md`
- `SECURITY.md`
- `VISION.md`
- `CHANGES.md`
- `AGENTS.md`
- `docs/plans/2026-06-15-oauth-user-case-folded-json-members.md`

## Verification

- Prove the current scanner accepts case-variant names that typed decoding maps
  to the same response, user, and identity fields.
- Run focused tests, repository-root and external-directory `make check`, race
  tests, and `go vet` with explicit timeouts.
- Reject isolated folding, integration, nested-regression, guidance, and
  completed-plan mutations.
- Audit the exact diff, generated artifacts, and credential-shaped additions.

## Non-Goals

- Do not change the accepted OAuth response schema, media types, or identity
  rules beyond decoder-aligned ambiguous-member rejection.
- Do not contact live Foursquare services or use credentials.
- Do not merge or close stacked pull requests without owner authorization.

## Status: In Progress
