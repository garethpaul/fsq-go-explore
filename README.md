Foursquare-Go-Explore
===================
An example of how to access Foursquare Data using GoLang and AppEngine. You can find this running in the wild via [fsq-go-explore.appspot.com](http://fsq-go-explore.appspot.com/).

![Screenshot](https://raw.githubusercontent.com/garethpaul/fsq-go-explore/master/static/images/screen.png)

Getting Started
----------
Here are some steps to get you started.

#### 1. Get your developer keys

Sign up for a <a href="https://developer.foursquare.com">Foursquare Developer</a> or <a href="https://enterprise.foursquare.com/contact-us">Enterprise Account</a>.

Once you get setup you should  have two important pieces of information:

1. CLIENT_ID - Unique to your registered application.
2. CLIENT_SECRET - Unique and private to your application.

Please remember to keep these keys in a secure location.


#### 2. Clone this repo

``` git clone https://github.com/garethpaul/fsq-go-explore.git```

#### 3. Amend App.yaml with your variables

These are directly environment variables.

```
env_variables:
  FSQ_CLIENT_ID: 'YOUR_FOURSQUARE_KEY' // found in step 1
  FSQ_CLIENT_SECRET: 'YOUR_CLIENT_SECRET' // found in step 1
  FSQ_VERSION: 'YYYYMMDD' // e.g. 20170101
```

##### 4. Run your application

```
goapp serve
```

> **Note:**

> - You'll need to install [GoLang](https://golang.org/doc/install)
> - You will need to download the [Google AppEngine SDK for GoLang](https://cloud.google.com/appengine/docs/go/download)


Caching Guidance
----------
We setup caching in via a header based approach and using memcache to provide a faster experience for users and in accordance with the [Foursquare Developer Platform Policy](https://foursquare.com/legal/api/platformpolicy).

Our caching strategy primarily focuses around a unique key for a given request. This is one primary method for any great caching strategy.  In this demo we opted for the following flow.

1. Convert struct into a json object.
2. Convert the json object into a string.
3. Encrypt the string with AES.
4. Base64Encode the string


#### Header Based Caching
The code for our use of header based caching is as follows.

```
    // Header Based Caching
    if match := r.Header.Get("If-None-Match"); match != "" {
      if strings.Contains(match, key) {
        w.WriteHeader(http.StatusNotModified)
        return
      }
    }

    // Execute Header Based Caching
    w.Header().Set("Etag", key)
    w.Header().Set("Cache-Control", "max-age=23200")
```

We look for an `etag` or `entity tag` that is a mechanism that HTTP provides to deal with caching. If you ask multiple times for the same resource you get the resource for free once it is cached. We therefore send the user a `304` rather send the request through an un-necessary loop.

Additionally we set `Cache-Control` to a specific date that tells the client that once the date expires the cache should revalidate the resource.

> **Recap:**

> - We use a key based system and utilze the `Etag` for cache-control.
> - We additionally use and set `Cache-Control` for setting the max age for the given resource.


More technical specifics can be found on this can be found within [cache.go](https://github.com/garethpaul/fsq-go-explore/blob/master/cache.go).

#### Memory Based Caching

Another primary method for caching is using Memory, in this case we opt for utilizing Memcache.

An example of some `sudo code` for this looks as follows:

```
  // Execute Memcached Based Caching
  ctx := appengine.NewContext(r)

  // Check for Key
  if item, _ := memcache.Get(ctx, key); err == memcache.ErrCacheMiss {

  	// Item is not in cache
  	service := fsq.NewFoursquareService(c)
  	resp := service.Search(v)
  	item := &memcache.Item{
  		Key:   key,
  		Object: resp,
  		Expiration: 120 * time.Minute,
  	}

  	// Set the object in memory.
	return resp

  }  else {

  	// Parse the request from the cache store as it has been found.
  	venues := new(fsq.VenueSearchResponse)
  	json.Unmarshal(item.Value, venues)
	return venues
  }
```

Here we use a very basic caching strategy that utilizes a `key` and looks up that key. Specifically when a key is found we return an `value` from memory.

> **Recap:**

> - We use a very simple `key` and `value` style pairing to determine what we do here.
> - When a `key` is not found, we send out an API request to Foursquare and then set it in cache.
> - Objects can be stored in their original form using AppEngine.
> - We set **two hours** as the maximum time the object is stored in cache.

More technical specifics can be found on this within [pages.go](https://github.com/garethpaul/fsq-go-explore/blob/master/pags.go). 
