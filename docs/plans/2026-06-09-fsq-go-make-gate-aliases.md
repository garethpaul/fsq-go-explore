# Foursquare Go Make Gate Aliases

status: completed

## Context

The repository had a working `make check` target, but the local pre-push gate
also expects `make lint`, `make test`, and `make build`. Without those aliases,
the first gate command fails before reaching Go formatting, tests, and static
guardrails.

## Objectives

- Provide stable Makefile targets for lint, test, build, and check.
- Keep the targets tied to the existing Go baseline.
- Document that lint, test, and build currently delegate to `make check`.
- Extend the static baseline so the aliases remain available.

## Verification

- `make lint`
- `make test`
- `make build`
- `make check`
- `git diff --check`
