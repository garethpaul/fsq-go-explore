---
title: OAuth User Duplicate JSON Members
type: security
status: planned
date: 2026-06-15
execution: code
---

# OAuth User Duplicate JSON Members

## Problem

Go's standard JSON unmarshaler accepts duplicate object members with
last-value-wins behavior. An OAuth user-profile response can therefore provide
ambiguous `response`, `user`, or `id` members before the selected identity is
cached and published in the session cookie.

## Approach

- Scan the bounded JSON body token by token before typed unmarshalling.
- Track member names independently for every nested object and reject the
  first duplicate with a stable error.
- Preserve the existing body limit, origin, media-type, identity, caching, and
  cookie behavior for valid responses.
- Add nested duplicate-member regressions plus mutation-sensitive source,
  guidance, and completed-plan contracts.

## Files

- `auth.go`
- `auth_test.go`
- `scripts/check-baseline.sh`
- `README.md`
- `SECURITY.md`
- `VISION.md`
- `CHANGES.md`
- `AGENTS.md`
- `docs/plans/2026-06-15-oauth-user-duplicate-json-members.md`

## Verification

- Prove the prior decoder accepts duplicate `response`, `user`, and `id`
  members, then prove all three fail before identity publication.
- Run focused tests, repository-root and external-directory `make check`, race
  tests, and `go vet` with explicit timeouts.
- Reject isolated scanner, integration, nested-regression, guidance, and
  completed-plan mutations.
- Audit the exact diff, generated artifacts, and credential-shaped additions.

## Non-Goals

- Do not reject unknown unique response fields or change accepted JSON media
  types.
- Do not contact live Foursquare services or use credentials.
- Do not merge or close stacked pull requests without owner authorization.
