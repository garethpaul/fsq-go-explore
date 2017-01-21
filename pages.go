// Package for application
package app

import (
  "fsq"
  "log"
  "os"
  "appengine"
  "appengine/urlfetch"
  "appengine/memcache"
  "encoding/json"
  "html/template"
  "net/http"
  "time"
)

// Setup SearchPageData struct for handling both requests and responses.
type SearchPageData struct {
	SearchPageRequest    fsq.VenueSearchRequest
  SearchPageResponse   *fsq.VenueSearchResponse
}

// [START Search_Page]
func SearchPage(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("templates/search.html")

  // Setup Foursquare Client Config
  c := &fsq.FoursquareConfig{
  	ClientId: 		os.Getenv("FSQ_CLIENT_ID"),
  	ClientSecret: os.Getenv("FSQ_CLIENT_SECRET"),
  	Client: 			getHttpClient(r),
  	Version: 			os.Getenv("FSQ_VERSION"),
  }

  v := new(fsq.VenueSearchRequest)
  // Take data from our search form and parse this into a struct
  query := r.FormValue("query")
  if query == "" {
    v.Query = "coffee"
  } else {
    v.Query = query
  }
  near := r.FormValue("near")
  if near == "" {
    city := r.Header.Get("X-AppEngine-City")
    region := r.Header.Get("X-AppEngine-Region")
    if city == "" {
      v.Near = "Chicago, IL"
    } else {
      near = city + ", " + region
      v.Near = near
    }
  } else {
    v.Near = near
  }

  log.Print(v)
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
  		Key:   key,
  		Object: resp,
  		Expiration: 120 * time.Minute,
  	}
  	memcache.JSON.Set(ctx, item)

    payload := SearchPageData{
      SearchPageRequest: *v,
      SearchPageResponse: resp,
    }
  	t.Execute(w, payload)
  } else if err != nil {

  	// Issue getting item from cache
  	log.Print(ctx, "error getting item: %v", err)
  } else {

  	// Parse from the cache store.
  	venues := new(fsq.VenueSearchResponse)
  	json.Unmarshal(item.Value, venues)
    payload := SearchPageData{
      SearchPageRequest: *v,
      SearchPageResponse: venues,
    }
    t.Execute(w, payload)
  }
}
// [END Search_Page]

func getHttpClient(r *http.Request) http.Client {
	return http.Client{Transport: &urlfetch.Transport{Context: appengine.NewContext(r)}}
}
