---
title: "fix: Restrict OAuth callbacks to GET"
type: fix
date: 2026-06-26
status: completed
---

# fix: Restrict OAuth callbacks to GET

## Context

The redirect handler called `FormValue` before checking the request method.
An unauthenticated POST could therefore trigger form or multipart body parsing
even though Foursquare documents the legacy v2 web authorization response as a
redirect URI carrying `code` in the URL.

Primary reference:

- `https://docs.foursquare.com/developer/reference/v2-authentication`

## Decision

Accept only `GET`. Return `405 Method Not Allowed` with `Allow: GET` before
logging, request-body reads, state validation, token exchange, or user-profile
requests.

## Alternatives

- Add a POST body size cap: rejected because POST is not part of the provider's
  documented callback contract and still expands the accepted surface.
- Parse query parameters directly but allow every method: rejected because it
  leaves ambiguous and unnecessary callback semantics.
- Apply the guard in routing middleware: rejected because the handler is also
  called directly by tests and should own its protocol boundary.

## Verification Completed

- The focused regression failed with HTTP 307 before implementation.
- The regression passes after the early method guard and proves zero body reads.
- Go 1.25.11 was downloaded from `go.dev` and its published SHA-256 verified.
- `go test -race -count=1 ./...`, `go vet ./...`, `go mod tidy -diff`, all
  Make aliases, an external-directory Make gate, and `git diff --check` passed.
- An isolated guard-removal mutation was rejected with the intended 307/405
  failure.

## Boundary

No live provider callback, credentials, token exchange, App Engine memcache,
or Foursquare API request is exercised locally.
