---
title: Refresh the Go toolchain to 1.25.11
type: security
date: 2026-06-18
status: completed
execution: code
---

# Refresh the Go toolchain to 1.25.11

## Priority

P1 standard-library security. `govulncheck` reports 18 reachable
vulnerabilities when the repository is built with its declared Go 1.25.3
toolchain. The affected paths include HTML template escaping, TLS and X.509
verification, HTTP/2 transport, URL parsing, form parsing, and MIME-header
handling.

## Current-State Findings

- `go.mod` declares Go 1.25.3, and GitHub Actions installs that exact release
  through `go-version-file`.
- The highest minimum patched release reported across the reachable findings is
  Go 1.25.11.
- All direct modules are already at their current releases; the reported
  reachable findings are in the standard library rather than direct module
  dependencies.
- The current branch passes race-enabled tests, vet, module tidiness, and the
  complete repository gate before the toolchain change.

## Approach

- Raise the `go.mod` language/toolchain declaration from 1.25.3 to 1.25.11.
- Keep GitHub Actions sourcing its exact version from `go.mod` so push and pull
  request checks use the patched toolchain.
- Extend the baseline contract to reject restoration of vulnerable Go patch
  releases and require completed security evidence in this plan.
- Preserve application behavior, module versions, response boundaries, and the
  existing CI event coverage.

## Implementation Units

### U1: Raise the patched toolchain floor

Update `go.mod` to Go 1.25.11 without changing module requirements or sums.

Verification:
- Go 1.25.11 runs all race-enabled package tests.
- `go vet`, `go mod tidy -diff`, and all Make aliases pass.
- `govulncheck` reports no reachable vulnerabilities.

### U2: Enforce the security boundary

Update the maintained baseline so the declared Go version and completed plan
evidence fail closed if weakened.

Verification:
- Repository-root and external-directory `make check` pass under Go 1.25.11.
- Isolated mutations to the Go patch floor, plan status, and vulnerability
  evidence are rejected.

## Scope Boundaries

- Do not change direct or transitive module versions in this toolchain refresh.
- Do not change Foursquare/OAuth request behavior, response parsing, templates,
  rate limiting, or public interfaces.
- Do not claim live provider validation; all network behavior remains covered
  by offline fixtures and deterministic policy tests.
- Keep PR #20 and its predecessors open and preserve base-first ordering.

## Success Criteria

- `go.mod` requires Go 1.25.11.
- Race-enabled tests, vet, module tidiness, Make gates, and `govulncheck` pass
  with the patched toolchain.
- Push and pull-request checks pass on the exact implementation head.
- The plan retains truthful local and hosted verification evidence.

## Verification Completed

- Go 1.25.11 ran `go test -race -count=1 ./...`, `go vet ./...`, and
  `go mod tidy -diff` successfully without changing module requirements or
  sums.
- `govulncheck` reported no reachable vulnerabilities with the patched
  toolchain; the pre-change Go 1.25.3 scan reported 18 reachable
  standard-library vulnerabilities.
- Repository-root and external-directory `make check` passed under Go 1.25.11.
- Four isolated hostile mutations were rejected across the declared toolchain
  floor, actual runtime floor, completed plan status, and no-vulnerability
  evidence.
