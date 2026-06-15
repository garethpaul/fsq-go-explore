package app

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestGetHTTPClientRefusesRedirects(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://example.com/login", nil)
	client := getHttpClient(req)

	if client.Transport == nil {
		t.Fatal("expected App Engine URL fetch transport")
	}
	if err := client.CheckRedirect(nil, nil); !errors.Is(err, http.ErrUseLastResponse) {
		t.Fatalf("redirect policy error = %v, want http.ErrUseLastResponse", err)
	}
}

func TestGetHTTPClientBoundsOAuthUserRequests(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://example.com/login", nil)
	client := getHttpClient(req)

	if client.Timeout != 10*time.Second {
		t.Fatalf("timeout = %s, want 10s", client.Timeout)
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

func TestDecodeOAuthUserResponseRejectsUnexpectedFinalURLBeforeRead(t *testing.T) {
	for name, responseURL := range map[string]string{
		"scheme":   "http://api.foursquare.com/v2/users/self",
		"host":     "https://example.com/v2/users/self",
		"userinfo": "https://user@api.foursquare.com/v2/users/self",
		"port":     "https://api.foursquare.com:443/v2/users/self",
		"path":     "https://api.foursquare.com/v2/users/other",
		"fragment": "https://api.foursquare.com/v2/users/self#profile",
	} {
		t.Run(name, func(t *testing.T) {
			body := &trackingReadCloser{reader: strings.NewReader(validOAuthUserResponse())}
			response := oauthUserResponse(http.StatusOK, "application/json", body)
			response.Request = httptest.NewRequest(http.MethodGet, responseURL, nil)

			_, err := decodeOAuthUserResponse(response)
			if !errors.Is(err, errOAuthUserResponseOrigin) {
				t.Fatalf("error = %v, want %v", err, errOAuthUserResponseOrigin)
			}
			if body.readCalls != 0 {
				t.Fatalf("body reads = %d, want zero", body.readCalls)
			}
		})
	}
}

func TestDecodeOAuthUserResponseAcceptsExpectedFinalURL(t *testing.T) {
	response := oauthUserResponse(http.StatusOK, "application/json", io.NopCloser(strings.NewReader(validOAuthUserResponse())))
	response.Request = httptest.NewRequest(http.MethodGet, "https://API.FOURSQUARE.COM/v2/users/self?oauth_token=redacted&v=20170101", nil)

	if _, err := decodeOAuthUserResponse(response); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeOAuthUserResponseRejectsNonJSONBeforeRead(t *testing.T) {
	for _, contentType := range []string{"", "text/html", "application/jsonp", "application/json, text/html", "not a media type"} {
		body := &trackingReadCloser{reader: strings.NewReader(validOAuthUserResponse())}
		response := oauthUserResponse(http.StatusOK, contentType, body)

		_, err := decodeOAuthUserResponse(response)
		if !errors.Is(err, errOAuthUserResponseMediaType) {
			t.Fatalf("content type %q error = %v, want %v", contentType, err, errOAuthUserResponseMediaType)
		}
		if body.readCalls != 0 {
			t.Fatalf("content type %q body reads = %d, want zero", contentType, body.readCalls)
		}
	}
}

func TestDecodeOAuthUserResponseRejectsDuplicateContentTypeBeforeRead(t *testing.T) {
	body := &trackingReadCloser{reader: strings.NewReader(validOAuthUserResponse())}
	response := oauthUserResponse(http.StatusOK, "application/json", body)
	response.Header.Add("Content-Type", "text/html")

	_, err := decodeOAuthUserResponse(response)
	if !errors.Is(err, errOAuthUserResponseMediaType) {
		t.Fatalf("error = %v, want %v", err, errOAuthUserResponseMediaType)
	}
	if body.readCalls != 0 {
		t.Fatalf("body reads = %d, want zero", body.readCalls)
	}
}

func TestDecodeOAuthUserResponseAcceptsSingleStructuredJSONContentType(t *testing.T) {
	response := oauthUserResponse(
		http.StatusOK,
		"application/problem+json; charset=utf-8",
		io.NopCloser(strings.NewReader(validOAuthUserResponse())),
	)

	if _, err := decodeOAuthUserResponse(response); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeOAuthUserResponseAcceptsExactLimit(t *testing.T) {
	body := validOAuthUserResponse()
	body += strings.Repeat(" ", maxOAuthUserResponseBytes-len(body))

	user, err := decodeOAuthUserResponse(oauthUserResponse(
		http.StatusOK,
		"application/problem+json; charset=utf-8",
		io.NopCloser(strings.NewReader(body)),
	))
	if err != nil {
		t.Fatal(err)
	}
	if user.User.ID != "user-1" {
		t.Fatalf("user ID = %q, want user-1", user.User.ID)
	}
}

func TestDecodeOAuthUserResponseRejectsOversizeBody(t *testing.T) {
	_, err := decodeOAuthUserResponse(oauthUserResponse(
		http.StatusOK,
		"application/json",
		io.NopCloser(strings.NewReader(strings.Repeat("x", maxOAuthUserResponseBytes+1))),
	))
	if !errors.Is(err, errOAuthUserResponseTooLarge) {
		t.Fatalf("error = %v, want %v", err, errOAuthUserResponseTooLarge)
	}
}

func TestDecodeOAuthUserResponsePreservesReadError(t *testing.T) {
	_, err := decodeOAuthUserResponse(oauthUserResponse(
		http.StatusOK,
		"application/json",
		io.NopCloser(failingOAuthProfileReader{}),
	))
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
			_, err := decodeOAuthUserResponse(oauthUserResponse(
				http.StatusOK,
				"application/json",
				io.NopCloser(strings.NewReader(body)),
			))
			if err == nil {
				t.Fatal("error = nil, want malformed JSON rejection")
			}
		})
	}
}

func TestDecodeOAuthUserResponseRejectsDuplicateJSONMembers(t *testing.T) {
	for name, body := range map[string]string{
		"response":     `{"response":{"user":{"id":"user-1"}},"response":{"user":{"id":"user-2"}}}`,
		"user":         `{"response":{"user":{"id":"user-1"},"user":{"id":"user-2"}}}`,
		"id":           `{"response":{"user":{"id":"user-1","id":"user-2"}}}`,
		"array object": `{"response":{"user":{"id":"user-1"},"items":[{"value":1,"value":2}]}}`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := decodeOAuthUserResponse(oauthUserResponse(
				http.StatusOK,
				"application/json",
				io.NopCloser(strings.NewReader(body)),
			))
			if !errors.Is(err, errOAuthUserResponseDuplicateKey) {
				t.Fatalf("error = %v, want %v", err, errOAuthUserResponseDuplicateKey)
			}
		})
	}
}

func TestDecodeOAuthUserResponseAcceptsUniqueUnknownNestedMembers(t *testing.T) {
	user, err := decodeOAuthUserResponse(oauthUserResponse(
		http.StatusOK,
		"application/json",
		io.NopCloser(strings.NewReader(
			`{"response":{"user":{"id":"user-1"},"items":[{"value":1},{"value":2}]}}`,
		)),
	))
	if err != nil {
		t.Fatal(err)
	}
	if user.User.ID != "user-1" {
		t.Fatalf("user ID = %q, want user-1", user.User.ID)
	}
}

func TestDecodeOAuthUserResponseRejectsInvalidUserIDs(t *testing.T) {
	for name, idField := range map[string]string{
		"missing":             ``,
		"empty":               `"id":""`,
		"whitespace only":     `"id":" \t\n"`,
		"leading whitespace":  `"id":"\u00a0user-1"`,
		"trailing whitespace": `"id":"user-1 "`,
	} {
		t.Run(name, func(t *testing.T) {
			separator := ""
			if idField != "" {
				separator = ","
			}
			body := fmt.Sprintf(
				`{"response":{"user":{%s%s"firstName":"Example"}}}`,
				idField,
				separator,
			)
			if _, err := decodeOAuthUserResponse(oauthUserResponse(
				http.StatusOK,
				"application/json",
				io.NopCloser(strings.NewReader(body)),
			)); !errors.Is(err, errOAuthUserResponseIdentity) {
				t.Fatalf("error = %v, want %v", err, errOAuthUserResponseIdentity)
			}
		})
	}
}

func TestDecodeOAuthUserResponseAllowsEmptyDisplayName(t *testing.T) {
	user, err := decodeOAuthUserResponse(oauthUserResponse(
		http.StatusOK,
		"application/json",
		io.NopCloser(strings.NewReader(`{"response":{"user":{"id":"user-1"}}}`)),
	))
	if err != nil {
		t.Fatal(err)
	}
	if user.User.ID != "user-1" || user.User.FirstName != "" {
		t.Fatalf("user = %#v, want canonical ID and empty display name", user.User)
	}
}

func validOAuthUserResponse() string {
	return fmt.Sprintf(`{"response":{"user":{"id":%q,"firstName":"Example"}}}`, "user-1")
}

func oauthUserResponse(status int, contentType string, body io.ReadCloser) *http.Response {
	response := &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       body,
		Request:    httptest.NewRequest(http.MethodGet, "https://api.foursquare.com/v2/users/self", nil),
	}
	if contentType != "" {
		response.Header.Set("Content-Type", contentType)
	}
	return response
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
