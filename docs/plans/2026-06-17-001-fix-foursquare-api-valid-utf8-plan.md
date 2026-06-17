---
title: "fix: Reject malformed Foursquare API UTF-8"
type: fix
date: 2026-06-17
---

# fix: Reject malformed Foursquare API UTF-8

## Summary

Reject malformed UTF-8 in bounded Foursquare venue/search JSON bodies before
envelope decoding so invalid provider bytes cannot become replacement
characters in rendered venue data.

## Problem Frame

The OAuth user-profile decoder validates the raw bounded body before JSON
tokenization because Go's JSON decoder replaces malformed UTF-8 inside strings.
The general Foursquare response decoder lacks the same boundary. Search and
venue-detail responses can therefore return silently altered names, addresses,
or other text to page and edit consumers.

## Requirements

- R1. General Foursquare JSON bodies must be valid UTF-8 after the existing
  2 MiB size check and before either envelope or target decoding.
- R2. Malformed UTF-8 in venue text or JSON member names must return a stable
  decode error and leave service results empty.
- R3. Valid Unicode venue data must continue to decode without behavior drift.
- R4. Maintained checks must prove the ordering and regression coverage so the
  boundary cannot be removed or moved after JSON decoding unnoticed.
- R5. Security and contributor guidance must describe the general API boundary
  without weakening the existing OAuth-specific contract.

## Key Technical Decisions

- KTD1. Validate the complete bounded byte slice in `decodeFoursquareResponse`
  before `json.Unmarshal`; this covers search and venue-detail consumers at the
  shared ingestion point.
- KTD2. Keep the existing size-limit precedence so oversized bodies remain
  classified as too large even if their trailing bytes are malformed.
- KTD3. Return a package-level sentinel error for malformed UTF-8 so tests and
  logs can distinguish integrity rejection from generic JSON syntax errors.
- KTD4. Preserve the OAuth decoder as an independent policy boundary; sharing
  helpers would add coupling without reducing meaningful complexity.

## Implementation Units

### U1. Enforce raw API response UTF-8 integrity

- **Goal:** Validate bounded venue/search response bytes before envelope
  decoding and expose a stable malformed-UTF-8 error.
- **Files:** `fsq/api.go`, `fsq/api_test.go`
- **Patterns:** Follow the ordering and valid-Unicode behavior established by
  the OAuth response decoder while retaining the general API's 2 MiB limit.
- **Test scenarios:** Reject malformed bytes inside a venue value; reject
  malformed bytes inside a JSON member name; accept a valid non-ASCII venue;
  preserve oversized-body precedence.
- **Verification:** Focused package tests and the repository's full Go/static
  maintenance gate pass.

### U2. Make the boundary durable in repository contracts

- **Goal:** Record the security behavior and add mutation-sensitive static
  checks for validation presence, ordering, regressions, and plan evidence.
- **Files:** `scripts/check-baseline.sh`, `README.md`, `SECURITY.md`, `VISION.md`,
  `AGENTS.md`, `CHANGES.md`, `docs/plans/2026-06-17-001-fix-foursquare-api-valid-utf8-plan.md`
- **Patterns:** Extend existing response-boundary contracts without duplicating
  runtime logic in shell.
- **Test scenarios:** Removing the raw validation, moving it after unmarshal,
  deleting either malformed-input regression, or omitting completion evidence
  must make the maintained gate fail.
- **Verification:** Root and external-directory gates pass, and isolated hostile
  mutations are rejected.

## Scope Boundaries

- No changes to response size limits, status/content-type/final-URL policy,
  redirects, client timeouts, or request construction.
- No duplicate-member or case-folded-key expansion for general venue payloads.
- No live Foursquare request, browser callback, or credential-bearing flow.
- PR stacking remains based on the terminal OAuth UTF-8 branch; no existing PR
  is merged or closed.

## Risks

- Provider payloads containing malformed UTF-8 that were previously accepted
  will now produce empty service results and a decode log, which is the intended
  fail-closed behavior.
- Local tests use synthetic responses and do not prove current provider output;
  hosted checks remain authoritative for exact-head integration evidence.
