---
title: OAuth User Response Boundary
type: security
status: planned
date: 2026-06-13
---

# OAuth User Response Boundary

## Summary

Validate the Foursquare OAuth user-profile response status and bound its body
before JSON decoding or authentication state is created.

## Requirements

- R1. Reject every non-2xx profile response before reading its body.
- R2. Limit profile response reads to 1 MiB plus one detection byte.
- R3. Preserve read errors and reject oversized, malformed wrapper, and
  malformed user payloads without logging response content.
- R4. Keep response closure, generic redirect behavior, token caching, cookie
  settings, and the existing API client unchanged.
- R5. Add focused helper tests for 1xx/3xx/4xx/5xx, exact-limit, oversize,
  malformed JSON, and read failure behavior.
- R6. Add ordering, documentation, completion, and mutation contracts.
- R7. Do not claim live OAuth, App Engine, memcache, or Foursquare validation.

## Verification Plan

- Run focused auth tests, uncached all-package tests, race tests, vet, tidy diff,
  and all four Make gates.
- Reject status-guard removal, decode-before-status, unbounded read, limit drift,
  weakened tests, stale plan, and missing evidence mutations.
- Run gofmt, shell syntax, diff, artifact, and secret checks.
- Take one bounded exact-head push/PR/code-scanning snapshot after push.

## Non-Goals

- Changing OAuth exchange, state cookies, cache keys, or session duration.
- Changing the shared `fsq` API response decoder or client timeout.
- Adding dependencies or modernizing the App Engine runtime.
