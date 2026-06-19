---
title: OAuth User Response Origin
date: 2026-06-14
status: completed
execution: code
---

## Context

The OAuth callback client may follow redirects before returning the Foursquare
user-profile response. The decoder currently validates status and body size,
but it can accept JSON from an unexpected final scheme, authority, or path and
use that payload to create authenticated application state.

## Requirements

- Require the final OAuth user-profile response URL to remain HTTPS,
  `api.foursquare.com`, without userinfo, an explicit port, or a fragment, and
  on the exact `/v2/users/self` path.
- Require a JSON response media type before reading the response body.
- Apply status, final-URL, and media-type validation before any body read while
  preserving the existing 1 MiB limit, error handling, and response closure.
- Add focused unread-body tests for rejected final URLs and media types.
- Add mutation-sensitive static contracts and synchronized repository guidance.

## Non-Goals

- Disabling redirects globally or changing the App Engine transport.
- Changing OAuth state, token exchange, cache keys, cookies, or session expiry.
- Logging URLs, query values, credentials, tokens, or response bodies.
- Claiming live OAuth, Foursquare, App Engine, or memcache validation.

## Verification Plan

- Run focused auth tests, all-package tests, race tests, vet, module integrity,
  all four Make gates, and the external-directory Make gate.
- Reject mutations that weaken scheme, host, path, media-type, ordering, test,
  documentation, or completed-plan contracts.
- Audit formatting, shell/Python syntax, the exact diff, generated artifacts,
  protected dependency/workflow files, whitespace, conflict markers, and
  changed-line credential patterns before commit and push.
- Take one bounded exact-head pull-request and security-alert snapshot without
  polling.

## Work Completed

- Added exact final-response URL validation for HTTPS,
  `api.foursquare.com`, absent userinfo/port/fragment, and the exact
  `/v2/users/self` path.
- Required JSON or structured-suffix JSON media types before OAuth profile body
  reads while preserving status-first validation and the existing 1 MiB limit.
- Added unread-body regressions for unexpected final URLs and media types,
  accepted-endpoint coverage, a mutation-sensitive baseline contract, and
  synchronized repository guidance.

## Verification Completed

- Focused auth tests, uncached all-package tests, the race detector, vet, and
  module integrity passed.
- All four Make gates and the external-directory Make gate passed.
- Eight hostile mutations were rejected across scheme, host, path, guard,
  media-type, structured-suffix, unread-body-test, and plan-evidence weakening.
- Formatting, shell syntax, exact diff, generated-artifact, protected
  dependency/workflow, whitespace, conflict-marker, and changed-line credential
  audits passed.
- No live OAuth flow, Foursquare request, App Engine service, or memcache access
  was performed.
- The hosted pull-request and security-alert result is recorded against the
  exact pushed head in the external engineering tracker.
