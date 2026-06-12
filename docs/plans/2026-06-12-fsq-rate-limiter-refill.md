# FSQ Rate Limiter Refill Semantics

status: completed

## Context

The reusable limiter constructs each token bucket with a refill interval equal
to the complete configured TTL. `NewLimiter(10, time.Minute)` therefore permits
an initial burst of ten requests but restores only one request per minute,
rather than the documented maximum of ten requests per minute.

## Priority

The limiter wraps the public search route. An incorrect sustained rate can
reject ordinary traffic long after a short burst and makes the `Max` and `TTL`
configuration contract misleading.

## Prioritized Backlog

1. Refill `Max` tokens over each `TTL` and cover the configured rate now.
2. Keep the existing per-key burst and least-recently-used key cap intact.
3. Review trusted proxy address configuration separately if the deployment
   topology changes.

## Implementation

- Construct token buckets at `Max / TTL` requests per second with a burst of
  `Max`.
- Reject non-positive maximums or durations by creating a closed bucket.
- Add deterministic tests for the initial burst, refill rate, and invalid
  configuration without wall-clock sleeps.
- Extend the repository baseline and operational documentation.

## Verification

- `gofmt -w limiter/config/config.go limiter/config/config_test.go`
- `go test ./...`
- `go vet ./...`
- `make lint`
- `make test`
- `make build`
- `make check`
- `git diff --check`
- Mutations restoring one-token-per-TTL refill or allowing invalid
  configurations must fail.
