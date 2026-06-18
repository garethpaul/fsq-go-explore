---
title: Protobuf Security Refresh
type: security
status: planned
date: 2026-06-18
execution: code
---

# Protobuf Security Refresh

## Problem Frame

The maintained Go 1.25.11 toolchain removes reachable standard-library
vulnerabilities, but the required module graph still contains
`google.golang.org/protobuf` v1.26.0. A fresh `govulncheck -show verbose ./...`
reports GO-2024-2611 in that module, fixed in v1.33.0. The current application
does not call the vulnerable symbols, but retaining a known-vulnerable runtime
in the required graph weakens dependency assurance and leaves future call paths
exposed to an avoidable risk.

The deprecated `github.com/golang/protobuf` compatibility module is also at
v1.5.2. Its current terminal release, v1.5.4, requires the fixed modern
protobuf runtime, while the current `google.golang.org/protobuf` release is
v1.36.11.

## Prioritized Engineering Tasks

1. P0. Remove GO-2024-2611 from the required module graph by refreshing both
   protobuf compatibility layers together.
2. P0. Preserve App Engine memcache behavior, package compilation, race-test
   safety, and all existing OAuth and Foursquare response boundaries.
3. P1. Make the fixed dependency floor and zero-module-vulnerability result
   mutation-sensitive in the maintained baseline and documentation.
4. P2. Keep broader App Engine migration outside this change; replacing the
   legacy integration requires a separate product and deployment plan.

## Requirements

- R1. `github.com/golang/protobuf` must resolve to v1.5.4 and
  `google.golang.org/protobuf` must resolve to v1.36.11 without changing direct
  application dependencies or public behavior.
- R2. `go mod tidy` must leave a deterministic module graph and must not add an
  unrelated dependency or toolchain migration.
- R3. `govulncheck -show verbose ./...` must report zero symbol, package, and
  module vulnerabilities, including removal of GO-2024-2611 from the required
  graph.
- R4. Race-enabled tests, vet, build, every Make alias, and the external
  Makefile gate must pass under the module-selected Go 1.25.11 toolchain.
- R5. The baseline, security guidance, changelog, agent guidance, and this plan
  must record the dependency boundary and actual completed verification.
- R6. The pull request must retain the exact nonempty H2 sections Summary,
  Baseline, Improvements, Validation, Risk, and Follow-Ups in that order.

## Scope Boundaries

- Do not replace `google.golang.org/appengine`, change cache semantics, or
  migrate deployment infrastructure.
- Do not change OAuth, cookie, redirect, response-size, status, endpoint,
  media-type, duplicate-member, or UTF-8 validation behavior.
- Do not add credentials, live-provider calls, or network-dependent tests.
- Do not merge or close this stacked pull request without explicit repository
  owner authorization.

## Key Technical Decisions

- Upgrade both protobuf modules together. The legacy compatibility module and
  modern runtime form one transitive boundary through App Engine; refreshing
  only one makes the graph harder to reason about.
- Use the current compatible releases rather than only the minimum fixed
  version. v1.5.4 is the terminal compatibility release and v1.36.11 is the
  current modern runtime release observed by the Go module resolver.
- Keep the vulnerability gate module-aware. A symbol-only success is
  insufficient for this task because the defect is the presence of a known
  vulnerable required module even without a currently reachable call path.

## Implementation Units

### U1. Refresh the protobuf graph

- **Goal:** Remove the vulnerable protobuf runtime while preserving the direct
  dependency surface.
- **Requirements:** R1, R2
- **Files:** `go.mod`, `go.sum`
- **Approach:** Update the legacy compatibility module and modern protobuf
  runtime together, then tidy under the module-selected toolchain.
- **Patterns to follow:** The focused Go toolchain refresh in
  `docs/plans/2026-06-18-go-1-25-11-security-refresh.md`.
- **Test scenarios:** The graph resolves the two intended versions, direct
  dependencies remain unchanged, and a second tidy produces no diff.
- **Verification:** Inspect the exact graph and sums before running behavioral
  gates.

### U2. Preserve runtime behavior and concurrency safety

- **Goal:** Prove the dependency refresh does not change application behavior.
- **Requirements:** R4
- **Files:** `auth_test.go`, `cache_test.go`, `edit_test.go`, `fsq/api_test.go`,
  `fsq/keys_test.go`, `search_test.go`
- **Approach:** Reuse the maintained tests without changing assertions unless
  the refreshed modules expose a real compatibility defect.
- **Test scenarios:** Race-enabled package tests, vet, build, normal package
  tests, and the static response contracts all pass under Go 1.25.11.
- **Verification:** Every repository Make alias and the external Makefile gate
  pass from clean state.

### U3. Make security evidence durable

- **Goal:** Prevent the vulnerable protobuf graph or incomplete evidence from
  returning unnoticed.
- **Requirements:** R3, R5, R6
- **Files:** `scripts/check-baseline.sh`, `AGENTS.md`, `README.md`, `SECURITY.md`,
  `CHANGES.md`, `docs/plans/2026-06-18-protobuf-security-refresh.md`
- **Approach:** Add precise version and completed-evidence contracts, document
  the module-level advisory boundary, and record exact local and hosted results
  only after they exist.
- **Test scenarios:** Mutations restoring either old protobuf version,
  removing zero-module-vulnerability evidence, or reopening the plan status are
  rejected with contract-specific messages.
- **Verification:** Diff, artifact, secret, module-graph, mutation, and hosted
  exact-head audits all pass.

## Verification Plan

- Confirm the selected Go runtime remains 1.25.11 with `GOTOOLCHAIN=auto`.
- Run race-enabled package tests, vet, build, tidy-diff, all Make aliases, and
  the external Makefile gate.
- Run `govulncheck -show verbose ./...` and require zero symbol, package, and
  module vulnerabilities.
- Reject isolated dependency and evidence mutations in initialized temporary
  Git repositories without modifying the live worktree.
- Audit intended paths, generated artifacts, credential-shaped additions, file
  modes, and whitespace before committing explicit paths.
- Require terminal canonical push and pull-request checks on the exact final
  head before terminal tracker reconciliation.

## Risks

- App Engine remains a legacy dependency and may constrain future dependency
  modernization even after this protobuf refresh.
- A newer protobuf runtime can expose latent generated-code incompatibilities;
  full package compilation and race tests are the primary local guard.
- Live Foursquare, OAuth, and App Engine behavior remain outside the offline
  verification boundary.

## Sources

- Go vulnerability database entry GO-2024-2611:
  <https://pkg.go.dev/vuln/GO-2024-2611>
- Module metadata resolved by `go list -m -u` and `go mod download` for
  `github.com/golang/protobuf` v1.5.4 and
  `google.golang.org/protobuf` v1.36.11.
