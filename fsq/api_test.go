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

type trackingReadCloser struct {
	readCalls int
}

func (r *trackingReadCloser) Read([]byte) (int, error) {
	r.readCalls++
	return 0, io.EOF
}

func (r *trackingReadCloser) Close() error { return nil }

type unboundedReadCloser struct {
	bytesRead int
}

func (r *unboundedReadCloser) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 'x'
	}
	r.bytesRead += len(p)
	return len(p), nil
}

func (r *unboundedReadCloser) Close() error { return nil }

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
	response, err := f(req)
	if response != nil && response.Request == nil {
		response.Request = req
	}
	return response, err
}

func testResponse(body string) *http.Response {
	return testResponseWithStatus(http.StatusOK, body)
}

func testResponseWithStatus(status int, body string) *http.Response {
	response := &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
	response.Header.Set("Content-Type", "application/json; charset=utf-8")
	return response
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

func TestFoursquareJSONResponseMediaTypes(t *testing.T) {
	for _, contentType := range []string{"application/json", "application/json; charset=utf-8", "application/vnd.foursquare+json"} {
		response := testResponse(`{"response":{}}`)
		response.Header.Set("Content-Type", contentType)
		if !isFoursquareJSONResponse(response) {
			t.Errorf("isFoursquareJSONResponse(%q) = false, want true", contentType)
		}
	}
	for _, contentType := range []string{"", "text/html", "text/json", "application/octet-stream", "not a type"} {
		response := testResponse(`{"response":{}}`)
		response.Header.Set("Content-Type", contentType)
		if isFoursquareJSONResponse(response) {
			t.Errorf("isFoursquareJSONResponse(%q) = true, want false", contentType)
		}
	}
}

func TestExpectedFoursquareResponseURL(t *testing.T) {
	tests := []struct {
		name         string
		responseURL  string
		expectedPath string
		want         bool
	}{
		{name: "search", responseURL: "https://api.foursquare.com/v2/venues/search?near=Portland", expectedPath: "/v2/venues/search", want: true},
		{name: "escaped venue", responseURL: "https://api.foursquare.com/v2/venues/venue%2F123?v=20260614", expectedPath: "/v2/venues/venue%2F123", want: true},
		{name: "http", responseURL: "http://api.foursquare.com/v2/venues/search", expectedPath: "/v2/venues/search"},
		{name: "different host", responseURL: "https://example.com/v2/venues/search", expectedPath: "/v2/venues/search"},
		{name: "userinfo", responseURL: "https://user@api.foursquare.com/v2/venues/search", expectedPath: "/v2/venues/search"},
		{name: "port", responseURL: "https://api.foursquare.com:443/v2/venues/search", expectedPath: "/v2/venues/search"},
		{name: "different path", responseURL: "https://api.foursquare.com/v2/venues/explore", expectedPath: "/v2/venues/search"},
		{name: "fragment", responseURL: "https://api.foursquare.com/v2/venues/search#result", expectedPath: "/v2/venues/search"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseURL, err := url.Parse(tt.responseURL)
			if err != nil {
				t.Fatalf("url.Parse(%q): %v", tt.responseURL, err)
			}
			response := &http.Response{Request: &http.Request{URL: responseURL}}
			if got := isExpectedFoursquareResponseURL(response, tt.expectedPath); got != tt.want {
				t.Fatalf("isExpectedFoursquareResponseURL(%q, %q) = %t, want %t", tt.responseURL, tt.expectedPath, got, tt.want)
			}
		})
	}
}

func TestFoursquareOperationsRejectUnexpectedFinalURLBeforeRead(t *testing.T) {
	tests := []struct {
		name string
		call func(*FoursquareService)
	}{
		{name: "search", call: func(service *FoursquareService) {
			service.Search(&VenueSearchRequest{Near: "Portland", Query: "coffee"})
		}},
		{name: "venue details", call: func(service *FoursquareService) {
			service.VenueDetails("venue-1")
		}},
		{name: "venue edit", call: func(service *FoursquareService) {
			service.VenueEdit("venue-1", url.Values{"name": []string{"New Name"}})
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := &trackingReadCloser{}
			unexpectedURL, err := url.Parse("https://example.com/v2/venues/search")
			if err != nil {
				t.Fatalf("url.Parse: %v", err)
			}
			service := NewFoursquareService(&FoursquareConfig{
				ClientId: "client-id", ClientSecret: "client-secret", AccessToken: "token", Version: "20260614",
				Client: http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       body,
						Header:     http.Header{"Content-Type": []string{"application/json"}},
						Request:    &http.Request{URL: unexpectedURL},
					}, nil
				})},
			})

			tt.call(service)
			if body.readCalls != 0 {
				t.Fatalf("response body reads = %d, want 0", body.readCalls)
			}
		})
	}
}

func TestSearchRejectsNonJSONResponseBeforeDecode(t *testing.T) {
	service := NewFoursquareService(&FoursquareConfig{
		ClientId: "client-id", ClientSecret: "client-secret", Version: "20260614",
		Client: http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			response := testResponse(`{"response":{"venues":[{"id":"must-not-decode"}]}}`)
			response.Header.Set("Content-Type", "text/html")
			return response, nil
		})},
	})
	response := service.Search(&VenueSearchRequest{Near: "San Francisco", Query: "coffee"})
	if len(response.Venues) != 0 {
		t.Fatalf("Search venues = %#v, want empty result for non-JSON response", response.Venues)
	}
}

func TestVenueDetailsRejectsNonJSONResponseBeforeDecode(t *testing.T) {
	service := NewFoursquareService(&FoursquareConfig{
		AccessToken: "token", Version: "20260614",
		Client: http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			response := testResponse(`{"response":{"venue":{"id":"must-not-decode"}}}`)
			response.Header.Set("Content-Type", "text/plain")
			return response, nil
		})},
	})
	response := service.VenueDetails("venue-1")
	if response.Venue.ID != "" {
		t.Fatalf("VenueDetails venue ID = %q, want empty result for non-JSON response", response.Venue.ID)
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

func TestNewFoursquareServiceRefusesRedirects(t *testing.T) {
	transport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return testResponse(`{"response":{}}`), nil
	})
	originalRedirect := func(req *http.Request, via []*http.Request) error { return nil }
	config := &FoursquareConfig{
		Client: http.Client{
			Transport:     transport,
			Timeout:       3 * time.Second,
			CheckRedirect: originalRedirect,
		},
	}

	service := NewFoursquareService(config)
	if err := service.Config.Client.CheckRedirect(nil, nil); !errors.Is(err, http.ErrUseLastResponse) {
		t.Fatalf("service redirect policy error = %v, want http.ErrUseLastResponse", err)
	}
	if service.Config.Client.Transport == nil || service.Config.Client.Timeout != 3*time.Second {
		t.Fatal("service redirect policy must preserve caller transport and timeout")
	}
	if err := config.Client.CheckRedirect(nil, nil); err != nil {
		t.Fatalf("caller redirect policy was mutated: %v", err)
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

func TestVenueEditRejectsNonSuccessBeforeRead(t *testing.T) {
	for _, status := range []int{http.StatusContinue, http.StatusMultipleChoices, http.StatusInternalServerError} {
		body := &trackingReadCloser{}
		service := NewFoursquareService(&FoursquareConfig{
			AccessToken: "token",
			Version:     "20260614",
			Client: http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: status, Body: body, Header: make(http.Header)}, nil
			})},
		})

		service.VenueEdit("venue-1", url.Values{"name": []string{"New Name"}})
		if body.readCalls != 0 {
			t.Fatalf("VenueEdit status %d body reads = %d, want 0", status, body.readCalls)
		}
	}
}

func TestVenueEditBoundsSuccessfulResponseBody(t *testing.T) {
	body := &unboundedReadCloser{}
	service := NewFoursquareService(&FoursquareConfig{
		AccessToken: "token",
		Version:     "20260614",
		Client: http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusOK, Body: body, Header: make(http.Header)}, nil
		})},
	})

	service.VenueEdit("venue-1", url.Values{"name": []string{"New Name"}})
	if body.bytesRead != maxFoursquareResponseBytes+1 {
		t.Fatalf("VenueEdit body bytes read = %d, want %d", body.bytesRead, maxFoursquareResponseBytes+1)
	}
}

func TestDiscardFoursquareResponseAcceptsExactLimit(t *testing.T) {
	if err := discardFoursquareResponse(strings.NewReader(strings.Repeat("x", maxFoursquareResponseBytes))); err != nil {
		t.Fatalf("discardFoursquareResponse exact limit: %v", err)
	}
}

func TestDiscardFoursquareResponseRejectsOversizeBody(t *testing.T) {
	err := discardFoursquareResponse(strings.NewReader(strings.Repeat("x", maxFoursquareResponseBytes+1)))
	if !errors.Is(err, errFoursquareResponseTooLarge) {
		t.Fatalf("discardFoursquareResponse error = %v, want %v", err, errFoursquareResponseTooLarge)
	}
}

func TestDiscardFoursquareResponsePreservesReadError(t *testing.T) {
	err := discardFoursquareResponse(&failingReader{})
	if !errors.Is(err, errTestReadFailure) {
		t.Fatalf("discardFoursquareResponse error = %v, want %v", err, errTestReadFailure)
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
