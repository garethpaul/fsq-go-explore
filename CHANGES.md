# Changes

## 2026-06-08

- Added a Go module and lockfile for the legacy App Engine sample.
- Updated local and App Engine imports so `go test ./...` compiles under modules.
- Replaced reversible cache-key construction with deterministic SHA-256 keys.
- Replaced static OAuth state with a per-login state cookie.
- Escaped venue IDs before building Foursquare detail and edit request paths.
- Returned explicit HTTP errors for search cache failure paths.
- Removed credential- and location-adjacent logging from OAuth, search, and Foursquare API flows.
- Added tests for cache-key behavior, OAuth state generation, venue path escaping, and App Engine location fallback parsing.
- Added `make check` and `scripts/check-baseline.sh` for formatting, tests, and static guardrails.
