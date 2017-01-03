// Package for application
package app

import (
  "net/http"
  "strings"
  "fsq"
)

// Process a request and cache using headers.
func HeaderCache(fn http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    // Set Key to utilize for caching
    spp := SearchParamParser(r)
    key := fsq.GetSearchKey(spp)

    // Header Based Caching for 2 days
    if match := r.Header.Get("If-None-Match"); match != "" {
      if strings.Contains(match, key) {
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
