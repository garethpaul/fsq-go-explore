package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/garethpaul/fsq-go-explore/fsq"
)

var errOAuthProfileRead = errors.New("oauth profile read failed")

type trackingReadCloser struct {
	reader    io.Reader
	readCalls int
}

func (r *trackingReadCloser) Read(p []byte) (int, error) {
	r.readCalls++
	return r.reader.Read(p)
}

func (r *trackingReadCloser) Close() error { return nil }

type failingOAuthProfileReader struct{}

func (failingOAuthProfileReader) Read([]byte) (int, error) {
	return 0, errOAuthProfileRead
}

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

func TestDecodeOAuthUserResponseRejectsNonSuccessBeforeRead(t *testing.T) {
	for _, status := range []int{http.StatusContinue, http.StatusMovedPermanently, http.StatusBadRequest, http.StatusBadGateway} {
		body := &trackingReadCloser{reader: strings.NewReader(validOAuthUserResponse())}
		response := &http.Response{StatusCode: status, Body: body}

		_, err := decodeOAuthUserResponse(response)
		if !errors.Is(err, errOAuthUserResponseStatus) {
			t.Fatalf("status %d error = %v, want %v", status, err, errOAuthUserResponseStatus)
		}
		if body.readCalls != 0 {
			t.Fatalf("status %d body reads = %d, want zero", status, body.readCalls)
		}
	}
}

func TestDecodeOAuthUserResponseAcceptsExactLimit(t *testing.T) {
	body := validOAuthUserResponse()
	body += strings.Repeat(" ", maxOAuthUserResponseBytes-len(body))

	user, err := decodeOAuthUserResponse(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if user.User.ID != "user-1" {
		t.Fatalf("user ID = %q, want user-1", user.User.ID)
	}
}

func TestDecodeOAuthUserResponseRejectsOversizeBody(t *testing.T) {
	_, err := decodeOAuthUserResponse(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(strings.Repeat("x", maxOAuthUserResponseBytes+1))),
	})
	if !errors.Is(err, errOAuthUserResponseTooLarge) {
		t.Fatalf("error = %v, want %v", err, errOAuthUserResponseTooLarge)
	}
}

func TestDecodeOAuthUserResponsePreservesReadError(t *testing.T) {
	_, err := decodeOAuthUserResponse(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(failingOAuthProfileReader{}),
	})
	if !errors.Is(err, errOAuthProfileRead) {
		t.Fatalf("error = %v, want %v", err, errOAuthProfileRead)
	}
}

func TestDecodeOAuthUserResponseRejectsMalformedPayloads(t *testing.T) {
	for name, body := range map[string]string{
		"wrapper": `{"response":`,
		"user":    `{"response":"invalid"}`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := decodeOAuthUserResponse(&http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
			})
			if err == nil {
				t.Fatal("error = nil, want malformed JSON rejection")
			}
		})
	}
}

func validOAuthUserResponse() string {
	return fmt.Sprintf(`{"response":{"user":{"id":%q,"firstName":"Example"}}}`, "user-1")
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

func TestLoginProtectRejectsMalformedAuthCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/edit", nil)
	req.AddCookie(&http.Cookie{Name: "fsq", Value: "not-a-user-cache-key"})
	rr := httptest.NewRecorder()
	called := false

	handler := LoginProtect(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	handler(rr, req)

	if called {
		t.Fatal("protected handler was called for malformed auth cookie")
	}
	if rr.Code != http.StatusTemporaryRedirect {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusTemporaryRedirect)
	}
	if location := rr.Header().Get("Location"); location != "/login" {
		t.Fatalf("redirect location = %q, want /login", location)
	}
}

func TestLoginProtectAllowsGeneratedUserCacheKey(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/edit", nil)
	req.AddCookie(&http.Cookie{
		Name:  "fsq",
		Value: fsq.GetUserKey(&fsq.FoursquareUser{ID: "user-1", Name: "Example", AccessToken: "token"}),
	})
	rr := httptest.NewRecorder()
	called := false

	handler := LoginProtect(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})
	handler(rr, req)

	if !called {
		t.Fatal("protected handler was not called for generated auth cache key")
	}
	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}
