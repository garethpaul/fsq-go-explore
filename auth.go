// Package for application
package app

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"unicode"

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

const (
	userCacheKeyPrefix        = "user:"
	maxOAuthUserResponseBytes = 1 * 1024 * 1024
	oauthUserRequestTimeout   = 10 * time.Second
	oauthUserResponseHost     = "api.foursquare.com"
	oauthUserResponsePath     = "/v2/users/self"
)

var (
	errOAuthUserResponseStatus        = errors.New("foursquare user response status was not successful")
	errOAuthUserResponseOrigin        = errors.New("foursquare user response origin was not expected")
	errOAuthUserResponseMediaType     = errors.New("foursquare user response media type was not JSON")
	errOAuthUserResponseTooLarge      = errors.New("foursquare user response exceeded the size limit")
	errOAuthUserResponseIdentity      = errors.New("foursquare user response identity was invalid")
	errOAuthUserResponseDuplicateKey  = errors.New("foursquare user response contained a duplicate JSON member")
	errOAuthUserResponseJSONStructure = errors.New("foursquare user response JSON structure was invalid")
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
	if !validUserCacheKey(key) {
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

func validUserCacheKey(key string) bool {
	if !strings.HasPrefix(key, userCacheKeyPrefix) {
		return false
	}
	digest := strings.TrimPrefix(key, userCacheKeyPrefix)
	if len(digest) != 64 {
		return false
	}
	for _, ch := range digest {
		if (ch < '0' || ch > '9') && (ch < 'a' || ch > 'f') {
			return false
		}
	}
	return true
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
	user, err := decodeOAuthUserResponse(p)
	if err != nil {
		log.Print("foursquare user response rejected")
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

func decodeOAuthUserResponse(response *http.Response) (*fsq.UserResponse, error) {
	if response == nil || response.Body == nil {
		return nil, errors.New("foursquare user response was missing")
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, errOAuthUserResponseStatus
	}
	if !isExpectedOAuthUserResponseURL(response) {
		return nil, errOAuthUserResponseOrigin
	}
	if !isOAuthUserJSONResponse(response) {
		return nil, errOAuthUserResponseMediaType
	}

	body, err := io.ReadAll(io.LimitReader(response.Body, maxOAuthUserResponseBytes+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxOAuthUserResponseBytes {
		return nil, errOAuthUserResponseTooLarge
	}
	if err := rejectDuplicateJSONMembers(body); err != nil {
		return nil, err
	}

	wrapper := new(fsq.Response)
	if err := json.Unmarshal(body, wrapper); err != nil {
		return nil, err
	}
	user := new(fsq.UserResponse)
	if err := json.Unmarshal(wrapper.Response, user); err != nil {
		return nil, err
	}
	if user.User.ID == "" || user.User.ID != strings.TrimSpace(user.User.ID) {
		return nil, errOAuthUserResponseIdentity
	}
	return user, nil
}

func rejectDuplicateJSONMembers(body []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := consumeUniqueJSONValue(decoder, 0); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		if err == nil {
			return errOAuthUserResponseJSONStructure
		}
		return err
	}
	return nil
}

// foldJSONMemberName mirrors encoding/json's case-insensitive field matching.
func foldJSONMemberName(name string) string {
	return strings.Map(func(r rune) rune {
		for {
			folded := unicode.SimpleFold(r)
			if folded <= r {
				return folded
			}
			r = folded
		}
	}, name)
}

func consumeUniqueJSONValue(decoder *json.Decoder, depth int) error {
	if depth > 10000 {
		return errOAuthUserResponseJSONStructure
	}

	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}

	switch delimiter {
	case '{':
		members := make(map[string]struct{})
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return errOAuthUserResponseJSONStructure
			}
			foldedKey := foldJSONMemberName(key)
			if _, exists := members[foldedKey]; exists {
				return errOAuthUserResponseDuplicateKey
			}
			members[foldedKey] = struct{}{}
			if err := consumeUniqueJSONValue(decoder, depth+1); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil {
			return err
		}
		if closing != json.Delim('}') {
			return errOAuthUserResponseJSONStructure
		}
	case '[':
		for decoder.More() {
			if err := consumeUniqueJSONValue(decoder, depth+1); err != nil {
				return err
			}
		}
		closing, err := decoder.Token()
		if err != nil {
			return err
		}
		if closing != json.Delim(']') {
			return errOAuthUserResponseJSONStructure
		}
	default:
		return errOAuthUserResponseJSONStructure
	}

	return nil
}

func isExpectedOAuthUserResponseURL(response *http.Response) bool {
	if response == nil || response.Request == nil || response.Request.URL == nil {
		return false
	}

	responseURL := response.Request.URL
	return responseURL.Scheme == "https" &&
		strings.EqualFold(responseURL.Hostname(), oauthUserResponseHost) &&
		responseURL.User == nil &&
		responseURL.Port() == "" &&
		responseURL.EscapedPath() == oauthUserResponsePath &&
		responseURL.Fragment == ""
}

func isOAuthUserJSONResponse(response *http.Response) bool {
	contentTypes := response.Header.Values("Content-Type")
	if len(contentTypes) != 1 {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(contentTypes[0])
	if err != nil {
		return false
	}
	mediaType = strings.ToLower(mediaType)
	return mediaType == "application/json" ||
		(strings.HasPrefix(mediaType, "application/") && strings.HasSuffix(mediaType, "+json"))
}

// Process a request and cache using headers.
func LoginProtect(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, _ := r.Cookie("fsq")
		if cookie == nil || !validUserCacheKey(cookie.Value) {
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
	return http.Client{
		Transport:     &urlfetch.Transport{Context: appengine.NewContext(r)},
		CheckRedirect: fsq.RefuseRedirect,
		Timeout:       oauthUserRequestTimeout,
	}
}
