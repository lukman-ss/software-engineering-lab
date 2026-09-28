package oidc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrMalformedJWT      = errors.New("malformed jwt token")
	ErrSignatureInvalid  = errors.New("invalid jwt signature")
	ErrIssuerMismatch    = errors.New("id token issuer mismatch")
	ErrAudienceMismatch  = errors.New("id token audience mismatch")
	ErrTokenExpired      = errors.New("id token has expired")
	ErrNonceMismatch     = errors.New("id token nonce mismatch")
	ErrIssuedInFuture    = errors.New("id token issued in future")
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type IDTokenClaims struct {
	Issuer          string `json:"iss"`
	Subject         string `json:"sub"`
	Audience        string `json:"aud"`
	Expiration      int64  `json:"exp"`
	IssuedAt        int64  `json:"iat"`
	AuthTime        int64  `json:"auth_time,omitempty"`
	Nonce           string `json:"nonce,omitempty"`
	AuthorizedParty string `json:"azp,omitempty"`
}

func SignIDToken(claims IDTokenClaims, secret []byte) (string, error) {
	header := Header{Alg: "HS256", Typ: "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	unsignedToken := fmt.Sprintf("%s.%s", headerB64, claimsB64)

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(unsignedToken))
	sig := mac.Sum(nil)
	sigB64 := base64.RawURLEncoding.EncodeToString(sig)

	return fmt.Sprintf("%s.%s", unsignedToken, sigB64), nil
}

func ParseAndVerifyIDToken(rawJWT string, secret []byte, expectedIssuer, expectedAudience, expectedNonce string, now time.Time) (*IDTokenClaims, error) {
	parts := strings.Split(rawJWT, ".")
	if len(parts) != 3 {
		return nil, ErrMalformedJWT
	}

	unsignedToken := fmt.Sprintf("%s.%s", parts[0], parts[1])
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, ErrMalformedJWT
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(unsignedToken))
	expectedSig := mac.Sum(nil)

	if !hmac.Equal(sig, expectedSig) {
		return nil, ErrSignatureInvalid
	}

	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrMalformedJWT
	}

	var claims IDTokenClaims
	if err := json.Unmarshal(claimsBytes, &claims); err != nil {
		return nil, ErrMalformedJWT
	}

	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("%w: expected %s, got %s", ErrIssuerMismatch, expectedIssuer, claims.Issuer)
	}

	if claims.Audience != expectedAudience {
		return nil, fmt.Errorf("%w: expected %s, got %s", ErrAudienceMismatch, expectedAudience, claims.Audience)
	}

	unixNow := now.Unix()
	if claims.Expiration <= unixNow {
		return nil, ErrTokenExpired
	}

	if claims.IssuedAt > unixNow+300 {
		return nil, ErrIssuedInFuture
	}

	if expectedNonce != "" && claims.Nonce != expectedNonce {
		return nil, fmt.Errorf("%w: expected %s, got %s", ErrNonceMismatch, expectedNonce, claims.Nonce)
	}

	return &claims, nil
}
