---
title: Foursquare Response Status Validation
type: security
status: completed
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

- Added one shared `successfulFoursquareStatus` predicate for the exact 2xx
  range.
- Returned the existing empty search or venue detail result before bounded JSON
  decoding when Foursquare returns a non-2xx status.
- Preserved response-body closure, numeric-only status logging, the 2 MiB
  decoder, request parameters, public method signatures, and escaped venue IDs.
- Added direct transport tests with valid-looking error envelopes and a
  199/200/299/300 predicate boundary test plus a method-scoped static
  guard/decoder ordering contract.
- Updated repository guidance and completed-plan enforcement for the status
  boundary.

## Verification Completed

- The status guard mutation failed after bypassing the search status predicate;
  the crafted 500 response populated a venue and failed the focused test.
- The decode ordering mutation failed after moving venue detail decoding ahead
  of status validation; the crafted 502 response populated a venue ID.
- The non-2xx test mutation failed after changing the crafted search response to
  HTTP 200; the test detected that its error-path premise had been weakened.
- Focused `fsq` tests and uncached all-package tests passed on Go 1.25.3.
- All packages passed with the race detector; `go vet ./...` and
  `go mod tidy -diff` passed.
- `make check`, `make lint`, `make test`, and `make build` passed.
- gofmt, shell syntax, `git diff --check`, and intended-path artifact and secret
  scans passed.
- The hosted pull-request check and code-scanning results are recorded against
  the exact pushed head in the external engineering tracker. Embedding that SHA
  here would create a new head without exact-head hosted evidence.
