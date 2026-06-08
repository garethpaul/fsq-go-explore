#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
PLAN="$ROOT_DIR/docs/plans/2026-06-08-fsq-go-explore-go-baseline.md"

require_file() {
  path=$1
  if [ ! -f "$ROOT_DIR/$path" ]; then
    printf '%s\n' "Required file missing: $path" >&2
    exit 1
  fi
}

for path in \
  ".gitignore" \
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
  "edit.go" \
  "main.go" \
  "pages.go" \
  "search.go" \
  "search_test.go" \
  "fsq/api.go" \
  "fsq/api_test.go" \
  "fsq/keys.go" \
  "fsq/keys_test.go" \
  "limiter/limiter.go" \
  "docs/plans/2026-06-08-fsq-go-explore-go-baseline.md"; do
  require_file "$path"
done

if command -v go >/dev/null 2>&1; then
  unformatted=$(find "$ROOT_DIR" -name '*.go' -not -path "$ROOT_DIR/.git/*" -print | xargs gofmt -l)
  if [ -n "$unformatted" ]; then
    printf '%s\n' "Go files need gofmt:" >&2
    printf '%s\n' "$unformatted" >&2
    exit 1
  fi
  (cd "$ROOT_DIR" && go test ./...)
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

if ! grep -Fq "go test ./..." "$ROOT_DIR/README.md" ||
  ! grep -Fq "make check" "$ROOT_DIR/README.md" ||
  ! grep -Fq "FSQ_CLIENT_ID" "$ROOT_DIR/README.md"; then
  printf '%s\n' "README must document Go verification and Foursquare env configuration." >&2
  exit 1
fi

if ! grep -Fq "scripts/check-baseline.sh" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "cache-key generation" "$ROOT_DIR/VISION.md" ||
  ! grep -Fq "Go module" "$ROOT_DIR/VISION.md"; then
  printf '%s\n' "VISION must describe the current Go baseline." >&2
  exit 1
fi

if ! grep -Fq "status: completed" "$PLAN"; then
  printf '%s\n' "Plan must be marked completed." >&2
  exit 1
fi

printf '%s\n' "fsq-go-explore Go baseline checks passed."
