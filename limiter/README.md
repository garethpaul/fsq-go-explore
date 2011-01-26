##  Limiter

This is a generic middleware to rate-limit HTTP requests

## Tutorial
```
package main

import (
    "github.com/garethpaul/fsq-go-explore/limiter"
    "net/http"
    "time"
)

func HelloHandler(w http.ResponseWriter, req *http.Request) {
    w.Write([]byte("Hello, World!"))
}

func main() {
    // Create a request limiter per handler.
    http.Handle("/", limiter.LimitFuncHandler(limiter.NewLimiter(1, time.Second), HelloHandler))
    http.ListenAndServe(":12345", nil)
}
```

## Features

1. Rate-limit by request's remote IP, path, methods, custom headers, & basic auth usernames.
    ```
    l := limiter.NewLimiter(1, time.Second)

    // Configure list of places to look for IP address.
    // By default it's: "RemoteAddr", "X-Forwarded-For", "X-Real-IP"
    // If your application is behind a proxy, set "X-Forwarded-For" first.
    l.IPLookups = []string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}

    // Limit only GET and POST requests.
    l.Methods = []string{"GET", "POST"}

    // Limit request headers containing certain values.
    // Typically, you prefetched these values from the database.
    l.Headers = make(map[string][]string)
    l.Headers["X-Access-Token"] = []string{"abc", "123"}

    // Limit based on basic auth usernames.
    // Typically, you prefetched these values from the database.
    l.BasicAuthUsers = []string{"bob", "joe", "didip"}
    ```

2. Each request handler can be rate-limited individually.

3. Compose your own middleware by using `LimitByKeys()`.

4. Limiter does not require external storage since it uses an algorithm called [Token Bucket](http://en.wikipedia.org/wiki/Token_bucket) [(Go library: golang.org/x/time/rate)](//godoc.org/golang.org/x/time/rate).
