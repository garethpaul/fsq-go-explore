package app

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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
