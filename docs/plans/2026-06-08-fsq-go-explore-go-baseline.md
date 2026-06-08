# Foursquare Go Explore Go Baseline Plan

status: completed

## Context

`fsq-go-explore` is a legacy Go/App Engine Foursquare sample. The current source uses GOPATH-era local imports and mixed App Engine import paths, so modern Go cannot compile or test the repository without additional dependency and import metadata.

## Objectives

- Add a reproducible Go module baseline without changing the sample's runtime behavior.
- Update local imports and App Engine imports so `go test ./...` can compile the repository.
- Add focused tests around cache-key generation and search request parsing.
- Replace credential-adjacent logging with safer diagnostics.
- Document the supported verification command and dependency expectations.

## Work Items

1. Added `go.mod` and `go.sum` for the existing packages.
2. Updated internal imports to use the module path and modern App Engine package paths.
3. Added `make check` with Go formatting, tests, and static credential checks.
4. Added focused tests for deterministic cache keys, OAuth state generation, and App Engine header fallback parsing.
5. Updated README, VISION, SECURITY, and CHANGES with the new baseline.
6. Replaced reversible cache-key encryption with stable SHA-256 digest keys.
7. Replaced static OAuth state with a per-login state cookie.
8. Removed raw token, URL, edit payload, and user record logging.

## Verification

- `make check`
- `go test ./...`
- `git diff --check`
