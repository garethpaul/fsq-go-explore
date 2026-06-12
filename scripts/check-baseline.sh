#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
PLAN="$ROOT_DIR/docs/plans/2026-06-08-fsq-go-explore-go-baseline.md"
EDIT_PLAN="$ROOT_DIR/docs/plans/2026-06-09-fsq-propose-edit-post-only.md"
VENUE_ID_PLAN="$ROOT_DIR/docs/plans/2026-06-09-fsq-venue-id-boundary.md"
SEARCH_PARAM_PLAN="$ROOT_DIR/docs/plans/2026-06-09-fsq-search-param-length.md"
EDIT_ID_FIRST_PLAN="$ROOT_DIR/docs/plans/2026-06-09-fsq-edit-page-id-first.md"
OAUTH_CODE_PLAN="$ROOT_DIR/docs/plans/2026-06-09-fsq-oauth-code-boundary.md"
MAKE_GATES_PLAN="$ROOT_DIR/docs/plans/2026-06-09-fsq-go-make-gate-aliases.md"
USER_CACHE_KEY_PLAN="$ROOT_DIR/docs/plans/2026-06-09-fsq-user-cache-key-boundary.md"
ETAG_MATCH_PLAN="$ROOT_DIR/docs/plans/2026-06-09-fsq-etag-exact-match.md"
LOGIN_PROTECT_KEY_PLAN="$ROOT_DIR/docs/plans/2026-06-09-fsq-login-protect-cache-key.md"
EDIT_FORM_PARSE_PLAN="$ROOT_DIR/docs/plans/2026-06-09-fsq-propose-edit-form-parse-boundary.md"
CI_PLAN="$ROOT_DIR/docs/plans/2026-06-10-ci-baseline.md"
RATE_LIMITER_KEY_CAP_PLAN="$ROOT_DIR/docs/plans/2026-06-10-fsq-rate-limiter-key-cap.md"
RATE_LIMITER_REFILL_PLAN="$ROOT_DIR/docs/plans/2026-06-12-fsq-rate-limiter-refill.md"
WORKFLOW="$ROOT_DIR/.github/workflows/check.yml"

require_file() {
  path=$1
  if [ ! -f "$ROOT_DIR/$path" ]; then
    printf '%s\n' "Required file missing: $path" >&2
    exit 1
  fi
}

for path in \
  ".gitignore" \
  ".github/workflows/check.yml" \
  "CHANGES.md" \
  "Makefile" \
  "README.md" \
  "SECURITY.md" \
  "VISION.md" \
  "go.mod" \
  "go.sum" \
  "app.yaml" \
  "auth.go" \
  "auth_test.go" \
  "cache.go" \
  "cache_test.go" \
  "edit.go" \
  "edit_test.go" \
  "main.go" \
  "pages.go" \
  "search.go" \
  "search_test.go" \
  "fsq/api.go" \
  "fsq/api_test.go" \
  "fsq/keys.go" \
  "fsq/keys_test.go" \
  "limiter/limiter.go" \
  "limiter/config/config_test.go" \
  "docs/plans/2026-06-12-fsq-rate-limiter-refill.md" \
  "docs/plans/2026-06-10-fsq-rate-limiter-key-cap.md" \
  "docs/plans/2026-06-09-fsq-login-protect-cache-key.md" \
  "docs/plans/2026-06-10-ci-baseline.md" \
  "docs/plans/2026-06-09-fsq-propose-edit-form-parse-boundary.md" \
  "docs/plans/2026-06-09-fsq-etag-exact-match.md" \
  "docs/plans/2026-06-09-fsq-search-param-length.md" \
  "docs/plans/2026-06-09-fsq-edit-page-id-first.md" \
  "docs/plans/2026-06-09-fsq-go-make-gate-aliases.md" \
  "docs/plans/2026-06-09-fsq-user-cache-key-boundary.md" \
  "docs/plans/2026-06-09-fsq-oauth-code-boundary.md" \
  "docs/plans/2026-06-09-fsq-venue-id-boundary.md" \
  "docs/plans/2026-06-09-fsq-propose-edit-post-only.md" \
  "docs/plans/2026-06-08-fsq-go-explore-go-baseline.md"; do
  require_file "$path"
done

makefile="$ROOT_DIR/Makefile"
if ! grep -Eq '^\.PHONY: .*build.*check.*lint.*test|^\.PHONY: .*build.*lint.*test.*check' "$makefile" ||
  ! grep -Fq "lint test build: check" "$makefile"; then
  printf '%s\n' "Makefile must expose lint, test, build, and check gate targets." >&2
  exit 1
fi

if command -v go >/dev/null 2>&1; then
  unformatted=$(find "$ROOT_DIR" -name '*.go' -not -path "$ROOT_DIR/.git/*" -print | xargs gofmt -l)
  if [ -n "$unformatted" ]; then
    printf '%s\n' "Go files need gofmt:" >&2
    printf '%s\n' "$unformatted" >&2
    exit 1
  fi
  (cd "$ROOT_DIR" && go vet ./...)
  (cd "$ROOT_DIR" && go test ./...)
  (cd "$ROOT_DIR" && go mod tidy -diff)
else
  printf '%s\n' "go is required for fsq-go-explore verification." >&2
  exit 1
fi

if git -C "$ROOT_DIR" grep -nE '^[[:space:]]*"(fsq|limiter|limiter/config|limiter/errors|limiter/libstring)"|^[[:space:]]*"(appengine|appengine/memcache)"' -- '*.go'; then
  printf '%s\n' "Go imports must use module-qualified local and modern App Engine paths." >&2
  exit 1
fi

if git -C "$ROOT_DIR" grep -nE 'log\.Print\(accessToken|log\.Print\(user\)|near=%q|fmt\.Print|key_text|NewCFBEncrypter|panic\(|template\.Must|stateString|STATE_STR|request failed: %v' -- '*.go'; then
  printf '%s\n' "Go source must avoid raw credential/location logging, reversible cache keys, and panic-based request paths." >&2
  exit 1
fi

if ! grep -Fq "newOAuthState" "$ROOT_DIR/auth.go" ||
  ! grep -Fq "oauthStateCookieName" "$ROOT_DIR/auth.go" ||
  ! grep -Fq "TestNewOAuthStateReturnsDistinctOpaqueValues" "$ROOT_DIR/auth_test.go"; then
  printf '%s\n' "OAuth redirects must use tested per-login state cookies." >&2
  exit 1
fi

if ! grep -Fq 'strings.TrimSpace(r.FormValue("code"))' "$ROOT_DIR/auth.go" ||
  ! grep -Fq "missing authorization code" "$ROOT_DIR/auth.go" ||
  ! grep -Fq "http.StatusBadRequest" "$ROOT_DIR/auth.go" ||
  ! grep -Fq "TestRedirectRejectsMissingAuthorizationCodeBeforeExchange" "$ROOT_DIR/auth_test.go"; then
  printf '%s\n' "OAuth callbacks must reject missing authorization codes before exchange work." >&2
  exit 1
fi

if ! grep -Fq "const userCacheKeyPrefix = \"user:\"" "$ROOT_DIR/auth.go" ||
  ! grep -Fq "func validUserCacheKey" "$ROOT_DIR/auth.go" ||
  ! grep -Fq "len(digest) != 64" "$ROOT_DIR/auth.go" ||
  ! grep -Fq "if !validUserCacheKey(key)" "$ROOT_DIR/auth.go" ||
  ! grep -Fq "cookie == nil || !validUserCacheKey(cookie.Value)" "$ROOT_DIR/auth.go" ||
  ! grep -Fq "TestValidUserCacheKeyAcceptsGeneratedUserKeys" "$ROOT_DIR/auth_test.go" ||
  ! grep -Fq "TestGetAccessTokenRejectsMalformedCacheKeysBeforeLookup" "$ROOT_DIR/auth_test.go" ||
  ! grep -Fq "TestLoginProtectRejectsMalformedAuthCookie" "$ROOT_DIR/auth_test.go"; then
  printf '%s\n' "Foursquare auth cookies must validate generated user cache keys before memcache lookup." >&2
  exit 1
fi

if grep -Eq 'VENUE_URL \+ (id|venueId)' "$ROOT_DIR/fsq/api.go" ||
  ! grep -Fq "url.PathEscape(id)" "$ROOT_DIR/fsq/api.go" ||
  ! grep -Fq "url.PathEscape(venueId)" "$ROOT_DIR/fsq/api.go"; then
  printf '%s\n' "Venue IDs must be path-escaped before Foursquare request URLs are built." >&2
  exit 1
fi

if ! grep -Fq "sha256.Sum256" "$ROOT_DIR/fsq/keys.go" ||
  ! grep -Fq 'cacheKey("search"' "$ROOT_DIR/fsq/keys.go" ||
  ! grep -Fq 'cacheKey("user"' "$ROOT_DIR/fsq/keys.go"; then
  printf '%s\n' "Cache key helpers must be deterministic, opaque, and namespaced." >&2
  exit 1
fi

if ! grep -Fq "appEngineLocationFallback" "$ROOT_DIR/search.go" ||
  ! grep -Fq "Chicago, IL" "$ROOT_DIR/search_test.go"; then
  printf '%s\n' "Search parameter parsing must have tested location fallback behavior." >&2
  exit 1
fi

if grep -Fq "strings.Contains(match, key)" "$ROOT_DIR/cache.go" ||
  ! grep -Fq "func etagMatches" "$ROOT_DIR/cache.go" ||
  ! grep -Fq "strings.Split(headerValue" "$ROOT_DIR/cache.go" ||
  ! grep -Fq "TestHeaderCacheRejectsPartialETagMatches" "$ROOT_DIR/cache_test.go" ||
  ! grep -Fq "TestHeaderCacheAcceptsExactETagFromList" "$ROOT_DIR/cache_test.go"; then
  printf '%s\n' "Header cache ETag matching must be exact and covered by tests." >&2
  exit 1
fi

if ! grep -Fq "const maxSearchParamRunes" "$ROOT_DIR/search.go" ||
  ! grep -Fq "func normalizeSearchParam" "$ROOT_DIR/search.go" ||
  ! grep -Fq "TestSearchParamParserLimitsLongInputs" "$ROOT_DIR/search_test.go"; then
  printf '%s\n' "Search parameter parsing must trim and bound query/location inputs with test coverage." >&2
  exit 1
fi

if ! grep -Fq "r.Method != http.MethodPost" "$ROOT_DIR/edit.go" ||
  ! grep -Fq "http.StatusMethodNotAllowed" "$ROOT_DIR/edit.go" ||
  ! grep -Fq "TestProposeEditRejectsNonPostRequests" "$ROOT_DIR/edit_test.go"; then
  printf '%s\n' "Venue edit submissions must reject non-POST requests with test coverage." >&2
  exit 1
fi

if ! grep -Fq 'strings.TrimSpace(r.FormValue("id"))' "$ROOT_DIR/edit.go" ||
  ! grep -Fq "missing venue id" "$ROOT_DIR/edit.go" ||
  ! grep -Fq "http.StatusBadRequest" "$ROOT_DIR/edit.go" ||
  ! grep -Fq "TestEditPageRejectsMissingVenueIDBeforeAuth" "$ROOT_DIR/edit_test.go" ||
  ! grep -Fq "TestProposeEditRejectsMissingVenueID" "$ROOT_DIR/edit_test.go"; then
  printf '%s\n' "Venue detail/edit handlers must reject missing or blank venue IDs with test coverage." >&2
  exit 1
fi

if ! grep -Fq 'strings.TrimSpace(r.Form.Get("id"))' "$ROOT_DIR/edit.go" ||
  ! grep -Fq "invalid venue edit form" "$ROOT_DIR/edit.go" ||
  ! grep -Fq "venue edit form parse failed" "$ROOT_DIR/edit.go" ||
  ! grep -Fq "TestProposeEditRejectsMalformedFormBeforeAuth" "$ROOT_DIR/edit_test.go"; then
  printf '%s\n' "Venue edit submissions must reject malformed forms before auth or Foursquare work." >&2
  exit 1
fi

if ! grep -Fq "defaultMaxTrackedKeys = 10000" "$ROOT_DIR/limiter/config/config.go" ||
  ! grep -Fq "list.New()" "$ROOT_DIR/limiter/config/config.go" ||
  ! grep -Fq "MoveToFront" "$ROOT_DIR/limiter/config/config.go" ||
  ! grep -Fq "tokenBucketOrder.Back()" "$ROOT_DIR/limiter/config/config.go" ||
  ! grep -Fq "delete(l.tokenBuckets, oldestKey)" "$ROOT_DIR/limiter/config/config.go" ||
  ! grep -Fq "TestLimiterCapsTrackedKeys" "$ROOT_DIR/limiter/config/config_test.go" ||
  ! grep -Fq "TestLimiterEvictsLeastRecentlyUsedKey" "$ROOT_DIR/limiter/config/config_test.go"; then
  printf '%s\n' "Rate limiter keys must remain capped with recency-sensitive eviction tests." >&2
  exit 1
fi

if ! grep -Fq "func newTokenBucket" "$ROOT_DIR/limiter/config/config.go" ||
  ! grep -Fq "float64(max) / ttl.Seconds()" "$ROOT_DIR/limiter/config/config.go" ||
  ! grep -Fq "max <= 0 || ttl <= 0" "$ROOT_DIR/limiter/config/config.go" ||
  ! grep -Fq "TestLimiterRefillsConfiguredMaximumAcrossTTL" "$ROOT_DIR/limiter/config/config_test.go" ||
  ! grep -Fq "TestLimiterRejectsInvalidRateConfiguration" "$ROOT_DIR/limiter/config/config_test.go"; then
  printf '%s\n' "Rate limiter buckets must refill Max requests over TTL and reject invalid configurations." >&2
  exit 1
fi

if ! grep -Fq "go test ./..." "$ROOT_DIR/README.md" ||
  ! grep -Fq "GitHub Actions" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make lint" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make test" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make build" "$ROOT_DIR/README.md" ||
  ! grep -Fq "make check" "$ROOT_DIR/README.md" ||
  ! grep -Fq "non-POST requests are rejected" "$ROOT_DIR/README.md" ||
  ! grep -Fq "missing venue IDs are rejected" "$ROOT_DIR/README.md" ||
  ! grep -Fq "Malformed edit forms are rejected" "$ROOT_DIR/README.md" ||
  ! grep -Fq "length-bounded" "$ROOT_DIR/README.md" ||
  ! grep -Fq "missing OAuth authorization codes are rejected" "$ROOT_DIR/README.md" ||
  ! grep -Fq "ETag comparisons are exact" "$ROOT_DIR/README.md" ||
  ! grep -Fq "Protected routes validate generated auth cookie cache keys" "$ROOT_DIR/README.md" ||
  ! grep -Fq "10,000 rate-limiter keys" "$ROOT_DIR/README.md" ||
  ! grep -Fq 'refills those `Max` requests' "$ROOT_DIR/README.md" ||
  ! grep -Fq "user cache keys" "$ROOT_DIR/README.md" ||
  ! grep -Fq "FSQ_CLIENT_ID" "$ROOT_DIR/README.md"; then
  printf '%s\n' "README must document Go verification and Foursquare env configuration." >&2
  exit 1
fi

if ! grep -Fq "scripts/check-baseline.sh" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "GitHub Actions" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "make lint" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "make test" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "make build" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "cache-key generation" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Venue edit submissions reject non-POST requests" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "missing venue IDs" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "malformed edit forms" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "missing OAuth authorization codes" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "ETag comparisons are exact" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Protected routes validate generated auth cookie cache keys" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "10,000 tracked request keys" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq 'refill `Max` requests' "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "user cache keys" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "length-bounded" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Go module" "$ROOT_DIR/VISION.md"; then
  printf '%s\n' "VISION must describe the current Go baseline." >&2
  exit 1
fi

if ! grep -Fq "Malformed venue edit forms should be rejected" "$ROOT_DIR/SECURITY.md" ||
  ! grep -Fq "least-recently-used" "$ROOT_DIR/SECURITY.md" ||
  ! grep -Fq 'refill `Max` requests over `TTL`' "$ROOT_DIR/SECURITY.md"; then
  printf '%s\n' "SECURITY must document the malformed venue edit form boundary." >&2
  exit 1
fi

exact_line_count() {
  awk -v expected="$2" '$0 == expected { count += 1 } END { print count + 0 }' "$1"
}

if [ "$(exact_line_count "$WORKFLOW" 'permissions:')" -ne 1 ] || \
  [ "$(exact_line_count "$WORKFLOW" '  contents: read')" -ne 1 ] || \
  grep -Eq '^[[:space:]]+permissions:' "$WORKFLOW" || \
  grep -Eq '(^|[[:space:]])write-all([[:space:]]|$)' "$WORKFLOW" || \
  grep -Eq '^[[:space:]]+[^#][^:]*:[[:space:]]*write([[:space:]]*(#.*)?)?$' "$WORKFLOW"; then
  printf '%s\n' "GitHub Actions must keep one top-level read-only permissions block." >&2
  exit 1
fi

if [ "$(grep -Fc 'uses: actions/checkout@' "$WORKFLOW")" -ne 1 ] || \
  ! grep -Fq 'uses: actions/checkout@df4cb1c069e1874edd31b4311f1884172cec0e10 # v6.0.3' "$WORKFLOW" || \
  [ "$(exact_line_count "$WORKFLOW" '          persist-credentials: false')" -ne 1 ]; then
  printf '%s\n' "GitHub Actions must keep one pinned, credential-free checkout step." >&2
  exit 1
fi

if [ "$(grep -Fc 'uses: actions/setup-go@' "$WORKFLOW")" -ne 1 ] || \
  ! grep -Fq 'uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6.4.0' "$WORKFLOW" || \
  [ "$(exact_line_count "$WORKFLOW" '          go-version-file: go.mod')" -ne 1 ] || \
  [ "$(exact_line_count "$WORKFLOW" '        run: make check')" -ne 1 ]; then
  printf '%s\n' "GitHub Actions must keep the pinned Go setup and canonical make check gate." >&2
  exit 1
fi

for workflow_contract in \
  '  workflow_dispatch:' \
  '  cancel-in-progress: true' \
  '    runs-on: ubuntu-24.04' \
  '    timeout-minutes: 10'; do
  if [ "$(exact_line_count "$WORKFLOW" "$workflow_contract")" -ne 1 ]; then
    printf '%s\n' "GitHub Actions is missing required workflow contract: $workflow_contract" >&2
    exit 1
  fi
done

if ! grep -Fq "status: completed" "$PLAN"; then
  printf '%s\n' "Plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$EDIT_PLAN"; then
  printf '%s\n' "Venue edit POST-only plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$VENUE_ID_PLAN"; then
  printf '%s\n' "Venue ID boundary plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$SEARCH_PARAM_PLAN"; then
  printf '%s\n' "Search parameter length plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$EDIT_ID_FIRST_PLAN"; then
  printf '%s\n' "Edit-page ID-first plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$OAUTH_CODE_PLAN"; then
  printf '%s\n' "OAuth code boundary plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$MAKE_GATES_PLAN"; then
  printf '%s\n' "Make gate alias plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$USER_CACHE_KEY_PLAN"; then
  printf '%s\n' "User cache key boundary plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$ETAG_MATCH_PLAN"; then
  printf '%s\n' "ETag exact match plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$LOGIN_PROTECT_KEY_PLAN"; then
  printf '%s\n' "Login protect cache key plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$EDIT_FORM_PARSE_PLAN"; then
  printf '%s\n' "Edit form parse boundary plan must be marked completed." >&2
  exit 1
fi

if ! grep -Fq "make check" "$LOGIN_PROTECT_KEY_PLAN"; then
  printf '%s\n' "Login protect cache key plan must record make check verification." >&2
  exit 1
fi

if ! grep -Fq "make check" "$EDIT_FORM_PARSE_PLAN"; then
  printf '%s\n' "Edit form parse boundary plan must record make check verification." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$CI_PLAN" ||
  ! grep -Fq "make check" "$CI_PLAN"; then
  printf '%s\n' "CI baseline plan must record completed make check verification." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$RATE_LIMITER_KEY_CAP_PLAN" ||
  ! grep -Fq "Mutations disabling the cap or recency refresh must fail" "$RATE_LIMITER_KEY_CAP_PLAN"; then
  printf '%s\n' "Rate limiter key-cap plan must record completed mutation verification." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$RATE_LIMITER_REFILL_PLAN" ||
  ! grep -Fq "Mutations restoring one-token-per-TTL refill" "$RATE_LIMITER_REFILL_PLAN"; then
  printf '%s\n' "Rate limiter refill plan must record completed mutation verification." >&2
  exit 1
fi

printf '%s\n' "fsq-go-explore Go baseline checks passed."
