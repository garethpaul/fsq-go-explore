// Foursquare GoLang SDK
package fsq

import (
	"encoding/json"
)

// Meta contains request information and error details
// https://developer.foursquare.com/overview/responses
type Response struct {
	Notifications []Notification  `json:"notifications"`
	Response      json.RawMessage `json:"response"`
	Meta          Meta            `json:"meta"`
}

// Meta contains request information and error details
// https://developer.foursquare.com/overview/responses
type Meta struct {
	Code        int    `json:"code"`
	ErrorType   string `json:"errorType"`
	ErrorDetail string `json:"errorDetail"`
	RequestID   string `json:"requestId"`
}

// Notification comes with all responses.
// https://developer.foursquare.com/docs/responses/notifications
type Notification struct {
	Type string  `json:"type"`
	Item Omitted `json:"item"`
}

// Group contains the default fields in a group. A lot of responses
// share these fields.
type Group struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// A Foursquare User Grouping
type FoursquareUser struct {
  ID               string       `json:"id"`
	Name             string       `json:"name"`
  AccessToken      string       `json:"access_token"`
}


// ID is a simple struct with just an id. VenuePage is an example.
type ID struct {
	ID string `json:"id"`
}

// Count is a simple struct with just a count. Followers and todo are examples
// that only have a count.
type Count struct {
	Count int `json:"count"`
}

// Omitted is for fields that do not have a known datastructure. If you find
// an example where this field is used please let me know. You will need to handle
// this in your application until that time.
type Omitted interface{}

// BoolAsAnInt is a bool that needs to be an int when transferred to an endpoint
type BoolAsAnInt int

// Option available for BoolAsAnInt
const (
	True = BoolAsAnInt(1)
)

// RateLimit is a struct of foursquare ratelimit data
type RateLimit struct {
	Limit     int
	Path      string
	Remaining int
}
