// Package for application
package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/garethpaul/fsq-go-explore/fsq"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/foursquare"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"
)

var (
	foursquareOauthConfig = &oauth2.Config{
		RedirectURL:  "https://fsq-go-explore.appspot.com/redirect",
		ClientID:     os.Getenv("FSQ_CLIENT_ID"),
		ClientSecret: os.Getenv("FSQ_CLIENT_SECRET"),
		Endpoint:     foursquare.Endpoint,
	}
	// Setup Foursquare Client Config
	config = &fsq.FoursquareConfig{
		ClientId:     os.Getenv("FSQ_CLIENT_ID"),
		ClientSecret: os.Getenv("FSQ_CLIENT_SECRET"),
		Version:      os.Getenv("FSQ_VERSION"),
		AuthConfig:   foursquareOauthConfig,
	}
	oauthStateCookieName = "fsq_oauth_state"
)

func newOAuthState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func secureCookie(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-AppEngine-Https") == "on" || r.Header.Get("X-Forwarded-Proto") == "https"
}

func getAccessToken(r *http.Request, key string) string {
	if key == "" {
		return ""
	}
	ctx := appengine.NewContext(r)
	item, err := memcache.Get(ctx, key)
	switch {
	case err == memcache.ErrCacheMiss:
		return ""
	case err != nil:
		log.Print("error getting access token cache item")
		return ""
	}

	user := new(fsq.FoursquareUser)
	if err := json.Unmarshal(item.Value, user); err != nil {
		log.Print("error decoding access token cache item")
		return ""
	}
	return user.AccessToken
}

func setAccessToken(r *http.Request, fsqUser *fsq.FoursquareUser) {
	if fsqUser == nil {
		return
	}
	key := fsq.GetUserKey(fsqUser)
	ctx := appengine.NewContext(r)
	item := &memcache.Item{
		Key:    key,
		Object: fsqUser,
	}
	if err := memcache.JSON.Set(ctx, item); err != nil {
		log.Print("error setting access token cache item")
	}
}

// [START Search_Page]
func Login(w http.ResponseWriter, r *http.Request) {
	state, err := newOAuthState()
	if err != nil {
		log.Print("failed to create oauth state")
		http.Error(w, "login unavailable", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   secureCookie(r),
		SameSite: http.SameSiteLaxMode,
	})

	url := config.AuthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	log.Print("received Foursquare callback")
	state := r.FormValue("state")
	stateCookie, err := r.Cookie(oauthStateCookieName)
	if err != nil || stateCookie.Value == "" || state != stateCookie.Value {
		log.Print("invalid oauth state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := strings.TrimSpace(r.FormValue("code"))
	if code == "" {
		log.Print("missing oauth code")
		http.SetCookie(w, &http.Cookie{
			Name:     oauthStateCookieName,
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   secureCookie(r),
			SameSite: http.SameSiteLaxMode,
		})
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     oauthStateCookieName,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secureCookie(r),
		SameSite: http.SameSiteLaxMode,
	})

	ctx := appengine.NewContext(r)
	token, err := foursquareOauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Print("oauth exchange failed")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	params := url.Values{}
	params.Set("v", "20170101")
	params.Set("oauth_token", token.AccessToken)
	userURL := "https://api.foursquare.com/v2/users/self?" + params.Encode()
	c := getHttpClient(r)
	p, err := c.Get(userURL)
	if err != nil {
		log.Print("foursquare user request failed")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer p.Body.Close()
	d, err := io.ReadAll(p.Body)
	if err != nil {
		log.Print("foursquare user response read failed")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	user := new(fsq.UserResponse)
	response := new(fsq.Response)
	if err := json.Unmarshal(d, response); err != nil {
		log.Print("foursquare user wrapper decode failed")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	if err := json.Unmarshal(response.Response, user); err != nil {
		log.Print("foursquare user decode failed")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	newUser := &fsq.FoursquareUser{
		ID:          user.User.ID,
		Name:        user.User.FirstName,
		AccessToken: token.AccessToken,
	}
	userKey := fsq.GetUserKey(newUser)
	setAccessToken(r, newUser)

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "fsq",
		Value:    userKey,
		Path:     "/",
		Expires:  expiration,
		HttpOnly: true,
		Secure:   secureCookie(r),
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// Process a request and cache using headers.
func LoginProtect(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("fsq")
		if cookie == nil {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		fn(w, r)
	}
}

// Process a request and cache using headers.
func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "fsq",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secureCookie(r),
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func getHttpClient(r *http.Request) http.Client {
	return http.Client{Transport: &urlfetch.Transport{Context: appengine.NewContext(r)}}
}
