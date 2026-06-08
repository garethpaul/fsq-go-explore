// Foursquare GoLang SDK
package fsq

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

// Setup CONSTANTS
const (
	SEARCH_URL = "https://api.foursquare.com/v2/venues/search?"
	VENUE_URL  = "https://api.foursquare.com/v2/venues/"
)

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
	svc := &FoursquareService{Config: config}
	return svc
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

	if r.StatusCode >= http.StatusBadRequest {
		log.Printf("foursquare search request returned status=%d", r.StatusCode)
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
	requestURL := VENUE_URL + url.PathEscape(id) + "?" + params.Encode()
	r, err := client.Get(requestURL)

	if err != nil {
		log.Print("foursquare venue details request failed")
		return venue
	}
	defer r.Body.Close()

	if r.StatusCode >= http.StatusBadRequest {
		log.Printf("foursquare venue details request returned status=%d", r.StatusCode)
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
	requestURL := VENUE_URL + url.PathEscape(venueId) + "/proposeedit?" + params.Encode()
	client := foursquareConfig.Client
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBufferString(vals.Encode()))
	if err != nil {
		log.Printf("foursquare venue edit request build failed: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Print("foursquare venue edit request failed")
		return
	}
	defer resp.Body.Close()
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		log.Printf("foursquare venue edit response drain failed: %v", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Printf("foursquare venue edit request returned status=%d", resp.StatusCode)
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
	data, err := io.ReadAll(body)
	if err != nil {
		return err
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
