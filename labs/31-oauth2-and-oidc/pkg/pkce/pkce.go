package pkce

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

var (
	ErrInvalidVerifierLength = errors.New("code_verifier length must be between 43 and 128 characters")
	ErrInvalidMethod         = errors.New("unsupported code_challenge_method")
	ErrChallengeMismatch     = errors.New("code_verifier does not match code_challenge")
)

type PKCEPair struct {
	CodeVerifier  string
	CodeChallenge string
	Method        string
}

func GeneratePKCEPair(method string) (*PKCEPair, error) {
	if method != "S256" && method != "plain" {
		return nil, ErrInvalidMethod
	}

	// 32 random bytes -> 43 base64url chars
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to generate random verifier: %w", err)
	}
	verifier := base64.RawURLEncoding.EncodeToString(buf)

	challenge, err := ComputeChallenge(verifier, method)
	if err != nil {
		return nil, err
	}

	return &PKCEPair{
		CodeVerifier:  verifier,
		CodeChallenge: challenge,
		Method:        method,
	}, nil
}

func ComputeChallenge(verifier string, method string) (string, error) {
	if len(verifier) < 43 || len(verifier) > 128 {
		return "", ErrInvalidVerifierLength
	}

	switch method {
	case "S256":
		h := sha256.Sum256([]byte(verifier))
		return base64.RawURLEncoding.EncodeToString(h[:]), nil
	case "plain":
		return verifier, nil
	default:
		return "", ErrInvalidMethod
	}
}

func Verify(verifier string, challenge string, method string) error {
	expected, err := ComputeChallenge(verifier, method)
	if err != nil {
		return err
	}
	if expected != challenge {
		return ErrChallengeMismatch
	}
	return nil
}
