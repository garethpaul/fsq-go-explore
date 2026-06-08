package app

import "testing"

func TestNewOAuthStateReturnsDistinctOpaqueValues(t *testing.T) {
	first, err := newOAuthState()
	if err != nil {
		t.Fatal(err)
	}
	second, err := newOAuthState()
	if err != nil {
		t.Fatal(err)
	}

	if first == "" || second == "" {
		t.Fatal("expected OAuth states to be populated")
	}
	if first == second {
		t.Fatal("expected distinct OAuth states")
	}
}
