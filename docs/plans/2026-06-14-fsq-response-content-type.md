---
title: Foursquare Response Content Type
type: security
status: completed
date: 2026-06-14
---

# Foursquare Response Content Type

## Summary

Require Foursquare search and venue-detail responses to declare a JSON media
type before bounded body decoding. Accept standard and structured application
JSON types while rejecting missing, malformed, text, and binary declarations.

## Requirements

- R1. Validate media type after status and before body decoding.
- R2. Accept `application/json` and `application/*+json` case-insensitively.
- R3. Reject missing, malformed, or incompatible content types with empty results.
- R4. Preserve timeouts, status checks, body limits, and request behavior.

## Verification

- Focused Go tests covered accepted and rejected media types plus search/detail
  fail-closed behavior.
- Seven hostile mutations were rejected across parsing, standard JSON,
  structured JSON, search use, details use, test coverage, and completed-plan
  evidence.
- `make check`, `make lint`, `make test`, and `make build` passed from the root;
  `make check` also passed through the absolute Makefile path from `/tmp`.
- No live Foursquare, OAuth, cache, or App Engine request was made.
- Exact intended-path, artifact, whitespace, conflict-marker, dependency,
  workflow, and changed-line credential-pattern audits passed.
