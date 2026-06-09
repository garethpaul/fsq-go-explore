---
title: FSQ Propose Edit Form Parse Boundary
date: 2026-06-09
status: completed
execution: code
---

## Context

`ProposeEdit` already rejects non-POST requests and missing venue IDs before
auth or Foursquare API work. Malformed `application/x-www-form-urlencoded`
bodies were parsed only after auth and Foursquare client setup, because the
handler used `FormValue` first and checked `ParseForm` later.

Malformed edit bodies should fail at the request boundary.

## Goals

- Reject malformed venue edit forms with `400 Bad Request`.
- Do that before auth-cookie lookup, token cache work, or Foursquare client
  setup.
- Preserve existing POST-only and missing-ID behavior.
- Extend tests, baseline, and docs for the form-parse boundary.

## Implementation

- Parse the edit form immediately after the POST method check.
- Read the venue ID from the parsed form data.
- Return `400 Bad Request` for malformed forms instead of redirecting after
  auth setup.
- Added a regression test for malformed form data before auth.
- Extended README, SECURITY, VISION, CHANGES, and `scripts/check-baseline.sh`.

## Verification

- `go test ./...`
- `sh -n scripts/check-baseline.sh`
- `scripts/check-baseline.sh`
- `make lint`
- `make test`
- `make build`
- `make check`
- `git diff --check`
