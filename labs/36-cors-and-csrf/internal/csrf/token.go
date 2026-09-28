package csrf

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("csrf: invalid or tampered token")
	ErrExpiredToken = errors.New("csrf: token has expired")
	ErrTokenMissing = errors.New("csrf: token is missing")
)

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenManager(secret []byte, ttl time.Duration) *TokenManager {
	if len(secret) == 0 {
		secret = make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			panic(fmt.Sprintf("failed to generate random secret: %v", err))
		}
	}
	if ttl <= 0 {
		ttl = 1 * time.Hour
	}
	return &TokenManager{
		secret: secret,
		ttl:    ttl,
	}
}

// GenerateToken creates a signed token: base64(sessionID:timestamp:nonce:hmac)
func (tm *TokenManager) GenerateToken(sessionID string) string {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}
	nonceStr := base64.RawURLEncoding.EncodeToString(nonce)
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	payload := sessionID + ":" + ts + ":" + nonceStr
	h := hmac.New(sha256.New, tm.secret)
	h.Write([]byte(payload))
	signature := h.Sum(nil)

	tokenRaw := payload + ":" + base64.RawURLEncoding.EncodeToString(signature)
	return base64.RawURLEncoding.EncodeToString([]byte(tokenRaw))
}

// ValidateToken verifies HMAC signature, sessionID match, and expiration.
func (tm *TokenManager) ValidateToken(tokenStr, expectedSessionID string) error {
	if tokenStr == "" {
		return ErrTokenMissing
	}

	decoded, err := base64.RawURLEncoding.DecodeString(tokenStr)
	if err != nil {
		return ErrInvalidToken
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 4 {
		return ErrInvalidToken
	}

	sessionID := parts[0]
	tsStr := parts[1]
	nonceStr := parts[2]
	sigB64 := parts[3]

	if sessionID != expectedSessionID {
		return ErrInvalidToken
	}

	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return ErrInvalidToken
	}

	if time.Since(time.Unix(ts, 0)) > tm.ttl {
		return ErrExpiredToken
	}

	payload := sessionID + ":" + tsStr + ":" + nonceStr
	h := hmac.New(sha256.New, tm.secret)
	h.Write([]byte(payload))
	expectedSig := h.Sum(nil)

	actualSig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return ErrInvalidToken
	}

	if subtle.ConstantTimeCompare(actualSig, expectedSig) != 1 {
		return ErrInvalidToken
	}

	return nil
}
