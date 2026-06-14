---
title: Foursquare Response Final URL Boundary
date: 2026-06-14
status: completed
execution: code
---

## Context

The shared `http.Client` follows redirects, while search, venue-detail, and
venue-edit handling currently trusts any final response that has a 2xx status.
That allows a redirected response from a different scheme, authority, or path
to reach JSON decoding or successful-body disposal.

## Requirements

- Require every successful Foursquare response to retain HTTPS,
  `api.foursquare.com`, no userinfo, no explicit port, no fragment, and the
  operation's exact escaped API path.
- Run final URL validation after the existing 2xx check and before content-type
  validation, JSON decoding, or successful response disposal.
- Preserve dynamic query parameters, existing credentials, request methods,
  timeouts, response limits, media handling, and public APIs.
- Add focused tests proving exact endpoints are accepted and unexpected final
  endpoints are rejected before response-body reads.
- Add mutation-sensitive static contracts, synchronized guidance, and truthful
  completed verification evidence.

## Non-Goals

- Disabling redirects globally or changing the caller-provided transport.
- Logging final URLs, query values, credentials, tokens, or response bodies.
- Changing dependencies, App Engine deployment behavior, or Foursquare API
  versions.
- Claiming live Foursquare or App Engine validation.

## Verification Plan

- Start with focused helper and operation tests, then run all-package tests,
  race tests, vet, module integrity, all four Make gates, and the
  external-directory Make gate.
- Reject mutations that remove an operation check, weaken scheme/host/path
  matching, move validation after reads, remove focused tests, or falsify plan
  completion evidence.
- Audit the exact diff, generated artifacts, formatting, protected dependency
  and workflow files, whitespace, conflict markers, and changed-line credential
  patterns before commit and push.
- Take one bounded exact-head pull-request and security-alert snapshot without
  polling.

## Work Completed

- Added one exact final-response URL validator covering HTTPS, the Foursquare
  API host, absent userinfo/port/fragment, and each operation's escaped path.
- Required search, venue-detail, and venue-edit responses to pass the validator
  after 2xx status and before media checks, decoding, or bounded disposal.
- Added focused endpoint and unread-body regressions, a mutation-sensitive
  static contract, baseline integration, and synchronized repository guidance.

## Verification Completed

- The focused final URL tests passed for exact endpoints and all three rejected
  operation responses.
- All-package tests, race tests, vet, and module integrity passed.
- All four Make gates and the external-directory Make gate passed.
- The six hostile mutations were rejected across scheme, host, operation-call,
  validation-order, executable-test, and completed-plan evidence weakening.
- Formatting, shell/Python syntax, exact diff, generated-artifact, protected
  dependency/workflow, whitespace, conflict-marker, and changed-line credential
  audits passed.
- No live Foursquare request or App Engine deployment was performed.
- The hosted pull-request and security-alert result is recorded against the
  exact pushed head in the external engineering tracker.
