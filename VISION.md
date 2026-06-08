## Foursquare Go Explore Vision

This document explains the current state and direction of the project.
Project overview and developer docs: [`README.md`](README.md)

Foursquare Go Explore is a Go/App Engine sample for calling Foursquare venue
search APIs with header and memcache-based caching.

The repository is useful as a practical API integration sample, with a reusable
rate-limiter package and detailed caching guidance. Project setup notes live in
[`README.md`](README.md).

The goal is to keep the sample secure, cache-aware, and easy to run with local
Foursquare credentials.

The current focus is:

Priority:

- Preserve the Foursquare search flow and App Engine deployment shape
- Keep `FSQ_CLIENT_ID`, `FSQ_CLIENT_SECRET`, and `FSQ_VERSION` in environment
  configuration
- Maintain caching behavior that respects platform policy
- Keep the embedded limiter package understandable and documented

Current baseline:

- The Go module is defined by `go.mod` and `go.sum` with App Engine/OAuth/rate-limit
  dependencies.
- `scripts/check-baseline.sh` and `make check` run `go test ./...`, Go
  formatting checks, module-import checks, and credential/privacy-log guardrails.
- Internal imports use `github.com/garethpaul/fsq-go-explore/...` module paths.
- The cache-key generation path uses stable SHA-256 digests rather than reversible
  encrypted payloads.
- Foursquare request logging avoids raw URLs, tokens, edit payloads, and user
  records, and search handling avoids raw location logging.
- OAuth redirects use per-login state cookies rather than a shared static state
  string.

Next priorities:

- Modernize App Engine Go runtime assumptions in a dedicated pass
- Clarify secret handling for local and hosted environments
- Separate reusable limiter concerns from demo-specific API code if needed

Contribution rules:

- One PR = one focused API, cache, limiter, or documentation change.
- Do not commit real Foursquare credentials.
- Verify local serving or tests before pushing Go changes.
- Update caching docs when request keys, TTLs, or headers change.

## Security

Canonical security policy and reporting:

- [`SECURITY.md`](SECURITY.md)

Foursquare credentials are private and must stay out of source control. Cache
keys should not expose raw secrets or unnecessary user data.

Rate limiting and caching should fail predictably under abuse or upstream
errors.

## What We Will Not Merge (For Now)

- Hardcoded client IDs or secrets
- Cache changes that violate platform policy or leak sensitive data
- Broad runtime migrations bundled with API behavior changes
- Limiter changes without examples or tests

This list is a roadmap guardrail, not a permanent rule.
Strong user demand and strong technical rationale can change it.
