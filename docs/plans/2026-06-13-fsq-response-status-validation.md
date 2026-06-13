---
title: Foursquare Response Status Validation
type: security
status: planned
date: 2026-06-13
---

# Foursquare Response Status Validation

## Summary

Reject non-success Foursquare HTTP responses before decoding their JSON bodies
into venue search or venue detail result structures.

## Priority

1. Keep error response payloads out of successful application data.
2. Preserve the existing 2 MiB response-body limit and generic logging policy.
3. Add direct transport tests without contacting Foursquare or App Engine.

## Requirements

- R1. Venue search and venue detail requests must accept only HTTP status codes
  in `200..<300` before response-body decoding begins.
- R2. Non-2xx responses must return the existing empty result value without
  decoding an error body into venue data.
- R3. The status rejection path must close the response body and log only the
  numeric status, without URLs, query parameters, credentials, tokens, or body
  content.
- R4. Successful 2xx responses must retain the existing 2 MiB bounded decoder,
  envelope handling, request parameter behavior, and escaped venue IDs.
- R5. Venue edit response behavior is outside this change because it does not
  decode a success payload; its status logging and body drain remain unchanged.
- R6. Focused tests must prove search and detail error envelopes cannot populate
  result structs, while successful transport behavior remains covered.
- R7. Static contracts must reject status-guard removal, guard-after-decode
  ordering, weakened non-2xx tests, and incomplete plan evidence.
- R8. README, SECURITY, VISION, CHANGES, and AGENTS must describe the status
  boundary without claiming live Foursquare or App Engine validation.

## Non-Goals

- Changing public method signatures, returning new error values, retrying
  requests, or adding user-facing error states.
- Changing the Foursquare endpoints, request parameters, OAuth tokens, cache
  behavior, rate limiting, response-size ceiling, or JSON models.
- Adding a global HTTP timeout or modernizing the legacy App Engine transport;
  those remain separate focused follow-ups.
- Exercising real Foursquare credentials, OAuth flows, APIs, memcache, or App
  Engine deployment.

## Implementation Units

### 1. Shared Status Boundary

Files: `fsq/api.go`

- Add a small success-status predicate or equivalent shared boundary.
- Return before `decodeFoursquareResponse` for non-2xx search and detail
  responses.

### 2. Focused Transport Tests

Files: `fsq/api_test.go`

- Return crafted non-2xx JSON envelopes that resemble valid search/detail
  payloads.
- Assert the returned search/detail results remain empty.
- Keep successful escaped-path and form-encoding coverage intact.

### 3. Static Contracts

Files: `scripts/check-baseline.sh`

- Require both status guards ahead of both decoder calls and retain named
  non-2xx tests.
- Require completed mutation and verification evidence.

### 4. Repository Guidance

Files: `README.md`, `SECURITY.md`, `VISION.md`, `CHANGES.md`, `AGENTS.md`

- Record the fail-closed status boundary and continuing live/deployment limits.

## Verification Plan

- Run the focused `fsq` tests, all-package tests, race tests, `go vet ./...`,
  `go mod tidy -diff`, and all four Make gates.
- Remove one status guard, move a decoder before its guard, and weaken a
  non-2xx test to return 200; the unit/static gates must reject each mutation.
- Run gofmt, shell syntax, `git diff --check`, and intended-path artifact and
  secret scans.
- Take bounded exact-head push, pull-request, and code-scanning snapshots after
  push; do not start a polling or watch loop.

## Work Completed

Pending implementation.

## Verification Completed

Pending implementation and verification.
