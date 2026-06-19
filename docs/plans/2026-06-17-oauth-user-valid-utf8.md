---
title: OAuth User Response UTF-8 Integrity
type: security
status: completed
date: 2026-06-17
execution: code
---

# OAuth User Response UTF-8 Integrity

## Status

Completed.

## Problem

Go's `encoding/json` accepts malformed UTF-8 in JSON strings and replaces invalid
bytes with `U+FFFD`. The OAuth user-profile decoder currently validates JSON
structure and identity only after that replacement, so a malformed provider
response can become an apparently valid user ID before cache-key and cookie
publication.

## Priorities

1. Reject malformed UTF-8 in the bounded raw response body before JSON token
   scanning or typed decoding.
2. Cover invalid bytes in both OAuth identity values and JSON member names while
   preserving valid Unicode identities and existing size, read-error, duplicate,
   and case-folding behavior.
3. Add a maintenance contract and synchronize contributor/security guidance so
   future refactors cannot silently remove the pre-decode boundary.

## Implementation

- `auth.go`: add a stable malformed-UTF-8 error and validate the bounded body
  with the standard library before `rejectDuplicateJSONMembers`.
- `auth_test.go`: add focused malformed-value and malformed-member-name
  regressions plus a valid non-ASCII identity control.
- `scripts/check-baseline.sh`: require the validation order, tests, and completed
  plan evidence.
- `AGENTS.md`, `README.md`, `SECURITY.md`, `VISION.md`, and `CHANGES.md`: record
  the raw-response encoding boundary and its verification.

## Validation

- Run focused OAuth response tests, `go test -race -count=1 ./...`, `go vet ./...`,
  and `make check` from the repository and an external directory.
- Reject isolated mutations that remove the UTF-8 check, move it after JSON
  parsing, remove either malformed-input regression, or reopen this plan.
- Run `gofmt`, `git diff --check`, generated-artifact inspection, and a
  changed-line credential-pattern audit before committing.

## Risks

- This intentionally rejects provider responses that were previously repaired
  by `encoding/json`; valid UTF-8 JSON behavior remains unchanged.
- Live OAuth exchange and credential-bearing Foursquare requests remain outside
  local verification.

## Work Completed

- Added a dedicated malformed-UTF-8 response error and rejected invalid raw
  bytes after the 1 MiB size boundary but before duplicate-member scanning or
  typed decoding.
- Added malformed identity-value and member-name regressions plus a valid
  Unicode identity control.
- Added mutation-sensitive source, ordering, test, guidance, and completed-plan
  contracts to the maintained baseline.

## Verification Completed

- Focused UTF-8 response tests passed.
- `go test -race -count=1 ./...` passed.
- `go vet ./...` passed.
- The repository-root and external-directory `make check` gates passed.
- Five isolated hostile mutations were rejected across the raw UTF-8 check,
  validation order, malformed-value regression, malformed-member regression,
  and completed-plan evidence.
- `gofmt`, module-integrity, diff, generated-artifact, and changed-line
  credential-pattern audits passed.
- No live OAuth callback was executed and no credentials were used.
