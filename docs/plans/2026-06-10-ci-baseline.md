# CI Baseline

status: completed

## Context

The repository had a local Go `make check` baseline for formatting, tests,
App Engine import shape, cache keys, and OAuth/Foursquare request boundaries,
but no hosted workflow ran it for pushes and pull requests.

## Changes

- Added a least-privilege GitHub Actions workflow that installs the exact Go
  version from `go.mod` and runs `make check`.
- Pinned checkout and setup-go by commit, cancelled superseded runs, and
  bounded the job with a timeout.
- Disabled checkout credential persistence and kept the workflow token limited
  to read-only repository contents.
- Extended the local gate to run formatting, `go vet`, `go test ./...`, and
  `go mod tidy -diff`.
- Extended the baseline guard and docs so the hosted CI path stays visible and
  rejects duplicate checkout steps, nested or write-capable permissions, and
  weakened action or command contracts.

## Verification

- `make check`
- Workflow YAML parse
- Hostile workflow mutation checks
- Hosted Go 1.25.3 GitHub Actions run
