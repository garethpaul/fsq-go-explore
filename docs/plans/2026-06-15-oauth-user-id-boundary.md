---
title: OAuth User ID Boundary
type: security
status: completed
date: 2026-06-15
execution: code
---

# OAuth User ID Boundary

## Problem Frame

The OAuth callback accepts any structurally decodable user profile. A response
with a missing, empty, whitespace-only, or edge-whitespace user ID therefore
reaches cache-key creation, stores the access token, and publishes an
authenticated cookie for an invalid identity.

## Prioritized Engineering Work

1. **P0 - Authentication identity:** reject invalid user IDs before the decoded
   profile can reach cache or cookie publication.
2. **P1 - Regression coverage:** prove missing and whitespace-padded IDs fail
   while a canonical ID and optional empty display name remain accepted.
3. **P2 - Maintained contract:** enforce validation ordering, generic error
   handling, synchronized guidance, and completed verification evidence.

## Requirements

- Require a nonempty Foursquare user ID after JSON decoding.
- Reject leading or trailing Unicode whitespace rather than silently
  normalizing the provider identity.
- Preserve optional first-name behavior and the canonical user ID unchanged.
- Keep status, final origin, JSON media type, 1 MiB body bound, redirect
  refusal, and 10-second request timeout behavior unchanged.
- Ensure invalid identities fail before `FoursquareUser`, cache-key, memcache,
  or authentication-cookie construction.
- Static contracts must reject validation removal, whitespace weakening,
  post-publication ordering, regression loss, guidance drift, and incomplete
  plan evidence.

## Implementation Units

### Validate Decoded OAuth Identity

Files:

- `auth.go`
- `auth_test.go`

Approach:

- Add one deterministic invalid-identity error.
- Validate the decoded ID before returning the profile from
  `decodeOAuthUserResponse`.
- Add focused table-driven invalid-ID coverage and a valid empty-name case.

### Enforce And Document The Boundary

Files:

- `scripts/check-baseline.sh`
- `AGENTS.md`
- `CHANGES.md`
- `README.md`
- `SECURITY.md`
- `VISION.md`
- `docs/plans/2026-06-15-oauth-user-id-boundary.md`

## Verification

- Focused OAuth response tests and the full Go test/race/vet/static gate.
- Repository and external-directory Make gates.
- Hostile mutations for validation removal, whitespace weakening, ordering,
  missing-ID regression, edge-whitespace regression, valid empty-name
  behavior, guidance, and plan completion evidence.
- Exact diff, generated artifact, conflict marker, executable mode, and
  changed-line credential audits.

## Scope Boundaries

- Do not change OAuth state, token exchange, response transport, cache-key
  hashing, cookie attributes, display-name requirements, or dependencies.
- Do not execute live Foursquare credentials, OAuth callbacks, App Engine, or
  memcache flows.

## Verification Completed

- Focused OAuth identity tests passed for five invalid ID forms and the valid
  empty-display-name case.
- Repository and external-directory Make gates passed with the full baseline.
- Eight hostile mutations failed across validation removal, whitespace
  weakening, publication ordering, regression coverage, guidance, and plan
  evidence.
- `go test ./...`, `go test -race ./...`, and `go vet ./...` passed.
- Diff, artifact, conflict-marker, mode, and changed-line credential audits
  passed.
- No live OAuth callback was executed because provider credentials and App
  Engine services are outside this local verification boundary.
