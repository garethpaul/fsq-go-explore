# CI Baseline

status: completed

## Context

The repository had a local Go `make check` baseline for formatting, tests,
App Engine import shape, cache keys, and OAuth/Foursquare request boundaries,
but no hosted workflow ran it for pushes and pull requests.

## Changes

- Added a GitHub Actions workflow that installs Go from `go.mod` and runs
  `make check`.
- Extended the baseline guard and docs so the hosted CI path stays visible.

## Verification

- `make check`
