---
title: Foursquare Response Body Parse Limit
date: 2026-06-13
status: in_progress
execution: code
---

## Context

`decodeFoursquareResponse` uses `io.ReadAll` for search and venue-detail
responses. Those bodies are controlled by the remote service or an
intermediary, so an unexpectedly large response can grow process memory without
an application-level bound before JSON parsing.

## Priority

This is the highest-impact remaining isolated request-path risk because it is a
direct unbounded allocation reachable through normal application traffic.
Dependency modernization is broader and does not by itself close this runtime
resource-exhaustion boundary.

## Prioritized Backlog

1. Limit Foursquare JSON response reads to 2 MiB before unmarshalling.
2. Return a stable sentinel error for oversized bodies and preserve existing
   behavior for valid, empty, malformed, and reader-error responses.
3. Cover exact-limit, one-byte-over-limit, and short-read failure cases.
4. Add static and hostile-mutation contracts plus repository guidance.
5. Handle response media-type validation, client deadlines, and App Engine
   modernization in separate focused changes.

## Implementation

- Add a package-level response-size constant and exported sentinel error in
  `fsq/api.go`.
- Read through an `io.LimitReader` capped at limit plus one byte, reject an
  observed body above the limit, and only then unmarshal the Foursquare
  envelope and response payload.
- Add focused package tests that exercise the decoder without network access.
- Extend `scripts/check-baseline.sh`, README, SECURITY, VISION, and CHANGES to
  preserve the parse boundary and truthful completed verification evidence.

## Verification Plan

- Run focused package tests, all-package tests with and without the race
  detector, vet, module-tidy diff, all four Make gates, formatting, diff, and
  intended-file secret checks.
- Remove the limit reader, change the limit, and remove the oversize test; each
  hostile mutation must fail the maintained gate.
- Take one bounded exact-head pull-request and CodeQL snapshot after push; do
  not poll.
