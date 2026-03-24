package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"time"
)

const (
	defaultSecret = "dev-change-this-session-secret"
)

type sessionPayload struct {
	Username string `json:"u"`
	Expires  int64  `json:"exp"`
	Nonce    string `json:"n"`
}

func sessionSecret() []byte {
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = defaultSecret
	}
	return []byte(secret)
}

func GenerateSessionToken(username string, ttl time.Duration) (string, error) {
	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}

	payload := sessionPayload{
		Username: username,
		Expires:  time.Now().Add(ttl).Unix(),
		Nonce:    base64.RawURLEncoding.EncodeToString(nonceBytes),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadJSON)
	signature := sign([]byte(encodedPayload))
	return encodedPayload + "." + signature, nil
}

func VerifySessionToken(token string) bool {
	if token == "" {
		return false
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return false
	}

	payloadPart := parts[0]
	signaturePart := parts[1]

	expectedSig := sign([]byte(payloadPart))
	if !hmac.Equal([]byte(signaturePart), []byte(expectedSig)) {
		return false
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return false
	}

	var payload sessionPayload
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return false
	}

	return payload.Expires > time.Now().Unix()
}

func sign(data []byte) string {
	mac := hmac.New(sha256.New, sessionSecret())
	mac.Write(data)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
