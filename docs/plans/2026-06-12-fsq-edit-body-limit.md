# FSQ Venue Edit Request Body Limit

status: planned

## Context

`ProposeEdit` parses the complete URL-encoded request body before access-token
lookup. The route wrapper rejects absent or malformed auth cookies, but any
caller can supply a syntactically valid generated-key shape and reach form
parsing without a live session.

Go applies a broad default limit to URL-encoded forms. The venue edit sample
needs a smaller application-level boundary because its expected fields are
short text values and the parsed form is forwarded to Foursquare.

## Priority

This is a public request boundary that can consume substantially more memory
and request-processing time than the feature requires. A focused cap reduces
denial-of-service exposure without changing valid edit behavior or the
Foursquare API contract.

## Prioritized Backlog

1. Limit venue edit request bodies to 64 KiB before parsing form data.
2. Return `413 Request Entity Too Large` when the cap is exceeded, while
   retaining `400 Bad Request` for other malformed forms.
3. Cover exact-boundary and oversized requests before auth or Foursquare work.
4. Record the boundary in repository security and maintenance documentation.

## Implementation

- Wrap `ProposeEdit` request bodies with `http.MaxBytesReader` before
  `ParseForm`.
- Detect `*http.MaxBytesError` with `errors.As` and use an explicit response
  that distinguishes oversized bodies from malformed encoding.
- Add focused handler tests using URL-encoded bodies and no valid session.
- Extend README, SECURITY, VISION, CHANGES, and the static baseline contract.

## Verification

- `gofmt -w edit.go edit_test.go`
- `go test -count=1 ./...`
- `go test -race -count=1 ./...`
- `go vet ./...`
- `go mod tidy -diff`
- `make check`
- `git diff --check`
- A mutation that removes `http.MaxBytesReader` or collapses oversized bodies
  into `400 Bad Request` must fail.
