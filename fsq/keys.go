// Foursquare GoLang SDK
package fsq

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// GetSearchKey returns a deterministic cache key for a venue search.
func GetSearchKey(q *VenueSearchRequest) string {
	return cacheKey("search", q)
}

// GetUserKey returns an opaque cache key for a Foursquare user.
func GetUserKey(q *FoursquareUser) string {
	return cacheKey("user", q)
}

func cacheKey(namespace string, value interface{}) string {
	out, err := json.Marshal(value)
	if err != nil {
		out = []byte(namespace)
	}
	sum := sha256.Sum256(out)
	return namespace + ":" + hex.EncodeToString(sum[:])
}
