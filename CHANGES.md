# Changes

## 2026-06-13

- Made Go verification independent of the caller's working directory by
  resolving the baseline checker from the loaded Makefile.
- Rejected non-2xx and over-1-MiB OAuth user-profile responses before session
  state is created.
- Added a 10-second default Foursquare HTTP client timeout while preserving
  explicit positive caller values and caller configuration ownership.
- Added focused constructor tests and mutation-sensitive timeout contracts.
- Rejected non-2xx Foursquare search and venue detail responses before JSON decoding,
  preventing error envelopes from populating successful result structures.
- Added focused transport tests and method-scoped static ordering contracts.
- Bounded Foursquare JSON response parsing to 2 MiB before unmarshalling and
  added exact-limit, oversized-body, and reader-error tests.

## 2026-06-12

- Limited venue edit request bodies to 64 KiB and return `413 Request Entity
  Too Large` before auth or Foursquare work when that boundary is exceeded.
- Corrected token-bucket refill rates so `Max` requests are restored over each
  `TTL`, while invalid non-positive configurations fail closed.
- Added deterministic limiter tests for burst, sustained refill, and invalid
  configuration behavior.
- Bound the in-process rate limiter to 10,000 tracked keys with
  least-recently-used eviction and focused capacity tests.

## 2026-06-10

- Added a pinned, least-privilege GitHub Actions workflow for the exact Go
  toolchain, formatting, vet, tests, module integrity, and static checks, with
  credential-free checkout.

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
