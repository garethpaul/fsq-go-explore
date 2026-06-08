// Package for application
package app

import (
	"net/http"
	"time"

	"github.com/garethpaul/fsq-go-explore/limiter"
)

func init() {
	http.Handle("/", limiter.LimitFuncHandler(limiter.NewLimiter(10, time.Minute), HeaderCache(SearchPage)))
	http.HandleFunc("/login", Login)
	http.HandleFunc("/redirect", Redirect)
	http.HandleFunc("/edit", LoginProtect(EditPage))
	http.HandleFunc("/propose_edit", LoginProtect(ProposeEdit))
	http.HandleFunc("/logout", Logout)
}
