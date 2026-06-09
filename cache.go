// Package for application
package app

import (
	"net/http"
	"strings"

	"github.com/garethpaul/fsq-go-explore/fsq"
)

// Process a request and cache using headers.
func HeaderCache(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set Key to utilize for caching
		spp := SearchParamParser(r)
		key := fsq.GetSearchKey(spp)

		// Header Based Caching for 2 days
		if match := r.Header.Get("If-None-Match"); match != "" {
			if etagMatches(match, key) {
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}

		// Execute Header Based Caching
		w.Header().Set("Etag", key)
		w.Header().Set("Cache-Control", "max-age=23200")

		fn(w, r)
	}
}

func etagMatches(headerValue, key string) bool {
	for _, candidate := range strings.Split(headerValue, ",") {
		candidate = strings.TrimSpace(candidate)
		candidate = strings.Trim(candidate, `"`)
		if candidate == key {
			return true
		}
	}
	return false
}
