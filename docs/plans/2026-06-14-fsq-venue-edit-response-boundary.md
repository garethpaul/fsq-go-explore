---
title: Foursquare Venue Edit Response Boundary
date: 2026-06-14
status: completed
execution: code
---

## Context

`VenueEdit` drains the remote response without a byte limit before evaluating
its status. A non-success or oversized response can therefore consume avoidable
network and CPU resources, and 3xx responses are not rejected by the local
status logic.

## Requirements

- Require the shared exact 2xx success predicate before reading the response.
- Bound successful response disposal to the existing 2 MiB response limit plus
  one detection byte.
- Preserve short successful bodies and read errors without parsing or logging
  response content.
- Add focused tests proving non-2xx bodies are unread and oversized bodies stop
  after the detection byte.
- Add mutation-sensitive static contracts, synchronized guidance, and truthful
  completed verification evidence.

## Non-Goals

- Parsing venue-edit response JSON or changing the request body.
- Changing credentials, endpoints, timeouts, redirect policy, dependencies, or
  legacy App Engine deployment behavior.
- Claiming live Foursquare or App Engine validation.

## Verification Plan

- Run focused and all-package Go tests, race tests, vet, module integrity, all
  four Make gates, and the external-directory Make gate.
- Reject mutations that remove status ordering, restore unbounded draining,
  remove oversize detection, weaken tests, or falsify completed plan evidence.
- Audit the exact diff, generated artifacts, formatting, and credential-like
  additions before commit and push.
- Take one bounded exact-head pull-request and code-scanning snapshot without
  polling.

## Work Completed

- Added a shared bounded response discard that reads at most 2 MiB plus one
  detection byte and reports the existing oversized-response sentinel.
- Required exact 2xx venue-edit status before invoking the discard, leaving
  rejected response bodies unread.
- Added focused runtime and static contracts plus synchronized guidance.

## Verification Completed

- Focused package tests, all-package tests, race tests, vet, module integrity,
  all four Make gates, and the external-directory Make gate passed.
- The status ordering mutation failed before a rejected body could be read.
- The unbounded discard mutation failed the source boundary contract.
- The oversize detection mutation failed the limit-plus-one contract.
- The focused test mutation failed the unique regression-test contract.
- The plan evidence mutation failed the completed-evidence contract.
- Formatting, shell syntax, exact diff, artifact, dependency/workflow,
  whitespace, conflict-marker, and changed-line credential audits passed.
- No live Foursquare request or App Engine deployment was performed.
- The hosted pull-request check and code-scanning result are recorded against
  the exact pushed head in the external engineering tracker.
