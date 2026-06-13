package fsq

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type failingReader struct {
	read bool
}

var errTestReadFailure = errors.New("read failed")

func (r *failingReader) Read(p []byte) (int, error) {
	if r.read {
		return 0, errTestReadFailure
	}
	r.read = true
	return copy(p, `{"response":`), nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func testResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestVenueDetailsEscapesVenueID(t *testing.T) {
	var gotPath string
	service := NewFoursquareService(&FoursquareConfig{
		AccessToken: "token",
		Version:     "20260608",
		Client: http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotPath = req.URL.EscapedPath()
			return testResponse(`{"response":{"venue":{"id":"venue/123"}}}`), nil
		})},
	})

	service.VenueDetails("venue/123")

	if !strings.Contains(gotPath, "venue%2F123") {
		t.Fatalf("VenueDetails path = %q, want escaped venue id", gotPath)
	}
}

func TestVenueEditEscapesVenueIDAndSendsForm(t *testing.T) {
	var gotPath string
	var gotContentType string
	var gotBody string
	service := NewFoursquareService(&FoursquareConfig{
		AccessToken: "token",
		Version:     "20260608",
		Client: http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			gotPath = req.URL.EscapedPath()
			gotContentType = req.Header.Get("Content-Type")
			body, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("reading request body: %v", err)
			}
			gotBody = string(body)
			return testResponse(`{"response":{}}`), nil
		})},
	})

	service.VenueEdit("venue/123", url.Values{"name": []string{"New Name"}})

	if !strings.Contains(gotPath, "venue%2F123/proposeedit") {
		t.Fatalf("VenueEdit path = %q, want escaped venue id and proposeedit suffix", gotPath)
	}
	if gotContentType != "application/x-www-form-urlencoded" {
		t.Fatalf("Content-Type = %q, want form encoding", gotContentType)
	}
	if gotBody != "name=New+Name" {
		t.Fatalf("body = %q, want encoded form", gotBody)
	}
}

func TestDecodeFoursquareResponseAcceptsExactLimit(t *testing.T) {
	body := `{"response":{}}`
	body += strings.Repeat(" ", maxFoursquareResponseBytes-len(body))

	if err := decodeFoursquareResponse(strings.NewReader(body), &VenueSearchResponse{}); err != nil {
		t.Fatalf("decodeFoursquareResponse exact limit: %v", err)
	}
}

func TestDecodeFoursquareResponseRejectsOversizeBody(t *testing.T) {
	body := strings.Repeat(" ", maxFoursquareResponseBytes+1)

	err := decodeFoursquareResponse(strings.NewReader(body), &VenueSearchResponse{})
	if !errors.Is(err, errFoursquareResponseTooLarge) {
		t.Fatalf("decodeFoursquareResponse error = %v, want %v", err, errFoursquareResponseTooLarge)
	}
}

func TestDecodeFoursquareResponsePreservesReadError(t *testing.T) {
	err := decodeFoursquareResponse(&failingReader{}, &VenueSearchResponse{})
	if !errors.Is(err, errTestReadFailure) {
		t.Fatalf("decodeFoursquareResponse error = %v, want %v", err, errTestReadFailure)
	}
}

func TestDecodeFoursquareResponseRejectsEmptyBody(t *testing.T) {
	if err := decodeFoursquareResponse(strings.NewReader(""), &VenueSearchResponse{}); err == nil {
		t.Fatal("decodeFoursquareResponse empty body error = nil, want JSON decode error")
	}
}

func TestDecodeFoursquareResponseRejectsMalformedJSON(t *testing.T) {
	if err := decodeFoursquareResponse(strings.NewReader(`{"response":`), &VenueSearchResponse{}); err == nil {
		t.Fatal("decodeFoursquareResponse malformed JSON error = nil, want JSON decode error")
	}
}
