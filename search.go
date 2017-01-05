// Package for application
package app

import (
	"fsq"
	"net/http"
)

// SearchParamParser takes a http request and returns a venue request struct.
func SearchParamParser(r *http.Request) (vsr *fsq.VenueSearchRequest) {
	// Take data from our search form and parse this into a struct
	query := r.FormValue("query")
	near := r.FormValue("near")

	v := new(fsq.VenueSearchRequest)
	v.Query = query

	if near == "" {
		city := r.Header.Get("X-AppEngine-City")
		region := r.Header.Get("X-AppEngine-Region")
		near = city + ", " + region
		v.Near = near
	} else {
		v.Near = near
	}
	return v
}
