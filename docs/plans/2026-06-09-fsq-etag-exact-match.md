---
title: FSQ ETag Exact Match
date: 2026-06-09
status: completed
execution: code
---

## Context

`HeaderCache` compared `If-None-Match` with `strings.Contains`, so a header that
only contained the generated ETag as a substring could return `304 Not
Modified`. Cache validators should compare full candidate tags, including
comma-separated header values.

## Goals

- Replace substring ETag matching with exact candidate matching.
- Keep unquoted and quoted ETag candidates compatible with the existing server
  header value.
- Add tests for partial-match rejection and comma-separated exact matches.
- Extend static verification and docs for the header-cache boundary.

## Implementation

- Added `etagMatches` to split `If-None-Match` on commas, trim whitespace, and
  compare unquoted candidates exactly.
- Added `cache_test.go` coverage for partial and exact ETag behavior.
- Extended `scripts/check-baseline.sh`, README, SECURITY, VISION, and CHANGES.

## Verification

- `go test ./...`
- `scripts/check-baseline.sh`
- `make check`
- `make lint`
- `make test`
- `make build`
- `git diff --check`
