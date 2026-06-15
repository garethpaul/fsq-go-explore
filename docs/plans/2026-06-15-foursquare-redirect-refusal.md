---
title: Foursquare Redirect Refusal
type: security
status: in_progress
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

Pending implementation.

## Verification Completed

Pending implementation and validation.
