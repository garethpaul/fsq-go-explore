package app

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestSearchParamParserUsesExplicitNear(t *testing.T) {
	req := httptest.NewRequest("GET", "/?query=+coffee+&near=+Seattle%2C+WA+", nil)

	got := SearchParamParser(req)

	if got.Query != "coffee" {
		t.Fatalf("Query = %q, want coffee", got.Query)
	}
	if got.Near != "Seattle, WA" {
		t.Fatalf("Near = %q, want Seattle, WA", got.Near)
	}
}

func TestSearchParamParserUsesAppEngineHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/?query=tea", nil)
	req.Header.Set("X-AppEngine-City", "Portland")
	req.Header.Set("X-AppEngine-Region", "OR")

	got := SearchParamParser(req)

	if got.Near != "Portland, OR" {
		t.Fatalf("Near = %q, want Portland, OR", got.Near)
	}
}

func TestSearchParamParserDefaultsLocation(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	got := SearchParamParser(req)

	if got.Query != "coffee" {
		t.Fatalf("Query = %q, want coffee", got.Query)
	}
	if got.Near != "Chicago, IL" {
		t.Fatalf("Near = %q, want Chicago, IL", got.Near)
	}
}

func TestSearchParamParserLimitsLongInputs(t *testing.T) {
	longValue := strings.Repeat("x", maxSearchParamRunes+25)
	values := url.Values{
		"query": []string{longValue},
		"near":  []string{longValue},
	}
	req := httptest.NewRequest("GET", "/?"+values.Encode(), nil)

	got := SearchParamParser(req)

	if len([]rune(got.Query)) != maxSearchParamRunes {
		t.Fatalf("Query length = %d, want %d", len([]rune(got.Query)), maxSearchParamRunes)
	}
	if len([]rune(got.Near)) != maxSearchParamRunes {
		t.Fatalf("Near length = %d, want %d", len([]rune(got.Near)), maxSearchParamRunes)
	}
}
