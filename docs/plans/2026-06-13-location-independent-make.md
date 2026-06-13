---
title: Location-Independent Go Verification
type: reliability
date: 2026-06-13
status: planned
execution: code
---

# Location-Independent Go Verification

## Summary

Resolve the maintained checker from the loaded Makefile and ensure its Go
commands run from the repository root so every documented gate works outside
the checkout.

## Requirements

- R1. Derive the repository root from `MAKEFILE_LIST`.
- R2. Invoke `scripts/check-baseline.sh` through its repository-rooted path.
- R3. Run repository-relative Go commands from the repository root if the
  checker currently depends on the caller's directory.
- R4. Add mutation-sensitive contracts for every corrected location boundary.
- R5. Preserve API, OAuth, request-body, response-body, status, timeout,
  dependency, workflow, and secret-handling behavior.
- R6. Record actual root and external-directory verification before completion.

## Verification Plan

- Run `make check`, `make lint`, `make test`, and `make build` at repository
  root and from `/tmp` through the absolute Makefile path.
- Reject isolated hostile Make-root, checker-path, Go-working-directory,
  documentation, plan-status, and verification-evidence mutations as applicable.
- Run Go formatting/vetting/testing, shell syntax, `git diff --check`, exact-path
  review, secret scanning, and generated-artifact inspection.

## Non-Goals

- Changing public API, OAuth flow, request/response limits, status handling,
  timeouts, dependencies, or workflow policy.
- Claiming credentialed Foursquare integration behavior.

## Work Completed

Pending implementation.

## Verification Completed

Pending implementation and verification.
