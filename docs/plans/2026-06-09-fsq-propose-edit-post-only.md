# FSQ Propose Edit POST-only Handler

status: completed

## Context

`ProposeEdit` changes venue data through the Foursquare API, but non-POST
requests reached the authentication path instead of being rejected at the method
boundary.

## Completed Scope

- Added handler coverage for non-POST `/propose_edit` requests.
- Returned `405 Method Not Allowed` with `Allow: POST` before auth or API work.
- Extended the static baseline and docs for the POST-only edit contract.

## Verification

- `go test ./...`
- `make check`
- `git diff --check`
