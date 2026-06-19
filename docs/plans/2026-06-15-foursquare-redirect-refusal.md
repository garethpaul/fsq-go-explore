---
title: Foursquare Redirect Refusal
type: security
status: completed
date: 2026-06-15
---

# Foursquare Redirect Refusal

## Summary

Refuse redirects for Foursquare API and OAuth user requests before query-string
client credentials or access tokens can be forwarded to another destination.

## Requirements

- Apply one shared `http.Client.CheckRedirect` policy that returns
  `http.ErrUseLastResponse`.
- Install the policy on copied Foursquare service clients while preserving
  caller transports, timeouts, and caller-owned configuration.
- Install the same policy on the App Engine `urlfetch` client used for OAuth
  user lookup.
- Preserve status, exact final URL, media type, response size, decoding, edit,
  and generic logging behavior.
- Add behavioral tests and mutation-sensitive static contracts.
- Record actual validation without claiming live credentials, API calls, OAuth
  redirects, or App Engine deployment.

## Non-Goals

- Changing Foursquare endpoints, query parameters, credentials, schemas,
  response limits, timeouts, rate limiting, cookies, or UI behavior.
- Replacing App Engine `urlfetch`, OAuth2, or the legacy deployment model.
- Following same-origin redirects or introducing a redirect allowlist.

## Verification Plan

- Run focused redirect-policy tests and the full `go test ./...` suite.
- Run `make check`, `make lint`, `make test`, and `make build` from the root,
  plus the absolute Makefile check from an external directory.
- Reject policy removal, redirect acceptance, service override removal,
  App Engine client override removal, and plan evidence drift.
- Run formatting, vet, diff, artifact, mode, and secret audits.
- Capture one bounded exact-head hosted snapshot after push.

## Work Completed

- Added shared `fsq.RefuseRedirect` policy returning
  `http.ErrUseLastResponse`.
- Applied the policy to copied Foursquare service clients while retaining
  caller transports, positive timeouts, and caller configuration ownership.
- Applied the same policy to the App Engine `urlfetch` client used for OAuth
  user lookup.
- Added focused behavior tests, static contracts, and maintained guidance.

## Verification Completed

- `go test -race -count=1 ./...` and `go vet ./...` passed in an isolated copy
  of the exact source tree.
- All four Make gates passed, and the absolute-Makefile check passed from an
  external directory.
- The policy removal mutation failed compilation and the structural contract.
- The redirect acceptance mutation failed focused behavior tests.
- The service override mutation failed focused behavior tests.
- The OAuth client override mutation failed the App Engine client contract.
- The focused test mutation failed the test-presence contract.
- The plan evidence mutation failed the completed-evidence contract.
- Formatting, shell syntax, and `git diff --check` passed before final audit.
- No live Foursquare credentials, API calls, OAuth redirects, or App Engine
  deployment were exercised.
- The hosted pull-request check is captured after push and recorded in the
  exact-head tracker evidence rather than claimed by this pre-push plan.
