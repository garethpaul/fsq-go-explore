package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garethpaul/fsq-go-explore/fsq"
)

func TestHeaderCacheRejectsPartialETagMatches(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?query=coffee&near=Seattle", nil)
	key := fsq.GetSearchKey(SearchParamParser(req))
	req.Header.Set("If-None-Match", "prefix"+key+"suffix")

	called := false
	handler := HeaderCache(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code == http.StatusNotModified {
		t.Fatal("partial ETag match returned 304")
	}
	if !called {
		t.Fatal("handler was not called for partial ETag match")
	}
}

func TestHeaderCacheAcceptsExactETagFromList(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?query=coffee&near=Seattle", nil)
	key := fsq.GetSearchKey(SearchParamParser(req))
	req.Header.Set("If-None-Match", `miss, "`+key+`"`)

	called := false
	handler := HeaderCache(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	rr := httptest.NewRecorder()

	handler(rr, req)

	if rr.Code != http.StatusNotModified {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNotModified)
	}
	if called {
		t.Fatal("handler was called for exact ETag match")
	}
}
