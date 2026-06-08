package fsq

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

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
