# Changes

- Bound the in-process rate limiter to 10,000 tracked keys with
  least-recently-used eviction and focused capacity tests.

## 2026-06-10

- Added a pinned, least-privilege GitHub Actions workflow for the exact Go
  toolchain, formatting, vet, tests, module integrity, and static checks.

## 2026-06-09

- Rejected malformed venue edit forms before auth or Foursquare edit work.
- Rejected malformed auth-cookie cache keys in the protected-route wrapper
  before handler work starts.
- Made header-cache ETag comparisons exact so partial `If-None-Match` values
  cannot trigger `304 Not Modified` responses.

## 2026-06-08

- Added `make lint`, `make test`, and `make build` aliases so local verification
  has the expected pre-push gate targets in addition to `make check`.
- Added a Go module and lockfile for the legacy App Engine sample.
- Updated local and App Engine imports so `go test ./...` compiles under modules.
- Replaced reversible cache-key construction with deterministic SHA-256 keys.
- Replaced static OAuth state with a per-login state cookie.
- Escaped venue IDs before building Foursquare detail and edit request paths.
- Returned explicit HTTP errors for search cache failure paths.
- Restricted the venue edit submission handler to POST requests.
- Rejected missing venue IDs before Foursquare venue detail or edit API work.
- Rejected missing edit-page venue IDs before auth and template work.
- Bounded search query and location parameters before Foursquare venue search
  requests are built.
- Rejected OAuth callbacks with missing authorization codes before token exchange
  work starts.
- Validated auth-cookie user cache keys before access-token memcache lookup.
- Removed credential- and location-adjacent logging from OAuth, search, and Foursquare API flows.
- Added tests for cache-key behavior, OAuth state generation, venue path escaping, and App Engine location fallback parsing.
- Added `make check` and `scripts/check-baseline.sh` for formatting, tests, and static guardrails.
