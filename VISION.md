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
- `scripts/check-baseline.sh`, `make lint`, `make test`, `make build`, and
  `make check` run `go test ./...`, Go formatting checks, module-import checks,
  and credential/privacy-log guardrails.
- Internal imports use `github.com/garethpaul/fsq-go-explore/...` module paths.
- The cache-key generation path uses stable SHA-256 digests rather than reversible
  encrypted payloads.
- Foursquare request logging avoids raw URLs, tokens, edit payloads, and user
  records, and search handling avoids raw location logging.
- OAuth redirects use per-login state cookies rather than a shared static state
  string.
- OAuth callbacks reject missing OAuth authorization codes before token exchange
  work starts.
- Auth cookies validate generated user cache keys before access-token memcache
  lookup starts.
- Protected routes validate generated auth cookie cache keys before handler
  work starts.
- ETag comparisons are exact before header-cache `304 Not Modified` responses.
- Venue edit submissions reject non-POST requests before auth or API work.
- Venue detail and edit handlers reject missing venue IDs before auth, template,
  or Foursquare API work where possible.
- Propose-edit rejects malformed edit forms before auth-cookie lookup or
  Foursquare edit API work.
- Search query and location parameters are trimmed and length-bounded before
  venue search requests.
- The local Makefile exposes lint, test, build, and check targets for a stable
  pre-push gate.

Next priorities:

- Modernize App Engine Go runtime assumptions in a dedicated pass
- Clarify secret handling for local and hosted environments
- Keep state-changing handlers method-constrained and covered by tests
- Keep missing OAuth authorization codes covered before exchange work is added
- Keep user cache keys validated before memcache lookup
- Keep protected-route auth cookie validation covered before handler work starts
- Keep ETag matching exact when changing header-cache behavior
- Keep missing venue IDs and malformed request boundaries covered by tests before
  auth or API side effects are introduced
- Keep malformed edit forms rejected before auth-cookie lookup and Foursquare
  edit API calls
- Keep search parameter bounds covered as request parsing changes
- Keep local verification targets available as the Go/App Engine toolchain
  evolves
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
