# Security Policy

## Supported Versions

The supported security scope for `fsq-go-explore` is the current default branch, `master`. Older commits, tags, branches, forks, demos, and generated artifacts are not actively supported unless the repository explicitly marks them as maintained.

Project summary: Foursquare GoLang Explore Sample

## Reporting a Vulnerability

Please report suspected vulnerabilities through GitHub's private vulnerability reporting or by opening a draft GitHub Security Advisory for `garethpaul/fsq-go-explore` when that option is available. If GitHub does not show a private reporting option for this repository, contact the repository owner through GitHub and avoid posting exploit details publicly until the issue can be assessed.

Do not open a public issue that includes exploit code, secrets, personal data, or detailed reproduction steps for an unpatched vulnerability.

## What to Include

Helpful reports include:

- the affected file, endpoint, permission, dependency, or workflow
- a concise impact statement explaining what an attacker could do
- reproduction steps using test data and accounts you control
- the branch, commit SHA, platform version, device, runtime, or dependency versions used
- logs, screenshots, or proof-of-concept snippets that demonstrate impact without exposing private data

## Project Security Posture

- This repository appears to be a public sample, documentation, or utility project. The active security scope is the code and documentation on the default branch.
- Review found authentication, token, or session-related code paths; changes in those areas should receive security-focused review before merge.
- Review found external API integrations or credential-adjacent configuration; changes in those areas should receive security-focused review before merge.
- Review found network clients, sockets, web APIs, or service endpoints; changes in those areas should receive security-focused review before merge.
- Review found mobile permission or privacy-sensitive data handling; changes in those areas should receive security-focused review before merge.
- Review found file, document, data, or media parsing flows; changes in those areas should receive security-focused review before merge.
- Go dependency manifests are tracked in `go.mod` and `go.sum`. Keep them in sync with code changes and review App Engine, OAuth, and rate-limit dependency updates carefully.
- GitHub Actions runs formatting, vet, tests, module-integrity checks, and the
  static security baseline with read-only permissions and credential-free
  checkout; keep it aligned when changing auth, cache, or Foursquare request
  code.

## Service and API Notes

For web services, APIs, sockets, or scraping workflows, prioritize reports involving authentication bypass, authorization errors, injection, server-side request forgery, unsafe deserialization, credential leakage, data exposure, or denial-of-service conditions. Use test accounts and minimal proof-of-concept traffic only.
Search query and location parameters should stay length-bounded before they are
sent to Foursquare or used in cache keys.
Malformed venue edit forms should be rejected before auth-cookie lookup, token
cache work, or Foursquare edit API calls.
Venue edit request bodies should remain limited to 64 KiB and return `413
Request Entity Too Large` before auth-cookie lookup, token cache work, or
Foursquare edit API calls when exceeded.
OAuth callbacks should reject missing authorization codes before token exchange
work starts, even when the state cookie matches.
OAuth user-profile responses should reject non-2xx statuses before body reads
and enforce a 1 MiB limit before JSON decoding or session creation.
Auth cookies should validate generated user cache keys before memcache lookup so
malformed cookie values do not reach access-token cache work.
Protected routes should reject malformed auth-cookie cache keys before handler
work starts.
ETag comparisons should stay exact before returning `304 Not Modified`.
Rate-limiter storage should remain capped and use least-recently-used eviction
so attacker-controlled key rotation cannot grow process memory without bound.
Rate-limiter buckets should refill `Max` requests over `TTL`, and non-positive
rate configurations should fail closed instead of disabling throttling.
Foursquare JSON response bodies must remain limited to 2 MiB before parsing;
oversized or failed reads should not reach JSON unmarshalling.
Non-2xx Foursquare search and venue detail responses must not reach JSON decoding;
status failures should log only the numeric status and return empty results.
Foursquare HTTP clients should use a 10-second default end-to-end timeout when
no positive caller timeout is configured; service construction must not mutate
caller-owned configuration.

## Dependency and Supply Chain Security

Dependency updates should come from trusted package managers and should keep lockfiles in sync when lockfiles exist. Do not commit credentials, private keys, tokens, generated secrets, or machine-local configuration. If a vulnerability depends on a compromised package, typosquatting risk, insecure transitive dependency, or unsafe build step, include the package name, affected version, and the path through which it is used.

## Safe Research Guidelines

Good-faith research is welcome when it stays within these boundaries:

- use only accounts, devices, data, and infrastructure that you own or have explicit permission to test
- avoid destructive actions, persistence, spam, phishing, social engineering, or denial-of-service testing
- minimize access to personal data and stop testing immediately if private data is exposed
- do not exfiltrate secrets or third-party data; report the minimum evidence needed to verify impact
- keep vulnerability details confidential until the maintainer has assessed the report

## Maintainer Response

The maintainer will review complete reports as availability allows, prioritize issues by exploitability and impact, and coordinate a fix or mitigation when the affected code is still maintained. For sample, archived, or educational repositories, the likely remediation may be documentation, dependency updates, or clearly marking unsupported code rather than a production-style patch release.
