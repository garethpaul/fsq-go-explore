---
title: OAuth User Content-Type Cardinality
type: security
status: completed
date: 2026-06-15
execution: code
---

# OAuth User Content-Type Cardinality

## Problem

OAuth user-profile decoding reads `Content-Type` through `http.Header.Get`.
When a response contains multiple media-type fields, Go returns one value and
the decoder can accept ambiguous metadata before processing an authentication
identity.

## Approach

- Require exactly one `Content-Type` field before parsing its media type.
- Continue accepting the existing JSON and structured `+json` media types when
  represented by one unambiguous field.
- Add focused duplicate, combined, missing, and valid-field coverage plus
  mutation-sensitive source, guidance, and plan contracts.

## Files

- `auth.go`
- `auth_test.go`
- `scripts/check-baseline.sh`
- `README.md`
- `SECURITY.md`
- `VISION.md`
- `CHANGES.md`
- `AGENTS.md`
- `docs/plans/2026-06-15-oauth-user-content-type-cardinality.md`

## Verification

- Prove duplicate metadata is accepted by the prior implementation, then
  rejected after the cardinality guard.
- Run focused tests, repository and external-directory `make check`, race tests,
  and `go vet` with explicit timeouts.
- Reject isolated source, regression, guidance, and completed-plan mutations.
- Audit the exact diff, generated artifacts, and secret patterns.

## Non-Goals

- Do not change OAuth state, token exchange, response body limits, accepted JSON
  media types, identity validation, caching, or cookie publication.
- Do not contact live Foursquare services or use credentials.
- Do not merge or close stacked pull requests without owner authorization.

## Status: Completed

## Work Completed

- Required exactly one OAuth user-profile `Content-Type` field before media
  parsing or response-body reads.
- Rejected duplicate and comma-combined metadata while preserving one
  structured `application/*+json` field.
- Added source, regression, guidance, and completed-plan baseline contracts.

## Verification Completed

- Focused OAuth media-type cardinality tests passed after the duplicate-field
  regression failed against the prior implementation.
- Repository and external-directory Make gates passed.
- `go test -race -count=1 ./...` and `go vet ./...` passed.
- Six hostile mutations failed for field cardinality, indexed parsing,
  duplicate coverage, combined coverage, guidance, and plan evidence.
- Exact diff, generated-artifact, conflict-marker, and secret-pattern audits passed.
- No live OAuth callback was executed.
