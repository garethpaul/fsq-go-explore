// Package for application
package app

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/garethpaul/fsq-go-explore/fsq"
)

// [START Edit_Page]
func EditPage(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		http.Error(w, "missing venue id", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("fsq")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	accessToken := getAccessToken(r, cookie.Value)
	if accessToken == "" {
		http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
		return
	}

	t, err := template.ParseFiles("templates/edit.html", "templates/navigation.html")
	if err != nil {
		log.Printf("edit template parse failed: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}

	// Setup Foursquare Client Config
	c := &fsq.FoursquareConfig{
		ClientId:     os.Getenv("FSQ_CLIENT_ID"),
		ClientSecret: os.Getenv("FSQ_CLIENT_SECRET"),
		Client:       getHttpClient(r),
		Version:      os.Getenv("FSQ_VERSION"),
		AccessToken:  accessToken,
	}

	service := fsq.NewFoursquareService(c)
	resp := service.VenueDetails(id)
	if err := t.Execute(w, resp); err != nil {
		log.Printf("edit template render failed: %v", err)
	}
}

// [END Edit_Page]

func ProposeEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Print("venue edit form parse failed")
		http.Error(w, "invalid venue edit form", http.StatusBadRequest)
		return
	}

	id := strings.TrimSpace(r.Form.Get("id"))
	if id == "" {
		http.Error(w, "missing venue id", http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("fsq")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	accessToken := getAccessToken(r, cookie.Value)
	if accessToken == "" {
		http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
		return
	}

	// Setup Foursquare Client Config
	c := &fsq.FoursquareConfig{
		ClientId:     os.Getenv("FSQ_CLIENT_ID"),
		ClientSecret: os.Getenv("FSQ_CLIENT_SECRET"),
		Client:       getHttpClient(r),
		Version:      os.Getenv("FSQ_VERSION"),
		AccessToken:  accessToken,
	}

	service := fsq.NewFoursquareService(c)
	service.VenueEdit(id, r.PostForm)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
