// Foursquare GoLang SDK
package fsq

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Setup CONSTANTS
const (
	SEARCH_URL = "https://api.foursquare.com/v2/venues/search?"
)

// Struct for FourceService to wrap around requests.
type FoursquareService struct {
	Config *FoursquareConfig
}

// Struct for configuration file.
type FoursquareConfig struct {
	ClientId     string
	ClientSecret string
	Client       http.Client
	Version      string
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

// Get the url params
func (fsqs *FoursquareConfig) clientParams() url.Values {
	params := url.Values{}
	params.Set("client_id", fsqs.ClientId)
	params.Set("client_secret", fsqs.ClientSecret)
	params.Set("v", fsqs.Version)
	return params
}
