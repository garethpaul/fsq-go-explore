// Package for application
package app

import (
	"net/http"
	"limiter"
	"time"
)

func init() {
	http.Handle("/", limiter.LimitFuncHandler(limiter.NewLimiter(10, time.Minute), HeaderCache(SearchPage)))
	http.HandleFunc("/login", Login)
	http.HandleFunc("/redirect", Redirect)
	http.HandleFunc("/edit", LoginProtect(EditPage))
	http.HandleFunc("/propose_edit", ProposeEdit)
	http.HandleFunc("/logout", Logout)
}
