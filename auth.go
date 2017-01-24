// Package for application
package app

import (
	"fsq"
	"net/http"
  "os"
  "log"
  "fmt"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/foursquare"
  "golang.org/x/net/context"
  "google.golang.org/appengine"
  "google.golang.org/appengine/urlfetch"
  "google.golang.org/appengine/memcache"
  "encoding/json"
	"io/ioutil"
  "time"
)

var (
  foursquareOauthConfig = &oauth2.Config{
    RedirectURL:	"http://localhost:8080/redirect",
    ClientID:     os.Getenv("FSQ_CLIENT_ID"),
    ClientSecret: os.Getenv("FSQ_CLIENT_SECRET"),
    Endpoint:     foursquare.Endpoint,
  }
  // Setup Foursquare Client Config
  config = &fsq.FoursquareConfig{
    ClientId: 		os.Getenv("FSQ_CLIENT_ID"),
    ClientSecret: os.Getenv("FSQ_CLIENT_SECRET"),
    Version: 			os.Getenv("FSQ_VERSION"),
    AuthConfig:  foursquareOauthConfig,
  }
  STATE_STR = "12312DWD213DW23D2SD"
)


func getAccessToken(r *http.Request, key string) (accessToken string) {
  var ctx context.Context = appengine.NewContext(r)
  //
  // Check for Key
  if item, err := memcache.Get(ctx, key); err == memcache.ErrCacheMiss {

    // Getting data
    log.Print(ctx, "None found")
    return ""

  } else if err != nil {

    // Issue getting item from cache
    log.Print(ctx, "error getting item: %v", err)
  } else {

    // Parse from the cache store.
    user := new(fsq.FoursquareUser)
    json.Unmarshal(item.Value, user)
    return user.AccessToken
  }
  return ""
}

func setAccessToken(r *http.Request, fsqUser *fsq.FoursquareUser) {
  key := fsq.GetUserKey(fsqUser)
  var ctx context.Context = appengine.NewContext(r)
  item := &memcache.Item{
    Key:   key,
    Object: fsqUser,
  }
  memcache.JSON.Set(ctx, item)
}

// [START Search_Page]
func Login(w http.ResponseWriter, r *http.Request) {
  url := config.AuthConfig.AuthCodeURL(STATE_STR)
  http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//
func Redirect(w http.ResponseWriter, r *http.Request) {
  log.Print("receivedFoursquareCallback")
	state := r.FormValue("state")
	if state != STATE_STR {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", STATE_STR, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

  var ctx context.Context = appengine.NewContext(r)
	code := r.FormValue("code")
	token, err := foursquareOauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Print("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

  log.Print("oAuth Confirmed")
  user_url := "https://api.foursquare.com/v2/users/self?v=20170101&oauth_token=" + token.AccessToken
  c := getHttpClient(r)
  p, err := c.Get(user_url)

  if err != nil {
    log.Println(err)
    http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
    return
  }
  defer p.Body.Close()
  d, _ := ioutil.ReadAll(p.Body)

  // Map the resp
  user := new(fsq.UserResponse)
  response := new(fsq.Response)

  // More details via https://developer.foursquare.com/overview/responses
  json.Unmarshal([]byte((string(d))), &response)
  json.Unmarshal(response.Response, user)

  log.Print(user)

  //
  // Setup Foursquare Client Config
  newUser := &fsq.FoursquareUser{
    ID: user.User.ID,
    Name: user.User.FirstName,
    AccessToken: token.AccessToken,
  }

  setAccessToken(r, newUser)
  accessToken := getAccessToken(r, fsq.GetUserKey(newUser))

  expiration := time.Now().Add(365 * 24 * time.Hour)
  cookie := http.Cookie{Name: "fsq",
                        Value: fsq.GetUserKey(newUser),
                        Expires:expiration}
  http.SetCookie(w, &cookie)
  log.Print(accessToken)
}


func getHttpClient(r *http.Request) http.Client {
	return http.Client{Transport: &urlfetch.Transport{Context: appengine.NewContext(r)}}
}
