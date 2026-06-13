package fsq

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
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
	return testResponseWithStatus(http.StatusOK, body)
}

func testResponseWithStatus(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestSuccessfulFoursquareStatusAcceptsOnly2xx(t *testing.T) {
	tests := []struct {
		status int
		want   bool
	}{
		{status: http.StatusContinue, want: false},
		{status: http.StatusOK, want: true},
		{status: 299, want: true},
		{status: http.StatusMultipleChoices, want: false},
	}

	for _, tt := range tests {
		if got := successfulFoursquareStatus(tt.status); got != tt.want {
			t.Errorf("successfulFoursquareStatus(%d) = %t, want %t", tt.status, got, tt.want)
		}
	}
}

func TestNewFoursquareServiceDefaultsClientTimeout(t *testing.T) {
	service := NewFoursquareService(&FoursquareConfig{})

	if service.Config.Client.Timeout != foursquareRequestTimeout {
		t.Fatalf("client timeout = %s, want %s", service.Config.Client.Timeout, foursquareRequestTimeout)
	}
}

func TestNewFoursquareServicePreservesExplicitClientTimeout(t *testing.T) {
	explicitTimeout := 3 * time.Second
	service := NewFoursquareService(&FoursquareConfig{
		Client: http.Client{Timeout: explicitTimeout},
	})

	if service.Config.Client.Timeout != explicitTimeout {
		t.Fatalf("client timeout = %s, want %s", service.Config.Client.Timeout, explicitTimeout)
	}
}

func TestNewFoursquareServiceDoesNotMutateCallerConfig(t *testing.T) {
	config := &FoursquareConfig{}
	service := NewFoursquareService(config)

	if config.Client.Timeout != 0 {
		t.Fatalf("caller client timeout = %s, want zero", config.Client.Timeout)
	}
	if service.Config == config {
		t.Fatal("service config aliases caller config")
	}
}

func TestSearchRejectsNonSuccessResponseBeforeDecode(t *testing.T) {
	service := NewFoursquareService(&FoursquareConfig{
		ClientId:     "client-id",
		ClientSecret: "client-secret",
		Version:      "20260613",
		Client: http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return testResponseWithStatus(
				http.StatusInternalServerError,
				`{"response":{"venues":[{"id":"must-not-decode"}]}}`,
			), nil
		})},
	})

	response := service.Search(&VenueSearchRequest{Near: "San Francisco", Query: "coffee"})

	if len(response.Venues) != 0 {
		t.Fatalf("Search venues = %#v, want empty result for non-2xx response", response.Venues)
	}
}

func TestVenueDetailsRejectsNonSuccessResponseBeforeDecode(t *testing.T) {
	service := NewFoursquareService(&FoursquareConfig{
		AccessToken: "token",
		Version:     "20260613",
		Client: http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return testResponseWithStatus(
				http.StatusBadGateway,
				`{"response":{"venue":{"id":"must-not-decode"}}}`,
			), nil
		})},
	})

	response := service.VenueDetails("venue-1")

	if response.Venue.ID != "" {
		t.Fatalf("VenueDetails venue ID = %q, want empty result for non-2xx response", response.Venue.ID)
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
