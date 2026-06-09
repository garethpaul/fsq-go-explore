package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/garethpaul/fsq-go-explore/fsq"
)

func TestNewOAuthStateReturnsDistinctOpaqueValues(t *testing.T) {
	first, err := newOAuthState()
	if err != nil {
		t.Fatal(err)
	}
	second, err := newOAuthState()
	if err != nil {
		t.Fatal(err)
	}

	if first == "" || second == "" {
		t.Fatal("expected OAuth states to be populated")
	}
	if first == second {
		t.Fatal("expected distinct OAuth states")
	}
}

func TestRedirectRejectsMissingAuthorizationCodeBeforeExchange(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/redirect?state=state-1", nil)
	req.AddCookie(&http.Cookie{Name: oauthStateCookieName, Value: "state-1"})
	rr := httptest.NewRecorder()

	Redirect(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
	if location := rr.Header().Get("Location"); location != "" {
		t.Fatalf("redirect location = %q, want none", location)
	}
	if !strings.Contains(rr.Body.String(), "missing authorization code") {
		t.Fatalf("body = %q, want missing authorization code", rr.Body.String())
	}
}

func TestValidUserCacheKeyAcceptsGeneratedUserKeys(t *testing.T) {
	key := fsq.GetUserKey(&fsq.FoursquareUser{ID: "user-1", Name: "Example", AccessToken: "token"})

	if !validUserCacheKey(key) {
		t.Fatalf("validUserCacheKey(%q) = false, want true", key)
	}
}

func TestGetAccessTokenRejectsMalformedCacheKeysBeforeLookup(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, key := range []string{
		"",
		"search:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		"user:not-hex",
		"user:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeg",
		"user:" + strings.Repeat("a", 300),
	} {
		if token := getAccessToken(req, key); token != "" {
			t.Fatalf("getAccessToken(%q) = %q, want empty token", key, token)
		}
	}
}
