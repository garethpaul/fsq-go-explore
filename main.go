// Package for application
package app

import (
	"net/http"
)

func init() {
	http.HandleFunc("/", HeaderCache(SearchPage))
}
