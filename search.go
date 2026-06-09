// Package for application
package app

import (
	"net/http"
	"strings"

	"github.com/garethpaul/fsq-go-explore/fsq"
)

const maxSearchParamRunes = 120

// SearchParamParser takes a http request and returns a venue request struct.
func SearchParamParser(r *http.Request) (vsr *fsq.VenueSearchRequest) {
	// Take data from our search form and parse this into a struct
	query := normalizeSearchParam(r.FormValue("query"))
	near := normalizeSearchParam(r.FormValue("near"))
	if query == "" {
		query = "coffee"
	}

	v := new(fsq.VenueSearchRequest)
	v.Query = query

	if near == "" {
		v.Near = appEngineLocationFallback(r)
	} else {
		v.Near = near
	}
	return v
}

func appEngineLocationFallback(r *http.Request) string {
	city := normalizeSearchParam(r.Header.Get("X-AppEngine-City"))
	region := normalizeSearchParam(r.Header.Get("X-AppEngine-Region"))

	switch {
	case city != "" && region != "":
		return city + ", " + region
	case city != "":
		return city
	case region != "":
		return region
	default:
		return "Chicago, IL"
	}
}

func normalizeSearchParam(value string) string {
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) <= maxSearchParamRunes {
		return value
	}
	return string(runes[:maxSearchParamRunes])
}
