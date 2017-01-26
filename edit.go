// Package for application
package app

import (
  "html/template"
  "net/http"
  "fsq"
  //"log"
  "os"
)

// [START Edit_Page]
func EditPage(w http.ResponseWriter, r *http.Request) {

  cookie, _ := r.Cookie("fsq")
  accessToken := getAccessToken(r, cookie.Value)
  if accessToken == "" {
    http.Redirect(w, r, "/logout", http.StatusTemporaryRedirect)
  }

  t := template.Must(template.ParseFiles("templates/edit.html", "templates/navigation.html"))

  // Setup Foursquare Client Config
  c := &fsq.FoursquareConfig{
    ClientId: 		os.Getenv("FSQ_CLIENT_ID"),
    ClientSecret: os.Getenv("FSQ_CLIENT_SECRET"),
    Client: 			getHttpClient(r),
    Version: 			os.Getenv("FSQ_VERSION"),
    AccessToken: accessToken,
  }


  id := r.FormValue("id")
  service := fsq.NewFoursquareService(c)
  resp := service.VenueDetails(id)
  t.Execute(w, resp)
}
// [END Edit_Page]


func ProposeEdit(w http.ResponseWriter, r *http.Request) {

  cookie, _ := r.Cookie("fsq")
  accessToken := getAccessToken(r, cookie.Value)

  // Setup Foursquare Client Config
  c := &fsq.FoursquareConfig{
    ClientId: 		os.Getenv("FSQ_CLIENT_ID"),
    ClientSecret: os.Getenv("FSQ_CLIENT_SECRET"),
    Client: 			getHttpClient(r),
    Version: 			os.Getenv("FSQ_VERSION"),
    AccessToken: accessToken,
  }

  id := r.FormValue("id")
  r.ParseForm()

  service := fsq.NewFoursquareService(c)
  service.VenueEdit(id, r.PostForm)

  http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
  //log.Print(r.PostForm)
}
