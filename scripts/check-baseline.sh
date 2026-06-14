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
EDIT_BODY_LIMIT_PLAN="$ROOT_DIR/docs/plans/2026-06-12-fsq-edit-body-limit.md"
RESPONSE_BODY_LIMIT_PLAN="$ROOT_DIR/docs/plans/2026-06-13-fsq-response-body-limit.md"
RESPONSE_STATUS_PLAN="$ROOT_DIR/docs/plans/2026-06-13-fsq-response-status-validation.md"
CLIENT_TIMEOUT_PLAN="$ROOT_DIR/docs/plans/2026-06-13-foursquare-client-timeout.md"
OAUTH_USER_RESPONSE_PLAN="$ROOT_DIR/docs/plans/2026-06-13-oauth-user-response-boundary.md"
LOCATION_INDEPENDENT_MAKE_PLAN="$ROOT_DIR/docs/plans/2026-06-13-location-independent-make.md"
RESPONSE_CONTENT_TYPE_PLAN="$ROOT_DIR/docs/plans/2026-06-14-fsq-response-content-type.md"
RESPONSE_CONTENT_TYPE_CHECK="$ROOT_DIR/scripts/check-response-content-type.py"
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
  "scripts/check-response-content-type.py" \
  "docs/plans/2026-06-14-fsq-response-content-type.md" \
  "docs/plans/2026-06-12-fsq-rate-limiter-refill.md" \
  "docs/plans/2026-06-12-fsq-edit-body-limit.md" \
  "docs/plans/2026-06-13-fsq-response-body-limit.md" \
  "docs/plans/2026-06-13-fsq-response-status-validation.md" \
  "docs/plans/2026-06-13-foursquare-client-timeout.md" \
  "docs/plans/2026-06-13-oauth-user-response-boundary.md" \
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

require_file "docs/plans/2026-06-13-location-independent-make.md"

python3 "$RESPONSE_CONTENT_TYPE_CHECK" \
  "$ROOT_DIR/fsq/api.go" \
  "$ROOT_DIR/fsq/api_test.go" \
  "$RESPONSE_CONTENT_TYPE_PLAN"

if ! grep -Fq 'ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))' "$ROOT_DIR/Makefile" ||
  ! grep -Fq '"$(ROOT)/scripts/check-baseline.sh"' "$ROOT_DIR/Makefile"; then
  printf '%s\n' "Makefile verification must resolve the checker from the loaded Makefile." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$LOCATION_INDEPENDENT_MAKE_PLAN" ||
  ! grep -Fq "from /tmp" "$LOCATION_INDEPENDENT_MAKE_PLAN" ||
  ! grep -Fq "absolute Makefile path" "$ROOT_DIR/README.md" ||
  ! grep -Fq "Made Go verification independent" "$ROOT_DIR/CHANGES.md"; then
  printf '%s\n' "Location-independent Make plan and guidance must record completed external verification." >&2
  exit 1
fi

makefile="$ROOT_DIR/Makefile"
if ! grep -Eq '^\.PHONY: .*build.*check.*lint.*test|^\.PHONY: .*build.*lint.*test.*check' "$makefile" ||
  ! grep -Fq "lint test build: check" "$makefile"; then
  printf '%s\n' "Makefile must expose lint, test, build, and check gate targets." >&2
  exit 1
fi

python3 - "$ROOT_DIR/fsq/api.go" "$ROOT_DIR/fsq/api_test.go" <<'PY'
import sys
from pathlib import Path

source = Path(sys.argv[1]).read_text()
tests = Path(sys.argv[2]).read_text()
source_contracts = (
    "maxFoursquareResponseBytes = 2 * 1024 * 1024",
    'errFoursquareResponseTooLarge = errors.New("foursquare response body exceeds 2 MiB")',
    "io.ReadAll(io.LimitReader(body, maxFoursquareResponseBytes+1))",
    "if len(data) > maxFoursquareResponseBytes",
    "return errFoursquareResponseTooLarge",
)
test_contracts = (
    "TestDecodeFoursquareResponseAcceptsExactLimit",
    "TestDecodeFoursquareResponseRejectsOversizeBody",
    "TestDecodeFoursquareResponsePreservesReadError",
    "TestDecodeFoursquareResponseRejectsEmptyBody",
    "TestDecodeFoursquareResponseRejectsMalformedJSON",
    "maxFoursquareResponseBytes+1",
    "errors.Is(err, errFoursquareResponseTooLarge)",
    "errors.Is(err, errTestReadFailure)",
)

if any(source.count(item) != 1 for item in source_contracts):
    raise SystemExit("Foursquare response decoding must keep one exact 2 MiB parse boundary.")
if any(tests.count(item) != 1 for item in test_contracts):
    raise SystemExit("Foursquare response parsing must keep exact-limit, oversize, and read-error tests.")
if "io.ReadAll(body)" in source:
    raise SystemExit("Foursquare response decoding must not read an unbounded body.")
PY

python3 - "$ROOT_DIR/fsq/api.go" "$ROOT_DIR/fsq/api_test.go" <<'PY'
import sys
from pathlib import Path

source = Path(sys.argv[1]).read_text()
tests = Path(sys.argv[2]).read_text()
constructor = source.split("func NewFoursquareService", 1)[-1].split("\n}\n", 1)[0]
required_source = (
    "foursquareRequestTimeout   = 10 * time.Second",
    "serviceConfig := *config",
    "client := config.Client",
    "if client.Timeout <= 0",
    "client.Timeout = foursquareRequestTimeout",
    "serviceConfig.Client = client",
    "return &FoursquareService{Config: &serviceConfig}",
)
if any(item not in source for item in required_source):
    raise SystemExit("Foursquare clients must receive the reviewed default timeout through cloned configuration.")
if "config.Client.Timeout =" in constructor:
    raise SystemExit("Foursquare service construction must not mutate the caller config.")

required_tests = (
    "TestNewFoursquareServiceDefaultsClientTimeout",
    "TestNewFoursquareServicePreservesExplicitClientTimeout",
    "TestNewFoursquareServiceDoesNotMutateCallerConfig",
    "explicitTimeout := 3 * time.Second",
    "service.Config == config",
)
if any(item not in tests for item in required_tests):
    raise SystemExit("Focused tests must preserve timeout defaulting, explicit preservation, and caller immutability.")
PY

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

python3 - "$ROOT_DIR/auth.go" "$ROOT_DIR/auth_test.go" <<'PY'
import sys
from pathlib import Path

source = Path(sys.argv[1]).read_text()
tests = Path(sys.argv[2]).read_text()
decoder = source.split("func decodeOAuthUserResponse", 1)[-1].split(
    "\n// Process a request and cache", 1
)[0]
required_source = (
    "maxOAuthUserResponseBytes = 1 * 1024 * 1024",
    "response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices",
    "io.ReadAll(io.LimitReader(response.Body, maxOAuthUserResponseBytes+1))",
    "len(body) > maxOAuthUserResponseBytes",
    "json.Unmarshal(body, wrapper)",
    "json.Unmarshal(wrapper.Response, user)",
)
if any(item not in source for item in required_source):
    raise SystemExit("OAuth user responses must keep status and 1 MiB decode boundaries.")
status = decoder.find("response.StatusCode < http.StatusOK")
read = decoder.find("io.ReadAll(io.LimitReader")
if status < 0 or read < 0 or status >= read:
    raise SystemExit("OAuth user response status validation must precede body reads.")

redirect = source.split("func Redirect", 1)[-1].split("func decodeOAuthUserResponse", 1)[0]
for item in ("defer p.Body.Close()", "user, err := decodeOAuthUserResponse(p)"):
    if redirect.count(item) != 1:
        raise SystemExit("OAuth callback must close and decode one user response.")

required_tests = (
    "TestDecodeOAuthUserResponseRejectsNonSuccessBeforeRead",
    "TestDecodeOAuthUserResponseAcceptsExactLimit",
    "TestDecodeOAuthUserResponseRejectsOversizeBody",
    "TestDecodeOAuthUserResponsePreservesReadError",
    "TestDecodeOAuthUserResponseRejectsMalformedPayloads",
    "body.readCalls != 0",
)
if any(tests.count(item) != 1 for item in required_tests):
    raise SystemExit("Focused OAuth user response boundary tests must remain unique.")
PY

if ! grep -Fq 'userCacheKeyPrefix        = "user:"' "$ROOT_DIR/auth.go" ||
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

python3 - "$ROOT_DIR/fsq/api.go" "$ROOT_DIR/fsq/api_test.go" <<'PY'
import sys
from pathlib import Path

source = Path(sys.argv[1]).read_text()
tests = Path(sys.argv[2]).read_text()
search = source.split("func (fsqs *FoursquareService) Search", 1)[-1].split(
    "\n// Details gets", 1
)[0]
details = source.split("func (fsqs *FoursquareService) VenueDetails", 1)[-1].split(
    "\n// ProposeEdit", 1
)[0]
guard = "if !successfulFoursquareStatus(r.StatusCode)"
decoder = "decodeFoursquareResponse(r.Body"

if source.count("func successfulFoursquareStatus(statusCode int) bool") != 1:
    raise SystemExit("Foursquare response status validation must use one shared success predicate.")
if "statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices" not in source:
    raise SystemExit("Foursquare success status validation must remain restricted to 2xx responses.")
for name, method in (("search", search), ("venue details", details)):
    if method.count(guard) != 1 or method.count(decoder) != 1:
        raise SystemExit(f"Foursquare {name} must retain one status guard and one bounded decoder call.")
    if method.index(guard) > method.index(decoder):
        raise SystemExit(f"Foursquare {name} status validation must run before response decoding.")

required_tests = (
    "TestSuccessfulFoursquareStatusAcceptsOnly2xx",
    "TestSearchRejectsNonSuccessResponseBeforeDecode",
    "TestVenueDetailsRejectsNonSuccessResponseBeforeDecode",
    "http.StatusContinue",
    "http.StatusMultipleChoices",
    "http.StatusInternalServerError",
    "http.StatusBadGateway",
    '"must-not-decode"',
)
if any(tests.count(item) < 1 for item in required_tests):
    raise SystemExit("Non-2xx Foursquare search and detail behavior must retain focused transport tests.")
PY

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

if ! grep -Fq "const maxVenueEditBodyBytes int64 = 64 << 10" "$ROOT_DIR/edit.go" ||
  ! grep -Fq "http.MaxBytesReader(w, r.Body, maxVenueEditBodyBytes)" "$ROOT_DIR/edit.go" ||
  ! grep -Fq "http.StatusRequestEntityTooLarge" "$ROOT_DIR/edit.go" ||
  ! grep -Fq "TestProposeEditAcceptsBodyAtLimitBeforeAuth" "$ROOT_DIR/edit_test.go" ||
  ! grep -Fq "TestProposeEditRejectsBodyOverLimitBeforeAuth" "$ROOT_DIR/edit_test.go"; then
  printf '%s\n' "Venue edit bodies must remain bounded and return 413 before auth work." >&2
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

python3 - "$RESPONSE_BODY_LIMIT_PLAN" <<'PY'
import re
import sys
from pathlib import Path

plan = Path(sys.argv[1]).read_text()
frontmatter = plan.split("---", 2)[1]
statuses = re.findall(r"^status: .+$", frontmatter, flags=re.MULTILINE)
verification = plan.split("## Verification Completed\n", 1)[-1]
required = (
    "unbounded read mutation failed",
    "limit drift mutation failed",
    "oversize test mutation failed",
    "hosted pull-request check",
)

if (
    statuses != ["status: completed"]
    or "## Verification Completed\n" not in plan
    or any(item not in verification for item in required)
    or re.search(r"\b(?:pending|todo|tbd|not run)\b", verification, re.IGNORECASE)
):
    raise SystemExit(
        "Foursquare response body limit plan must remain completed with actual verification recorded."
    )
PY

python3 - "$RESPONSE_STATUS_PLAN" <<'PY'
import re
import sys
from pathlib import Path

plan = Path(sys.argv[1]).read_text()
frontmatter = plan.split("---", 2)[1]
statuses = re.findall(r"^status: .+$", frontmatter, flags=re.MULTILINE)
verification = plan.split("## Verification Completed\n", 1)[-1]
required = (
    "status guard mutation failed",
    "decode ordering mutation failed",
    "non-2xx test mutation failed",
    "hosted pull-request check",
)

if (
    statuses != ["status: completed"]
    or "## Verification Completed\n" not in plan
    or any(item not in verification for item in required)
    or re.search(r"\b(?:pending|todo|tbd|not run)\b", verification, re.IGNORECASE)
):
    raise SystemExit(
        "Foursquare response status plan must remain completed with actual verification recorded."
    )
PY

if ! grep -Fq "Foursquare JSON response bodies are limited to 2 MiB" "$ROOT_DIR/README.md" ||
  ! grep -Fq "Foursquare JSON response bodies must remain limited to 2 MiB" "$ROOT_DIR/SECURITY.md" ||
  ! grep -Fq "Foursquare JSON response parsing is limited to 2 MiB" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Bounded Foursquare JSON response parsing to 2 MiB" "$ROOT_DIR/CHANGES.md"; then
  printf '%s\n' "Project guidance must document the Foursquare response parse boundary." >&2
  exit 1
fi

if ! grep -Fq "Non-2xx Foursquare search and venue detail responses are rejected before JSON decoding" "$ROOT_DIR/README.md" ||
  ! grep -Fq "Non-2xx Foursquare search and venue detail responses must not reach JSON decoding" "$ROOT_DIR/SECURITY.md" ||
  ! grep -Fq "Non-2xx Foursquare search and venue detail responses are rejected before decoding" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Rejected non-2xx Foursquare search and venue detail responses before JSON decoding" "$ROOT_DIR/CHANGES.md"; then
  printf '%s\n' "Project guidance must document the Foursquare response status boundary." >&2
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

python3 - "$EDIT_BODY_LIMIT_PLAN" <<'PY'
import re
import sys
from pathlib import Path

plan = Path(sys.argv[1]).read_text()
statuses = re.findall(r"^status: .+$", plan, flags=re.MULTILINE)
verification = plan.split("## Verification Completed\n", 1)[-1]
required = (
    "go test -race -count=1 ./...",
    "push run `27392746868`",
    "pull-request run `27392750651`",
    "push run `27392769894`",
    "CodeQL run `27402320766`",
)

if (
    statuses != ["status: completed"]
    or any(item not in verification for item in required)
    or re.search(r"\b(?:pending|todo|tbd|not run)\b", verification, re.IGNORECASE)
):
    raise SystemExit(
        "Venue edit body-limit plan must remain completed with actual verification recorded."
    )
PY

python3 - "$CLIENT_TIMEOUT_PLAN" <<'PY'
import re
import sys
from pathlib import Path

plan = Path(sys.argv[1]).read_text()
frontmatter = plan.split("---", 2)[1]
statuses = re.findall(r"^status: .+$", frontmatter, flags=re.MULTILINE)
verification = plan.split("## Verification Completed\n", 1)[-1]
required = (
    "timeout removal mutation failed",
    "timeout drift mutation failed",
    "unconditional override mutation failed",
    "caller mutation mutation failed",
    "focused test mutation failed",
    "plan evidence mutation failed",
    "hosted pull-request check",
)
if (
    statuses != ["status: completed"]
    or "## Verification Completed\n" not in plan
    or any(item not in verification for item in required)
    or re.search(r"\b(?:pending|todo|tbd|not run)\b", verification, re.IGNORECASE)
):
    raise SystemExit("Foursquare client timeout plan must remain completed with actual verification recorded.")
PY

python3 - "$OAUTH_USER_RESPONSE_PLAN" <<'PY'
import re
import sys
from pathlib import Path

plan = Path(sys.argv[1]).read_text()
frontmatter = plan.split("---", 2)[1]
statuses = re.findall(r"^status: .+$", frontmatter, flags=re.MULTILINE)
required = (
    "seven hostile mutations were rejected",
    "all four Make gates passed",
    "race detector passed",
    "No live OAuth flow",
)
if statuses != ["status: completed"] or any(item not in plan for item in required):
    raise SystemExit("OAuth user response plan must record completed local verification.")
PY

if ! grep -Fq "OAuth user-profile responses require a 2xx status" "$ROOT_DIR/README.md" ||
  ! grep -Fq "OAuth user-profile responses should reject non-2xx" "$ROOT_DIR/SECURITY.md" ||
  ! grep -Fq "OAuth user-profile responses require 2xx status" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "over-1-MiB OAuth user-profile responses" "$ROOT_DIR/CHANGES.md" ||
  ! grep -Fq "Reject non-2xx OAuth user-profile responses" "$ROOT_DIR/AGENTS.md"; then
  printf '%s\n' "Project docs must preserve OAuth user response boundaries." >&2
  exit 1
fi

printf '%s\n' "fsq-go-explore Go baseline checks passed."
