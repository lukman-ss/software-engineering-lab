package tests

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"sync"
	"testing"
	"time"

	"labs/31-oauth2-and-oidc/pkg/client"
	"labs/31-oauth2-and-oidc/pkg/oidc"
	"labs/31-oauth2-and-oidc/pkg/pkce"
	"labs/31-oauth2-and-oidc/pkg/server"
)

func TestPKCE_S256_Valid(t *testing.T) {
	pair, err := pkce.GeneratePKCEPair("S256")
	if err != nil {
		t.Fatalf("GeneratePKCEPair failed: %v", err)
	}

	if len(pair.CodeVerifier) < 43 || len(pair.CodeVerifier) > 128 {
		t.Errorf("invalid verifier length: %d", len(pair.CodeVerifier))
	}

	if err := pkce.Verify(pair.CodeVerifier, pair.CodeChallenge, "S256"); err != nil {
		t.Errorf("Verify failed for valid pair: %v", err)
	}
}

func TestPKCE_InvalidMethod(t *testing.T) {
	_, err := pkce.GeneratePKCEPair("INVALID")
	if err == nil {
		t.Errorf("expected error on invalid PKCE method")
	}
}

func TestPKCE_Mismatch(t *testing.T) {
	pair1, _ := pkce.GeneratePKCEPair("S256")
	pair2, _ := pkce.GeneratePKCEPair("S256")

	if err := pkce.Verify(pair1.CodeVerifier, pair2.CodeChallenge, "S256"); err == nil {
		t.Errorf("expected challenge mismatch error")
	}
}

func TestOIDC_IDToken_Valid(t *testing.T) {
	secret := []byte("test-signing-key-123456789012")
	issuer := "https://auth.example.com"
	audience := "client-abc"
	nonce := "random-nonce-123"

	now := time.Now()
	claims := oidc.IDTokenClaims{
		Issuer:     issuer,
		Subject:    "user_123",
		Audience:   audience,
		Expiration: now.Add(time.Hour).Unix(),
		IssuedAt:   now.Unix(),
		Nonce:      nonce,
	}

	signed, err := oidc.SignIDToken(claims, secret)
	if err != nil {
		t.Fatalf("SignIDToken failed: %v", err)
	}

	parsed, err := oidc.ParseAndVerifyIDToken(signed, secret, issuer, audience, nonce, now)
	if err != nil {
		t.Fatalf("ParseAndVerifyIDToken failed: %v", err)
	}

	if parsed.Subject != "user_123" {
		t.Errorf("expected subject user_123, got %s", parsed.Subject)
	}
}

func TestOIDC_IDToken_TamperedSignature(t *testing.T) {
	secret := []byte("test-signing-key-123456789012")
	claims := oidc.IDTokenClaims{
		Issuer:     "https://auth.example.com",
		Subject:    "user_123",
		Audience:   "client-abc",
		Expiration: time.Now().Add(time.Hour).Unix(),
		IssuedAt:   time.Now().Unix(),
	}

	signed, _ := oidc.SignIDToken(claims, secret)
	parts := strings.Split(signed, ".")
	tampered := parts[0] + "." + parts[1] + ".invalidSig=="

	_, err := oidc.ParseAndVerifyIDToken(tampered, secret, claims.Issuer, claims.Audience, "", time.Now())
	if err == nil {
		t.Errorf("expected signature verification failure")
	}
}

func TestOIDC_IDToken_Expired(t *testing.T) {
	secret := []byte("test-signing-key-123456789012")
	claims := oidc.IDTokenClaims{
		Issuer:     "https://auth.example.com",
		Subject:    "user_123",
		Audience:   "client-abc",
		Expiration: time.Now().Add(-10 * time.Minute).Unix(),
		IssuedAt:   time.Now().Add(-20 * time.Minute).Unix(),
	}

	signed, _ := oidc.SignIDToken(claims, secret)
	_, err := oidc.ParseAndVerifyIDToken(signed, secret, claims.Issuer, claims.Audience, "", time.Now())
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Errorf("expected token expired error, got: %v", err)
	}
}

func TestOIDC_IDToken_MismatchClaims(t *testing.T) {
	secret := []byte("test-signing-key-123456789012")
	claims := oidc.IDTokenClaims{
		Issuer:     "https://auth.example.com",
		Subject:    "user_123",
		Audience:   "client-abc",
		Expiration: time.Now().Add(time.Hour).Unix(),
		IssuedAt:   time.Now().Unix(),
		Nonce:      "nonce-correct",
	}

	signed, _ := oidc.SignIDToken(claims, secret)

	// Test issuer mismatch
	if _, err := oidc.ParseAndVerifyIDToken(signed, secret, "https://wrong.com", claims.Audience, claims.Nonce, time.Now()); err == nil {
		t.Errorf("expected issuer mismatch error")
	}

	// Test audience mismatch
	if _, err := oidc.ParseAndVerifyIDToken(signed, secret, claims.Issuer, "client-wrong", claims.Nonce, time.Now()); err == nil {
		t.Errorf("expected audience mismatch error")
	}

	// Test nonce mismatch
	if _, err := oidc.ParseAndVerifyIDToken(signed, secret, claims.Issuer, claims.Audience, "nonce-wrong", time.Now()); err == nil {
		t.Errorf("expected nonce mismatch error")
	}
}

func TestOAuth2_FullFlowAndPKCEInterception(t *testing.T) {
	key := []byte("secret-key-12345678901234567890")
	as := server.NewAuthorizationServer("https://auth.example.com", key)
	as.RegisterClient("client_1", "https://app.com/cb")

	cli := client.NewClient("client_1", "https://app.com/cb", as, key)
	ch, err := cli.BuildAuthorizationRequest("openid profile")
	if err != nil {
		t.Fatalf("BuildAuthorizationRequest failed: %v", err)
	}

	ac, err := as.Authorize("client_1", "https://app.com/cb", "openid profile", "user_sub_1", ch, "S256", cli.Nonce)
	if err != nil {
		t.Fatalf("Authorize failed: %v", err)
	}

	// Attacker tries to exchange intercepted code with wrong verifier
	_, err = as.ExchangeCode(ac.Code, "client_1", "https://app.com/cb", "attacker_wrong_verifier_12345678901234567890123")
	if err == nil {
		t.Fatalf("expected PKCE failure for attacker")
	}

	// Legitimate client exchanges code
	resp, err := cli.Exchange(ac.Code)
	if err != nil {
		t.Fatalf("Legitimate client exchange failed: %v", err)
	}

	if resp.AccessToken == "" || resp.RefreshToken == "" || resp.IDToken == "" {
		t.Fatalf("expected all tokens returned")
	}

	// Reusing the same authorization code must fail
	_, err = cli.Exchange(ac.Code)
	if err == nil {
		t.Fatalf("expected auth code already redeemed error")
	}
}

func TestOAuth2_RefreshTokenRotation_AndReplayDetection(t *testing.T) {
	key := []byte("secret-key-12345678901234567890")
	as := server.NewAuthorizationServer("https://auth.example.com", key)
	as.RegisterClient("client_1", "https://app.com/cb")

	cli := client.NewClient("client_1", "https://app.com/cb", as, key)
	ch, _ := cli.BuildAuthorizationRequest("openid")
	ac, _ := as.Authorize("client_1", "https://app.com/cb", "openid", "user_sub_1", ch, "S256", cli.Nonce)
	_, _ = cli.Exchange(ac.Code)

	firstRefreshToken := cli.RefreshToken

	// Normal refresh (rotation)
	resp2, err := cli.RefreshTokens()
	if err != nil {
		t.Fatalf("RefreshTokens failed: %v", err)
	}

	if resp2.RefreshToken == firstRefreshToken {
		t.Errorf("refresh token was not rotated")
	}

	// Attacker reuses firstRefreshToken (already consumed)
	_, err = as.Refresh(firstRefreshToken, "client_1")
	if err == nil || !strings.Contains(err.Error(), "reuse detected") {
		t.Fatalf("expected replay detection error, got: %v", err)
	}

	// Legitimate client attempts to use active rotated token -> must be revoked because entire family was compromised!
	_, err = cli.RefreshTokens()
	if err == nil || !strings.Contains(err.Error(), "revoked") {
		t.Fatalf("expected token family revoked error, got: %v", err)
	}
}

func TestPKCE_Plain_Method(t *testing.T) {
	pair, err := pkce.GeneratePKCEPair("plain")
	if err != nil {
		t.Fatalf("GeneratePKCEPair plain failed: %v", err)
	}
	if pair.CodeChallenge != pair.CodeVerifier {
		t.Fatalf("plain challenge must equal verifier")
	}
	if err := pkce.Verify(pair.CodeVerifier, pair.CodeChallenge, "plain"); err != nil {
		t.Fatalf("Verify failed for plain method: %v", err)
	}
}

func TestOIDC_MalformedJWT(t *testing.T) {
	secret := []byte("test-signing-key-123456789012")
	if _, err := oidc.ParseAndVerifyIDToken("invalid.jwt", secret, "iss", "aud", "", time.Now()); err == nil {
		t.Errorf("expected malformed jwt error for 2-part string")
	}
	if _, err := oidc.ParseAndVerifyIDToken("bad!header.bad!payload.bad!sig", secret, "iss", "aud", "", time.Now()); err == nil {
		t.Errorf("expected malformed jwt error for non-base64 input")
	}
}

func TestOAuth2_NegativePaths(t *testing.T) {
	key := []byte("secret-key-12345678901234567890")
	as := server.NewAuthorizationServer("https://auth.example.com", key)
	as.RegisterClient("client_1", "https://app.com/cb")

	// 1. Authorize with bad client
	if _, err := as.Authorize("bad_client", "https://app.com/cb", "openid", "sub1", "ch", "S256", "n"); err == nil {
		t.Errorf("expected error for unregistered client")
	}

	// 2. Authorize with bad redirect uri
	if _, err := as.Authorize("client_1", "https://evil.com/cb", "openid", "sub1", "ch", "S256", "n"); err == nil {
		t.Errorf("expected error for wrong redirect uri")
	}

	// 3. Authorize with empty challenge
	if _, err := as.Authorize("client_1", "https://app.com/cb", "openid", "sub1", "", "S256", "n"); err == nil {
		t.Errorf("expected error for empty challenge")
	}

	// 4. Authorize with invalid challenge method
	if _, err := as.Authorize("client_1", "https://app.com/cb", "openid", "sub1", "ch", "INVALID_METHOD", "n"); err == nil {
		t.Errorf("expected error for invalid challenge method")
	}

	// Valid authorize
	cli := client.NewClient("client_1", "https://app.com/cb", as, key)
	ch, _ := cli.BuildAuthorizationRequest("openid profile read:data")
	ac, err := as.Authorize("client_1", "https://app.com/cb", "openid profile read:data", "user1", ch, "S256", cli.Nonce)
	if err != nil {
		t.Fatalf("Authorize failed: %v", err)
	}

	// 5. Exchange with invalid code
	if _, err := as.ExchangeCode("non_existent_code", "client_1", "https://app.com/cb", cli.Verifier); err == nil {
		t.Errorf("expected error for non-existent code")
	}

	// 6. Exchange with mismatched clientID / redirectURI
	if _, err := as.ExchangeCode(ac.Code, "wrong_client", "https://app.com/cb", cli.Verifier); err == nil {
		t.Errorf("expected error for mismatched client ID")
	}
	if _, err := as.ExchangeCode(ac.Code, "client_1", "https://wrong.com/cb", cli.Verifier); err == nil {
		t.Errorf("expected error for mismatched redirect URI")
	}

	// Legitimate exchange
	resp, err := cli.Exchange(ac.Code)
	if err != nil {
		t.Fatalf("Exchange failed: %v", err)
	}

	// 7. Validate access token scopes
	if _, err := as.ValidateAccessToken(resp.AccessToken, "read:data"); err != nil {
		t.Errorf("expected scope read:data to pass: %v", err)
	}
	if _, err := as.ValidateAccessToken(resp.AccessToken, "write:admin"); err == nil {
		t.Errorf("expected scope check to fail for ungranted write:admin")
	}
	if _, err := as.ValidateAccessToken("invalid_token", "read:data"); err == nil {
		t.Errorf("expected error for nonexistent token")
	}

	// 8. Refresh negative cases
	if _, err := as.Refresh("non_existent_rt", "client_1"); err == nil {
		t.Errorf("expected error for non-existent refresh token")
	}
	if _, err := as.Refresh(resp.RefreshToken, "wrong_client"); err == nil {
		t.Errorf("expected error for refresh with wrong client ID")
	}
}

func TestOAuth2_ConcurrencyAndRace(t *testing.T) {
	key := []byte("secret-key-12345678901234567890")
	as := server.NewAuthorizationServer("https://auth.example.com", key)
	as.RegisterClient("client_conc", "https://app.com/cb")

	var wg sync.WaitGroup
	workers := 20

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			cli := client.NewClient("client_conc", "https://app.com/cb", as, key)
			ch, err := cli.BuildAuthorizationRequest("openid")
			if err != nil {
				t.Errorf("BuildAuthorizationRequest failed: %v", err)
				return
			}

			rawBytes := make([]byte, 8)
			rand.Read(rawBytes)
			sub := "user_" + hex.EncodeToString(rawBytes)

			ac, err := as.Authorize("client_conc", "https://app.com/cb", "openid", sub, ch, "S256", cli.Nonce)
			if err != nil {
				t.Errorf("Authorize failed: %v", err)
				return
			}

			resp, err := cli.Exchange(ac.Code)
			if err != nil {
				t.Errorf("Exchange failed: %v", err)
				return
			}

			_, err = as.ValidateAccessToken(resp.AccessToken, "openid")
			if err != nil {
				t.Errorf("ValidateAccessToken failed: %v", err)
				return
			}

			_, err = cli.RefreshTokens()
			if err != nil {
				t.Errorf("RefreshTokens failed: %v", err)
				return
			}
		}(i)
	}

	wg.Wait()
}
