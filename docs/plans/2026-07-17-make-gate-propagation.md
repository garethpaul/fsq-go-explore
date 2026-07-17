---
title: Make Gate Failure Propagation
type: reliability
date: 2026-07-17
status: completed
execution: code
---

# Make Gate Failure Propagation

## Summary

The baseline checker asserted that the Makefile *mentions*
`scripts/check-baseline.sh`, but a substring grep still matches a neutered
recipe. Appending `|| true` to the recipe left every pinned string
byte-identical while making `make check` — the sole CI entry point — exit 0
regardless of what the checker found. Observe the gate executing and propagating
its exit status instead of asserting that its name appears.

## Requirements

- R1. Pin the gate invocation as an anchored, whole recipe line, not a substring.
- R2. Prove the real Makefile propagates a failing gate's exit status for every
  gate target (`check`, `lint`, `test`, `build`).
- R3. Prove every gate target actually executes `scripts/check-baseline.sh`
  rather than echoing or skipping it.
- R4. Give the checker's verdict a path to CI that a neutered Makefile cannot
  swallow.
- R5. Pin that direct CI path so it cannot be silently dropped.
- R6. Preserve API, OAuth, request-body, response-body, status, timeout,
  dependency, workflow, and secret-handling behavior.

## Verification Plan

- Run `make check`, `make lint`, `make test`, and `make build` on a clean tree.
- Reject isolated hostile recipe-neuter mutations: `|| true`, `@echo` prefix,
  dash error-ignore prefix, `.IGNORE:`, `MAKEFLAGS += -i`, and a duplicate
  `check:` rule.
- Reject deletion of the direct CI invocation.
- Run Go formatting, vetting, tests, and module-integrity checks.

## Non-Goals

- Changing public API, OAuth flow, request/response limits, status handling,
  timeouts, dependencies, or workflow policy.
- Claiming credentialed Foursquare integration behavior.

## Work Completed

- Reused the existing `exact_line_count` helper (previously applied only to
  workflow contracts) to pin the gate recipe as a whole line, and hoisted its
  definition above first use.
- Added `verify_makefile_gate_propagation`: it copies the real Makefile into a
  throwaway directory, injects a passing stub that records execution and a
  failing stub that must turn every gate target red. This catches make-level
  neuters that leave the recipe byte-identical.
- Added a `Verify baseline gate wiring` CI step that runs the checker directly,
  and pinned that step's exact line, so a swallowed Makefile verdict still
  reaches CI.

## Verification Completed

- `make check`, `make lint`, `make test`, and `make build` passed on a clean
  tree; `go vet ./...`, `go test ./...`, `gofmt -l`, and `go mod tidy -diff`
  passed.
- Six isolated hostile Makefile neuters (`|| true`, `@echo` prefix, dash prefix,
  `.IGNORE:`, `MAKEFLAGS += -i`, duplicate `check:` rule) were each rejected by
  the checker; before this change all six left the checker green.
- Deleting the direct CI step alone was rejected by `make check`.
- Credentialed Foursquare integration behavior was unavailable and is not
  claimed.

## Known Limitation

`make check` cannot police the Makefile that defines it: each neuter above
destroys the exit-status channel the verdict travels on, so `make check` itself
stays green. Detection reaches CI only through the direct checker step. A
simultaneous Makefile neuter *and* deletion of that step would defeat both
layers; closing that requires an oracle outside the code the entry point
invokes. Any single-point neuter is caught.
