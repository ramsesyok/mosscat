package auth

import "testing"

func TestHashAndCompare(t *testing.T) {
	hash, err := Hash("pass")
	if err != nil {
		t.Fatalf("hash error: %v", err)
	}
	if err := Compare(hash, "pass"); err != nil {
		t.Fatalf("compare error: %v", err)
	}
	if err := Compare(hash, "wrong"); err == nil {
		t.Fatalf("expected error on wrong password")
	}
}
