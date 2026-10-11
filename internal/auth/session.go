// Package auth issues and verifies signed session cookies using HMAC-SHA256.
package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

// Manager signs and verifies session tokens.
type Manager struct {
	secret []byte
	ttl    time.Duration
}

// NewManager builds a session manager. The secret must be long and random.
func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{secret: []byte(secret), ttl: ttl}
}

// TTL returns the session lifetime.
func (m *Manager) TTL() time.Duration { return m.ttl }

// payload is what gets signed. Keep it small.
type payload struct {
	Username string `json:"u"`
	Expires  int64  `json:"exp"`
	Nonce    string `json:"n"`
}

// Issue returns a signed token of the form "<payload>.<signature>".
func (m *Manager) Issue(username string) (string, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	body, err := json.Marshal(payload{
		Username: username,
		Expires:  time.Now().Add(m.ttl).Unix(),
		Nonce:    base64.RawURLEncoding.EncodeToString(nonce),
	})
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(body)
	return encoded + "." + m.sign(encoded), nil
}

// Verify checks the signature and expiry, returning the username on success.
func (m *Manager) Verify(token string) (string, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return "", false
	}
	if !hmac.Equal([]byte(parts[1]), []byte(m.sign(parts[0]))) {
		return "", false
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", false
	}
	var p payload
	if err := json.Unmarshal(raw, &p); err != nil {
		return "", false
	}
	if p.Username == "" || p.Expires < time.Now().Unix() {
		return "", false
	}
	return p.Username, true
}

func (m *Manager) sign(data string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
