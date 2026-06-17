// Foursquare GoLang SDK
package fsq

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/oauth2"
)

// Setup CONSTANTS
const (
	SEARCH_URL                 = "https://api.foursquare.com/v2/venues/search?"
	VENUE_URL                  = "https://api.foursquare.com/v2/venues/"
	foursquareAPIHost          = "api.foursquare.com"
	foursquareSearchPath       = "/v2/venues/search"
	foursquareVenuePathPrefix  = "/v2/venues/"
	maxFoursquareResponseBytes = 2 * 1024 * 1024
	foursquareRequestTimeout   = 10 * time.Second
)

var errFoursquareResponseTooLarge = errors.New("foursquare response body exceeds 2 MiB")
var errFoursquareResponseInvalidUTF8 = errors.New("foursquare response body was not valid UTF-8")

func RefuseRedirect(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

func successfulFoursquareStatus(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}

func isExpectedFoursquareResponseURL(response *http.Response, expectedEscapedPath string) bool {
	if response == nil || response.Request == nil || response.Request.URL == nil {
		return false
	}

	responseURL := response.Request.URL
	return responseURL.Scheme == "https" &&
		strings.EqualFold(responseURL.Hostname(), foursquareAPIHost) &&
		responseURL.User == nil &&
		responseURL.Port() == "" &&
		responseURL.EscapedPath() == expectedEscapedPath &&
		responseURL.Fragment == ""
}

func isFoursquareJSONResponse(response *http.Response) bool {
	mediaType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return false
	}
	mediaType = strings.ToLower(mediaType)
	return mediaType == "application/json" ||
		(strings.HasPrefix(mediaType, "application/") && strings.HasSuffix(mediaType, "+json"))
}

func discardFoursquareResponse(body io.Reader) error {
	written, err := io.Copy(io.Discard, io.LimitReader(body, maxFoursquareResponseBytes+1))
	if err != nil {
		return err
	}
	if written > maxFoursquareResponseBytes {
		return errFoursquareResponseTooLarge
	}
	return nil
}

// Struct for FourceService to wrap around requests.
type FoursquareService struct {
	Config *FoursquareConfig
}

// Struct for configuration file.
type FoursquareConfig struct {
	ClientId     string
	ClientSecret string
	AccessToken  string
	Client       http.Client
	Version      string
	AuthConfig   *oauth2.Config
}

// Create a new fource service with a given config file.
func NewFoursquareService(config *FoursquareConfig) *FoursquareService {
	serviceConfig := *config
	client := config.Client
	if client.Timeout <= 0 {
		client.Timeout = foursquareRequestTimeout
	}
	client.CheckRedirect = RefuseRedirect
	serviceConfig.Client = client
	return &FoursquareService{Config: &serviceConfig}
}

// See https://developer.foursquare.com/docs/venues/search
func (fsqs *FoursquareService) Search(vsr *VenueSearchRequest) (resp *VenueSearchResponse) {
	venues := new(VenueSearchResponse)
	foursquareConfig := fsqs.Config

	// Setup Params
	params := foursquareConfig.clientParams()
	params.Set("near", vsr.Near)
	params.Set("query", vsr.Query)
	params.Set("limit", "100")

	client := foursquareConfig.Client
	url := SEARCH_URL + params.Encode()
	r, err := client.Get(url)
	if err != nil {
		log.Print("foursquare search request failed")
		return venues
	}
	defer r.Body.Close()

	if !successfulFoursquareStatus(r.StatusCode) {
		log.Printf("foursquare search request returned status=%d", r.StatusCode)
		return venues
	}
	if !isExpectedFoursquareResponseURL(r, foursquareSearchPath) {
		log.Print("foursquare search response final URL was rejected")
		return venues
	}
	if !isFoursquareJSONResponse(r) {
		log.Print("foursquare search response content type was rejected")
		return venues
	}
	if err := decodeFoursquareResponse(r.Body, venues); err != nil {
		log.Printf("foursquare search response decode failed: %v", err)
	}
	return venues
}

// Details gets all the data for a venue
// https://developer.foursquare.com/docs/venues/venues
func (fsqs *FoursquareService) VenueDetails(id string) (resp *VenueResponse) {
	venue := new(VenueResponse)

	foursquareConfig := fsqs.Config

	params := foursquareConfig.userParams()
	client := foursquareConfig.Client
	venuePath := foursquareVenuePathPrefix + url.PathEscape(id)
	requestURL := VENUE_URL + url.PathEscape(id) + "?" + params.Encode()
	r, err := client.Get(requestURL)

	if err != nil {
		log.Print("foursquare venue details request failed")
		return venue
	}
	defer r.Body.Close()

	if !successfulFoursquareStatus(r.StatusCode) {
		log.Printf("foursquare venue details request returned status=%d", r.StatusCode)
		return venue
	}
	if !isExpectedFoursquareResponseURL(r, venuePath) {
		log.Print("foursquare venue details response final URL was rejected")
		return venue
	}
	if !isFoursquareJSONResponse(r) {
		log.Print("foursquare venue details response content type was rejected")
		return venue
	}
	if err := decodeFoursquareResponse(r.Body, venue); err != nil {
		log.Printf("foursquare venue details response decode failed: %v", err)
	}
	return venue
}

// ProposeEdit
// More details via https://developer.foursquare.com/docs/venues/proposeedit
func (fsqs *FoursquareService) VenueEdit(venueId string, vals url.Values) {
	foursquareConfig := fsqs.Config
	params := foursquareConfig.userParams()
	venueEditPath := foursquareVenuePathPrefix + url.PathEscape(venueId) + "/proposeedit"
	requestURL := VENUE_URL + url.PathEscape(venueId) + "/proposeedit?" + params.Encode()
	client := foursquareConfig.Client
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBufferString(vals.Encode()))
	if err != nil {
		log.Print("foursquare venue edit request build failed")
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Print("foursquare venue edit request failed")
		return
	}
	defer resp.Body.Close()
	if !successfulFoursquareStatus(resp.StatusCode) {
		log.Printf("foursquare venue edit request returned status=%d", resp.StatusCode)
		return
	}
	if !isExpectedFoursquareResponseURL(resp, venueEditPath) {
		log.Print("foursquare venue edit response final URL was rejected")
		return
	}
	if err := discardFoursquareResponse(resp.Body); err != nil {
		log.Printf("foursquare venue edit response discard failed: %v", err)
	}
}

// Get the url params
func (fsqs *FoursquareConfig) userParams() url.Values {
	params := url.Values{}
	params.Set("oauth_token", fsqs.AccessToken)
	params.Set("v", fsqs.Version)
	return params
}

// Get the url params
func (fsqs *FoursquareConfig) clientParams() url.Values {
	params := url.Values{}
	params.Set("client_id", fsqs.ClientId)
	params.Set("client_secret", fsqs.ClientSecret)
	params.Set("v", fsqs.Version)
	return params
}

func decodeFoursquareResponse(body io.Reader, target interface{}) error {
	data, err := io.ReadAll(io.LimitReader(body, maxFoursquareResponseBytes+1))
	if err != nil {
		return err
	}
	if len(data) > maxFoursquareResponseBytes {
		return errFoursquareResponseTooLarge
	}
	if !utf8.Valid(data) {
		return errFoursquareResponseInvalidUTF8
	}
	response := new(Response)
	if err := json.Unmarshal(data, response); err != nil {
		return err
	}
	if len(response.Response) == 0 {
		return nil
	}
	return json.Unmarshal(response.Response, target)
}
