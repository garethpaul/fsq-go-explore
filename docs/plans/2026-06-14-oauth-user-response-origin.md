---
title: OAuth User Response Origin
date: 2026-06-14
status: planned
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

Pending implementation.

## Verification Completed

Pending implementation and verification.
