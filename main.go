// Package for application
package app

import (
	"net/http"
)

func init() {
	http.HandleFunc("/", HeaderCache(SearchPage))
	http.HandleFunc("/login", Login)
	http.HandleFunc("/redirect", Redirect)
	http.HandleFunc("/edit", LoginProtect(EditPage))
	http.HandleFunc("/propose_edit", ProposeEdit)
	http.HandleFunc("/logout", Logout)
}
