# FSQ Rate Limiter Key Cap

status: completed

## Context

The local token-bucket limiter stores one bucket for every distinct request key
and never removes entries. Requests that rotate remote addresses, paths, or
configured key material can grow the map for the lifetime of the process.

## Priority

The limiter wraps the public search route, and its key dimensions include
request-controlled data. Bounding retained keys directly limits memory growth
without weakening per-key request limits.

## Implementation

- Add a default cap of 10,000 tracked keys.
- Maintain least-recently-used ordering with `container/list` for O(1) access,
  refresh, and eviction.
- Evict the least recently used key before inserting beyond the cap.
- Preserve existing token-bucket behavior for retained keys.
- Add package tests for the hard cap and recency-sensitive eviction.
- Extend the repository baseline and operational documentation.

## Verification

- `go test ./...`
- `go vet ./...`
- `make check`
- `make lint`
- `make test`
- `make build`
- `git diff --check`
- Mutations disabling the cap or recency refresh must fail.
- Hosted Go module integrity workflow.

## Follow-up: cap is now upstream's, asserted behaviorally

`limiter/` was an untested vendored copy of
`github.com/garethpaul/go-ratelimiter`. The application now depends on upstream,
which carries this cap (`defaultMaxTrackedKeys`, `list.New()`, `MoveToFront`,
`tokenBucketOrder`, LRU eviction) along with key-encoding, atomicity, IP
canonicalization and status-clamping fixes the copy had drifted away from.

The cap is asserted through the public API in `limiter_contract_test.go`
(`TestRateLimiterBoundsTrackedKeys`) rather than by grepping vendored source, so
an equivalent implementation is no longer rejected. `check-baseline.sh` also pins
the dependency and fails if `limiter/` is re-vendored, since the copy's zero
coverage is why the drift went unnoticed.
