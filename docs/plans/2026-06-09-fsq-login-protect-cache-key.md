---
title: FSQ Login Protect Cache Key
date: 2026-06-09
status: completed
execution: code
---

## Context

Access-token lookup already validates that auth cookie values are generated
`user:` cache keys before memcache lookup. `LoginProtect` only checked that the
cookie existed, so malformed cookie values could still reach protected handler
code before later token lookup redirected them.

## Goals

- Reject missing auth cookies in the protected-route wrapper.
- Reject malformed auth-cookie cache keys before protected handler work starts.
- Preserve generated user cache keys as the accepted auth-cookie shape.
- Add Go test and static baseline coverage for the wrapper boundary.

## Implementation

- Updated `LoginProtect` to require `validUserCacheKey(cookie.Value)`.
- Added tests for malformed auth cookies and generated user cache keys.
- Extended `scripts/check-baseline.sh`, README, SECURITY, VISION, and CHANGES.

## Verification

- `go test ./...`
- `scripts/check-baseline.sh`
- `make lint`
- `make test`
- `make build`
- `make check`
- `git diff --check`
