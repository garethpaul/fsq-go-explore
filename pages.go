// Package for application
package app

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/garethpaul/fsq-go-explore/fsq"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

// Setup SearchPageData struct for handling both requests and responses.
type SearchPageData struct {
	SearchPageRequest  fsq.VenueSearchRequest
	SearchPageResponse *fsq.VenueSearchResponse
}

// [START Search_Page]
func SearchPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/search.html")
	if err != nil {
		log.Printf("search template parse failed: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	// Setup Foursquare Client Config
	c := &fsq.FoursquareConfig{
		ClientId:     os.Getenv("FSQ_CLIENT_ID"),
		ClientSecret: os.Getenv("FSQ_CLIENT_SECRET"),
		Client:       getHttpClient(r),
		Version:      os.Getenv("FSQ_VERSION"),
	}

	v := SearchParamParser(r)

	log.Print("search request received")
	// Setup some intelligent caching
	key := fsq.GetSearchKey(v)

	// Execute Memcached Based Caching
	ctx := appengine.NewContext(r)

	// Check for Key
	if item, err := memcache.Get(ctx, key); err == memcache.ErrCacheMiss {

		// Item not in cache
		service := fsq.NewFoursquareService(c)
		resp := service.Search(v)
		item := &memcache.Item{
			Key:        key,
			Object:     resp,
			Expiration: 120 * time.Minute,
		}
		if err := memcache.JSON.Set(ctx, item); err != nil {
			log.Printf("search response cache write failed: %v", err)
		}

		payload := SearchPageData{
			SearchPageRequest:  *v,
			SearchPageResponse: resp,
		}
		if err := t.Execute(w, payload); err != nil {
			log.Printf("search template render failed: %v", err)
		}

	} else if err != nil {

		// Issue getting item from cache
		log.Printf("error getting cached search response: %v", err)
		http.Error(w, "search unavailable", http.StatusInternalServerError)
	} else {

		// Parse from the cache store.
		venues := new(fsq.VenueSearchResponse)
		if err := json.Unmarshal(item.Value, venues); err != nil {
			log.Printf("search response cache decode failed: %v", err)
			http.Error(w, "search unavailable", http.StatusInternalServerError)
			return
		}
		payload := SearchPageData{
			SearchPageRequest:  *v,
			SearchPageResponse: venues,
		}
		if err := t.Execute(w, payload); err != nil {
			log.Printf("search template render failed: %v", err)
		}
	}
}

// [END Search_Page]
