# fsq-go-explore

<!-- README-OVERVIEW-IMAGE -->
![Project overview](docs/readme-overview.svg)

## Overview

`garethpaul/fsq-go-explore` is a Go project. Foursquare GoLang Explore Sample

This README is based on the checked-in source, manifests, scripts, and repository metadata on the `master` branch. The project language mix found during review was: Go (14).

## Repository Contents

- `README.md` - project overview and local usage notes
- `CHANGES.md` - concise history of maintenance changes
- `Makefile` - local verification entry point
- `go.mod` and `go.sum` - Go module dependency metadata
- `fsq` - source or example code
- `limiter` - source or example code
- `scripts/check-baseline.sh` - Go formatting, test, import, and credential/privacy-log checks
- `SECURITY.md` - security reporting and disclosure guidance
- `static` - source or example code
- `templates` - source or example code
- `VISION.md` - project direction and maintenance guardrails

Additional scan context:

- Source directories: fsq, limiter, static, templates
- Dependency and build manifests: go.mod, go.sum
- Entry points or build surfaces: main.go, app.yaml
- Test-looking files: auth_test.go, search_test.go, fsq/api_test.go, fsq/keys_test.go

## Getting Started

### Prerequisites

- Git
- Go 1.25 or a compatible modern Go toolchain

### Setup

```bash
git clone https://github.com/garethpaul/fsq-go-explore.git
cd fsq-go-explore
go mod download
```

The setup commands above are derived from repository files. Legacy mobile, Python, or JavaScript samples may require older SDKs or package versions than a modern workstation uses by default.

## Running or Using the Project

- This is a legacy Google App Engine sample. Use the App Engine tooling that
  matches `app.yaml` for local serving or deployment.
- Configure `FSQ_CLIENT_ID`, `FSQ_CLIENT_SECRET`, and `FSQ_VERSION` through the
  environment or deployment configuration. Do not commit real values.

## Testing and Verification

Run the baseline:

```bash
make check
```

The baseline runs `go test ./...`, verifies Go formatting, checks that module
imports are used instead of GOPATH-era local imports, and guards against
credential- and location-adjacent logging. It also covers state-changing venue
edit submissions so non-POST requests are rejected before auth or Foursquare API
work, and missing venue IDs are rejected before venue API requests are built.
Search query and location values are trimmed and length-bounded before venue
search requests are built.

When the required SDK or runtime is unavailable, use static checks and source review first, then verify on a machine that has the matching platform toolchain.

## Configuration and Secrets

- Required Foursquare settings: `FSQ_CLIENT_ID`, `FSQ_CLIENT_SECRET`, and
  `FSQ_VERSION`.
- Keep API keys, OAuth credentials, access tokens, `.env` files, and
  deployment-specific config out of source control.

## Security and Privacy Notes

- Review changes touching authentication or token handling; examples from the scan include auth.go, limiter/config/config.go, limiter/limiter.go, main.go.
- Review changes touching external API calls or credential-adjacent configuration; examples from the scan include auth.go, edit.go, fsq/api.go, fsq/common.go, and 6 more.
- Review changes touching network requests, sockets, or service endpoints; examples from the scan include auth.go, fsq/api.go, fsq/common.go, fsq/venue.go, and 3 more.
- Review changes touching file, media, JSON, XML, CSV, OCR, or data parsing; examples from the scan include auth.go, fsq/api.go, fsq/common.go, fsq/keys.go, and 3 more.
- Cache keys are deterministic SHA-256 digests and should not expose raw query,
  token, or user fields.
- OAuth login uses per-request state values and HTTP-only cookies for callback
  validation.
- Venue edit submissions are POST-only; non-POST requests receive `405 Method
  Not Allowed`.
- Missing venue IDs are rejected with `400 Bad Request` before Foursquare venue
  detail or edit API work.
- Search query and location parameters are length-bounded before being sent to
  Foursquare or used in cache keys.

## Maintenance Notes

- Run `make check` before pushing changes that touch Foursquare API calls,
  OAuth, cache keys, rate limiting, or App Engine imports.
- See `SECURITY.md` for vulnerability reporting and safe research guidance.
- See `VISION.md` for project direction and contribution guardrails.
- See `docs/plans/2026-06-09-fsq-propose-edit-post-only.md` for the venue edit
  method guard.
- See `docs/plans/2026-06-09-fsq-venue-id-boundary.md` for the venue ID request
  boundary.
- See `docs/plans/2026-06-09-fsq-search-param-length.md` for search parameter
  length guardrails.

## Contributing

Keep changes small and tied to the project that is already present in this repository. For code changes, document the toolchain used, avoid committing generated dependency directories or local configuration, and update this README when setup or verification steps change.
