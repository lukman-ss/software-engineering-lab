package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"labs/31-oauth2-and-oidc/pkg/oidc"
	"labs/31-oauth2-and-oidc/pkg/pkce"
)

var (
	ErrInvalidRequest       = errors.New("invalid request")
	ErrUnauthorizedClient   = errors.New("unauthorized client or redirect uri")
	ErrInvalidGrant         = errors.New("invalid or expired grant")
	ErrCodeAlreadyUsed      = errors.New("authorization code already redeemed")
	ErrInvalidPKCE          = errors.New("pkce verification failed")
	ErrTokenReplayDetected  = errors.New("refresh token reuse detected: entire token family revoked")
	ErrRefreshTokenRevoked  = errors.New("refresh token is revoked")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

type AuthCode struct {
	Code                string
	ClientID            string
	RedirectURI         string
	Scope               string
	Subject             string
	CodeChallenge       string
	CodeChallengeMethod string
	Nonce               string
	ExpiresAt           time.Time
	Used                bool
}

type RefreshTokenMeta struct {
	FamilyID  string
	Subject   string
	ClientID  string
	Scope     string
	Revoked   bool
	ExpiresAt time.Time
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	Scope        string `json:"scope"`
}

type AuthorizationServer struct {
	mu           sync.Mutex
	Issuer       string
	SigningKey   []byte
	Clients      map[string]string // clientID -> redirectURI
	authCodes    map[string]*AuthCode
	tokens       map[string]string            // accessToken -> subject
	tokenScopes  map[string]string            // accessToken -> scope
	refreshMeta  map[string]*RefreshTokenMeta // tokenString -> meta
	revokedFams  map[string]bool              // familyID -> true
}

func NewAuthorizationServer(issuer string, signingKey []byte) *AuthorizationServer {
	return &AuthorizationServer{
		Issuer:      issuer,
		SigningKey:  signingKey,
		Clients:     make(map[string]string),
		authCodes:   make(map[string]*AuthCode),
		tokens:      make(map[string]string),
		tokenScopes: make(map[string]string),
		refreshMeta: make(map[string]*RefreshTokenMeta),
		revokedFams: make(map[string]bool),
	}
}

func (s *AuthorizationServer) RegisterClient(clientID, redirectURI string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Clients[clientID] = redirectURI
}

func (s *AuthorizationServer) Authorize(clientID, redirectURI, scope, subject, challenge, challengeMethod, nonce string) (*AuthCode, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	expectedRedirect, exists := s.Clients[clientID]
	if !exists || expectedRedirect != redirectURI {
		return nil, ErrUnauthorizedClient
	}

	if challenge == "" {
		return nil, fmt.Errorf("%w: code_challenge is mandatory", ErrInvalidRequest)
	}

	if challengeMethod != "S256" && challengeMethod != "plain" {
		return nil, fmt.Errorf("%w: invalid code_challenge_method", ErrInvalidRequest)
	}

	codeBytes := make([]byte, 16)
	if _, err := rand.Read(codeBytes); err != nil {
		return nil, err
	}
	codeStr := hex.EncodeToString(codeBytes)

	authCode := &AuthCode{
		Code:                codeStr,
		ClientID:            clientID,
		RedirectURI:         redirectURI,
		Scope:               scope,
		Subject:             subject,
		CodeChallenge:       challenge,
		CodeChallengeMethod: challengeMethod,
		Nonce:               nonce,
		ExpiresAt:           time.Now().Add(5 * time.Minute),
		Used:                false,
	}

	s.authCodes[codeStr] = authCode
	return authCode, nil
}

func (s *AuthorizationServer) ExchangeCode(code, clientID, redirectURI, codeVerifier string) (*TokenResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ac, exists := s.authCodes[code]
	if !exists || time.Now().After(ac.ExpiresAt) {
		return nil, ErrInvalidGrant
	}

	if ac.Used {
		return nil, ErrCodeAlreadyUsed
	}

	if ac.ClientID != clientID || ac.RedirectURI != redirectURI {
		return nil, ErrUnauthorizedClient
	}

	if err := pkce.Verify(codeVerifier, ac.CodeChallenge, ac.CodeChallengeMethod); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPKCE, err)
	}

	ac.Used = true

	// Issue Access Token
	atBytes := make([]byte, 24)
	rand.Read(atBytes)
	accessToken := "at_" + hex.EncodeToString(atBytes)
	s.tokens[accessToken] = ac.Subject
	s.tokenScopes[accessToken] = ac.Scope

	// Issue Refresh Token with lineage / family tracking
	familyBytes := make([]byte, 8)
	rand.Read(familyBytes)
	familyID := "fam_" + hex.EncodeToString(familyBytes)

	rtBytes := make([]byte, 24)
	rand.Read(rtBytes)
	refreshToken := "rt_" + hex.EncodeToString(rtBytes)

	s.refreshMeta[refreshToken] = &RefreshTokenMeta{
		FamilyID:  familyID,
		Subject:   ac.Subject,
		ClientID:  ac.ClientID,
		Scope:     ac.Scope,
		Revoked:   false,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	resp := &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: refreshToken,
		Scope:        ac.Scope,
	}

	// Issue ID Token if OpenID scope requested
	if containsOpenID(ac.Scope) {
		now := time.Now()
		claims := oidc.IDTokenClaims{
			Issuer:     s.Issuer,
			Subject:    ac.Subject,
			Audience:   ac.ClientID,
			Expiration: now.Add(1 * time.Hour).Unix(),
			IssuedAt:   now.Unix(),
			AuthTime:   now.Unix(),
			Nonce:      ac.Nonce,
		}
		signedIDToken, err := oidc.SignIDToken(claims, s.SigningKey)
		if err != nil {
			return nil, fmt.Errorf("failed to sign id token: %w", err)
		}
		resp.IDToken = signedIDToken
	}

	return resp, nil
}

func (s *AuthorizationServer) Refresh(refreshToken, clientID string) (*TokenResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	meta, exists := s.refreshMeta[refreshToken]
	if !exists {
		return nil, ErrRefreshTokenNotFound
	}

	// Check if family is revoked (due to prior replay detection)
	if s.revokedFams[meta.FamilyID] {
		return nil, ErrTokenReplayDetected
	}

	// Refresh Token Rotation: If already revoked/used, someone is replaying a stolen token!
	if meta.Revoked {
		// Invalidate the entire token family
		s.revokedFams[meta.FamilyID] = true
		return nil, ErrTokenReplayDetected
	}

	if meta.ClientID != clientID {
		return nil, ErrUnauthorizedClient
	}

	if time.Now().After(meta.ExpiresAt) {
		return nil, ErrInvalidGrant
	}

	// Invalidate the consumed refresh token
	meta.Revoked = true

	// Issue new rotated refresh token in the same family
	newRtBytes := make([]byte, 24)
	rand.Read(newRtBytes)
	newRefreshToken := "rt_" + hex.EncodeToString(newRtBytes)

	s.refreshMeta[newRefreshToken] = &RefreshTokenMeta{
		FamilyID:  meta.FamilyID,
		Subject:   meta.Subject,
		ClientID:  meta.ClientID,
		Scope:     meta.Scope,
		Revoked:   false,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	// Issue new access token
	atBytes := make([]byte, 24)
	rand.Read(atBytes)
	newAccessToken := "at_" + hex.EncodeToString(atBytes)
	s.tokens[newAccessToken] = meta.Subject
	s.tokenScopes[newAccessToken] = meta.Scope

	return &TokenResponse{
		AccessToken:  newAccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: newRefreshToken,
		Scope:        meta.Scope,
	}, nil
}

func (s *AuthorizationServer) ValidateAccessToken(accessToken, requiredScope string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sub, exists := s.tokens[accessToken]
	if !exists {
		return "", errors.New("invalid or expired access token")
	}

	return sub, nil
}

func containsOpenID(scope string) bool {
	for _, part := range splitSpaces(scope) {
		if part == "openid" {
			return true
		}
	}
	return false
}

func splitSpaces(s string) []string {
	var res []string
	start := -1
	for i, r := range s {
		if r == ' ' || r == '\t' {
			if start >= 0 {
				res = append(res, s[start:i])
				start = -1
			}
		} else if start < 0 {
			start = i
		}
	}
	if start >= 0 {
		res = append(res, s[start:])
	}
	return res
}
