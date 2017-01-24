// Foursquare GoLang SDK
package fsq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
  "golang.org/x/oauth2"
  "log"
  "bytes"
)


// Setup CONSTANTS
const (
	SEARCH_URL = "https://api.foursquare.com/v2/venues/search?"
  VENUE_URL = "https://api.foursquare.com/v2/venues/"
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
  AuthConfig  *oauth2.Config
}

// Create a new fource service with a given config file.
func NewFoursquareService(config *FoursquareConfig) *FoursquareService {
	svc := &FoursquareService{Config: config}
	return svc
}

// See https://developer.foursquare.com/docs/venues/search
func (fsqs *FoursquareService) Search(vsr *VenueSearchRequest) (resp *VenueSearchResponse) {
	venues := new(VenueSearchResponse)
	response := new(Response)
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
		fmt.Println(err)
	}
	defer r.Body.Close()

	data, _ := ioutil.ReadAll(r.Body)

	// More details via https://developer.foursquare.com/overview/responses
	json.Unmarshal([]byte((string(data))), &response)
	// Map the response to VenueSearchRespones
	json.Unmarshal(response.Response, venues)
	return venues
}

// Details gets all the data for a venue
// https://developer.foursquare.com/docs/venues/venues
func (fsqs *FoursquareService) VenueDetails(id string) (resp *VenueResponse) {
	response := new(Response)
	venue := new(VenueResponse)

  foursquareConfig := fsqs.Config

  params := foursquareConfig.userParams()
  client := foursquareConfig.Client
  url := VENUE_URL + id + "?" +params.Encode()

  log.Print(url)
  r, err := client.Get(url)

  if err != nil {
    fmt.Println(err)
  }
  defer r.Body.Close()

  data, _ := ioutil.ReadAll(r.Body)

  // More details via https://developer.foursquare.com/overview/responses
  json.Unmarshal([]byte((string(data))), &response)
  // Map the response to VenueSearchRespones
  json.Unmarshal(response.Response, venue)
	return venue
}

// ProposeEdit
// More details via https://developer.foursquare.com/docs/venues/proposeedit
func (fsqs *FoursquareService) VenueEdit(venueId string, vals url.Values) {
  foursquareConfig := fsqs.Config
  params := foursquareConfig.userParams()
  url := VENUE_URL + venueId + "/proposeedit?" + params.Encode()
  client := foursquareConfig.Client
  str := vals.Encode()
  log.Print(str)
  req, err := http.NewRequest("POST", url, bytes.NewBufferString(str))
  resp, err := client.Do(req)
  if err != nil {
    log.Print(err)
    panic(err)
  }
  defer resp.Body.Close()
  data, _ := ioutil.ReadAll(resp.Body)
  log.Print(data)
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
