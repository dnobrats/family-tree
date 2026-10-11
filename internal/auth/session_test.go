package auth

import (
	"testing"
	"time"
)

func TestIssueVerify_RoundTrip(t *testing.T) {
	m := NewManager("a-very-secret-value-1234", time.Hour)
	token, err := m.Issue("admin")
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	user, ok := m.Verify(token)
	if !ok || user != "admin" {
		t.Fatalf("verify failed: user=%q ok=%v", user, ok)
	}
}

func TestVerify_RejectsTampered(t *testing.T) {
	m := NewManager("a-very-secret-value-1234", time.Hour)
	token, _ := m.Issue("admin")

	if _, ok := m.Verify(token + "x"); ok {
		t.Fatal("tampered signature must be rejected")
	}
	if _, ok := m.Verify("not-a-token"); ok {
		t.Fatal("malformed token must be rejected")
	}
	if _, ok := m.Verify(""); ok {
		t.Fatal("empty token must be rejected")
	}
}

func TestVerify_RejectsExpired(t *testing.T) {
	m := NewManager("a-very-secret-value-1234", -time.Minute)
	token, _ := m.Issue("admin")
	if _, ok := m.Verify(token); ok {
		t.Fatal("expired token must be rejected")
	}
}

func TestVerify_RejectsWrongSecret(t *testing.T) {
	issuer := NewManager("secret-one-1234567890", time.Hour)
	verifier := NewManager("secret-two-1234567890", time.Hour)
	token, _ := issuer.Issue("admin")
	if _, ok := verifier.Verify(token); ok {
		t.Fatal("token signed with another secret must be rejected")
	}
}
