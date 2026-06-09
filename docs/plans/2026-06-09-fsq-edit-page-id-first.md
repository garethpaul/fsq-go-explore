# FSQ Edit Page ID-First Boundary

status: completed

## Context

The edit page eventually rejected missing venue IDs, but it performed auth
checks and parsed templates first. Malformed requests should fail at the request
boundary before invoking unrelated work or redirect behavior.

## Completed Scope

- Moved edit-page venue ID trimming and empty checks ahead of auth and template
  parsing.
- Added handler coverage showing missing or blank edit-page IDs return
  `400 Bad Request` without redirecting to login.
- Extended the static baseline and docs to preserve the ID-first request
  boundary.

## Verification

- `go test ./...`
- `make check`
- `git diff --check`
