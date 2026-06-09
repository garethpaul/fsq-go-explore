package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
