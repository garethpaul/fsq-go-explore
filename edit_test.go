package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEditPageRejectsMissingVenueIDBeforeAuth(t *testing.T) {
	for _, path := range []string{"/edit", "/edit?id=+++"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rr := httptest.NewRecorder()

		EditPage(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("%s status = %d, want %d", path, rr.Code, http.StatusBadRequest)
		}
		if location := rr.Header().Get("Location"); location != "" {
			t.Fatalf("%s redirect location = %q, want none", path, location)
		}
	}
}

func TestProposeEditRejectsNonPostRequests(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/propose_edit?id=venue-1", nil)
	rr := httptest.NewRecorder()

	ProposeEdit(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
	if got := rr.Header().Get("Allow"); got != http.MethodPost {
		t.Fatalf("Allow header = %q, want %q", got, http.MethodPost)
	}
}

func TestProposeEditRejectsMissingVenueID(t *testing.T) {
	for _, path := range []string{"/propose_edit", "/propose_edit?id=+++"} {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		rr := httptest.NewRecorder()

		ProposeEdit(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatalf("%s status = %d, want %d", path, rr.Code, http.StatusBadRequest)
		}
	}
}
