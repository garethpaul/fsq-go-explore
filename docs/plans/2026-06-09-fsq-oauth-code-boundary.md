# Foursquare Go Explore OAuth Code Boundary

status: completed

## Context

OAuth callback handling already validates a per-login state cookie before
exchanging the authorization code. A callback with matching state but no `code`
parameter should still stop before App Engine context setup or token exchange
work.

## Objectives

- Reject callbacks that do not include a non-blank authorization code.
- Clear the transient OAuth state cookie on the rejected callback path.
- Add unit coverage that exercises the missing-code boundary without invoking
  token exchange.
- Extend the static baseline and docs so the callback boundary stays visible.

## Verification

- `go test ./...`
- `scripts/check-baseline.sh`
- `make check`
- `git diff --check`
