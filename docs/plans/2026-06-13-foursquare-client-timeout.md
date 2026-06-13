---
title: Foursquare Client Timeout
type: reliability
status: planned
date: 2026-06-13
---

# Foursquare Client Timeout

## Summary

Apply a bounded end-to-end timeout to Foursquare HTTP clients that do not
already declare one, while preserving explicit caller timeouts and configuration
ownership.

## Priority

1. Bound stalled search, details, and venue-edit requests.
2. Preserve explicit positive caller timeout choices.
3. Avoid mutating the caller-owned `FoursquareConfig` or `http.Client` value.

## Requirements

- R1. The package must define one `10 * time.Second` default request timeout.
- R2. `NewFoursquareService` must clone the supplied config and client value.
- R3. A non-positive client timeout must become the package default.
- R4. An explicit positive timeout must remain unchanged.
- R5. Construction must not mutate the caller's config or client timeout.
- R6. Focused tests must cover defaulting, explicit preservation, and caller
  immutability.
- R7. Static contracts must reject timeout removal, value drift, unconditional
  override, caller mutation, weakened tests, documentation drift, or incomplete
  plan evidence.
- R8. README, SECURITY, VISION, CHANGES, and AGENTS must describe the default
  timeout and the continuing absence of live integration validation.

## Non-Goals

- Adding per-operation timeout values, retries, backoff, or circuit breakers.
- Changing response status, body-size, decode, cache, OAuth, or rate-limit
  behavior.
- Replacing the legacy App Engine transport or changing deployment metadata.
- Making live Foursquare requests or exercising OAuth credentials.

## Implementation Units

### 1. Service Construction

Files: `fsq/api.go`

- Clone the config and client, default only non-positive timeout values, and
  retain the cloned config on the service.

### 2. Focused Tests

Files: `fsq/api_test.go`

- Verify default, explicit, and immutable construction behavior.

### 3. Static Contracts

Files: `scripts/check-baseline.sh`

- Require exact value, conditional defaulting, clone ownership, tests,
  documentation, and completed evidence.

### 4. Repository Guidance

Files: `README.md`, `SECURITY.md`, `VISION.md`, `CHANGES.md`, `AGENTS.md`

- Record the request timeout and remaining live/deployment limits.

## Verification Plan

- Run focused fsq tests, all-package tests, race tests, `go vet ./...`,
  `go mod tidy -diff`, and all four Make gates.
- Remove the timeout, drift its value, override explicit timeouts, mutate the
  caller config, weaken focused tests, and regress plan evidence; each mutation
  must fail.
- Run gofmt, shell syntax, `git diff --check`, and intended-file secret/artifact
  scans.
- Take bounded exact-head push, pull-request, and code-scanning snapshots after
  push; do not start a watch loop.

## Work Completed

Pending implementation.

## Verification Completed

Pending implementation and verification.
