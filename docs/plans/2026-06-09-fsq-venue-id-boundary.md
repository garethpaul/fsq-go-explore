# FSQ Venue ID Boundary

status: completed

## Context

Venue detail and edit flows build Foursquare venue URLs from the request `id`
parameter. The URL builder path-escapes IDs, but an empty or whitespace-only ID
still produced a request path for no concrete venue.

## Completed Scope

- Trimmed venue IDs in edit-page and propose-edit handlers.
- Rejected missing venue IDs with `400 Bad Request` before outbound Foursquare
  detail or edit API work.
- Added handler coverage for POST `/propose_edit` without an ID.
- Extended the static baseline and docs to preserve the venue ID boundary.

## Verification

- `go test ./...`
- `make check`
- `git diff --check`
