package fsq

import (
	"strings"
	"testing"
)

func TestGetSearchKeyIsDeterministicAndOpaque(t *testing.T) {
	req := &VenueSearchRequest{Near: "Chicago, IL", Query: "coffee"}

	first := GetSearchKey(req)
	second := GetSearchKey(req)

	if first == "" {
		t.Fatal("GetSearchKey returned an empty key")
	}
	if first != second {
		t.Fatalf("GetSearchKey was not deterministic: %q != %q", first, second)
	}
	if strings.Contains(first, req.Query) || strings.Contains(first, req.Near) {
		t.Fatalf("GetSearchKey exposed request fields in %q", first)
	}
}

func TestGetUserKeyDoesNotExposeAccessToken(t *testing.T) {
	user := &FoursquareUser{ID: "user-1", Name: "Example", AccessToken: "secret-token"}

	key := GetUserKey(user)

	if key == "" {
		t.Fatal("GetUserKey returned an empty key")
	}
	if strings.Contains(key, user.AccessToken) {
		t.Fatalf("GetUserKey exposed access token in %q", key)
	}
}
