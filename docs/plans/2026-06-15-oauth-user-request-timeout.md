---
title: OAuth User Request Timeout
type: security
status: in_progress
date: 2026-06-15
---

# OAuth User Request Timeout

## Summary

Bound the App Engine HTTP client used to fetch the authenticated Foursquare
user profile so a stalled upstream request cannot hold the OAuth callback open
indefinitely.

## Requirements

- Define one 10-second end-to-end timeout for the OAuth user-profile client.
- Apply the timeout alongside the existing App Engine transport and redirect
  refusal policy.
- Preserve OAuth state, token exchange, response status, final endpoint, media
  type, body-size, decoding, cookie, cache, and logging behavior.
- Add a direct behavioral assertion and mutation-sensitive static contracts.
- Record actual validation without claiming live credentials, API calls, OAuth
  callbacks, or App Engine deployment.

## Non-Goals

- Changing the shared Foursquare package client timeout or caller ownership.
- Adding retries, backoff, cancellation plumbing, or per-operation deadlines.
- Replacing App Engine `urlfetch`, OAuth2, or the legacy deployment model.

## Verification Plan

- Run the focused OAuth client test and the full race-enabled Go test suite.
- Run `go vet`, all Make gates, and the absolute Makefile check externally.
- Reject timeout removal, value drift, client assignment removal, weakened
  focused tests, missing guidance, and incomplete plan evidence.
- Run formatting, shell syntax, module, diff, artifact, mode, and secret audits.
- Capture one bounded exact-head hosted snapshot after push.

## Work Completed

Pending implementation.

## Verification Completed

Pending implementation and verification.
