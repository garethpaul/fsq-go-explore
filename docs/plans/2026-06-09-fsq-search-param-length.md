# FSQ Search Parameter Length

status: completed

## Context

Search query and location values are forwarded into Foursquare venue search
parameters. They were trimmed, but unbounded values from request parameters or
App Engine location headers could still be carried through to upstream API
requests and cache-key inputs.

## Completed Scope

- Added a shared search parameter normalizer.
- Capped search query, explicit location, and App Engine location fallback
  values to a bounded rune length.
- Added parser coverage for long query and location values.
- Extended the static baseline and docs to preserve the length boundary.

## Verification

- `go test ./...`
- `make check`
- `git diff --check`
