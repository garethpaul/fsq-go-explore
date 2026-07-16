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

## Follow-up: guarantees now come from upstream go-ratelimiter

`limiter/` was a vendored copy of `github.com/garethpaul/go-ratelimiter` with no
test coverage, and it had drifted **in both directions**: upstream carried key
encoding, atomic multi-key consumption, rejection-status clamping and
`netip`-based IP canonicalization fixes the copy lacked, while the copy carried
this plan's refill work under names upstream does not use.

The application now depends on upstream. Upstream satisfies this plan's
guarantees — verified behaviorally rather than by matching source strings:

- Refill: `NewLimiter(2, 200ms)` allows 2, then allows 2 again after the TTL.
- Invalid configuration: `max=0`, `max=-1`, and `ttl=0` each allow **0 of 5**
  requests, so invalid configuration fails closed rather than becoming
  unlimited.

The previous `check-baseline.sh` greps asserted the vendored implementation's
internals (`float64(max) / ttl.Seconds()`, `max <= 0 || ttl <= 0`,
`func newTokenBucket`). Upstream implements the same behavior differently, so
those greps rejected a functionally equivalent limiter — the contract blocked
its own fix. They are replaced by `limiter_contract_test.go`, which asserts the
behavior through the public API, plus a pinned-dependency check and a guard
against re-vendoring `limiter/`.
