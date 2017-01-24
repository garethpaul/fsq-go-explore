// Package for application
package app

import (
	"net/http"
)

func init() {
	http.HandleFunc("/", SearchPage)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/redirect", Redirect)
	http.HandleFunc("/edit", EditPage)
	http.HandleFunc("/propose_edit", ProposeEdit)
}
