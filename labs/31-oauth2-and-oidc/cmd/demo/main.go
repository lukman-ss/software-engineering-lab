package main

import (
	"fmt"
	"os"

	"labs/31-oauth2-and-oidc/pkg/client"
	"labs/31-oauth2-and-oidc/pkg/server"
)

func main() {
	fmt.Println("=== LAB 31: OAuth 2.0 & OpenID Connect (OIDC) Demo ===")
	fmt.Println()

	signingKey := []byte("super-secret-hmac-key-32-bytes!!")
	as := server.NewAuthorizationServer("https://auth.example.com", signingKey)

	clientID := "spa-client-123"
	redirectURI := "https://app.example.com/callback"
	as.RegisterClient(clientID, redirectURI)

	cli := client.NewClient(clientID, redirectURI, as, signingKey)

	fmt.Println("[Step 1] Initiate Authorization Request with PKCE & OIDC scope...")
	challenge, err := cli.BuildAuthorizationRequest("openid profile email")
	if err != nil {
		fmt.Printf("Failed to build auth req: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(" Generated PKCE code_challenge (S256): %s\n", challenge)
	fmt.Printf(" Generated OIDC nonce: %s\n", cli.Nonce)
	fmt.Println()

	fmt.Println("[Step 2] Authorization Server issues Authorization Code...")
	authCode, err := as.Authorize(clientID, redirectURI, "openid profile email", "user_alice_99", challenge, "S256", cli.Nonce)
	if err != nil {
		fmt.Printf("Authorize failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(" Issued Code: %s (expires: %s)\n", authCode.Code, authCode.ExpiresAt.Format("15:04:05"))
	fmt.Println()

	fmt.Println("[Step 3] Client exchanges Code + code_verifier for Tokens...")
	resp, err := cli.Exchange(authCode.Code)
	if err != nil {
		fmt.Printf("Exchange failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(" Access Token : %s\n", resp.AccessToken)
	fmt.Printf(" Refresh Token: %s\n", resp.RefreshToken)
	fmt.Printf(" ID Token     : %s...\n", resp.IDToken[:30])
	fmt.Println()

	fmt.Println("[Step 4] Validate OIDC ID Token claims...")
	fmt.Printf(" ID Token Subject  : %s\n", cli.IDClaims.Subject)
	fmt.Printf(" ID Token Issuer   : %s\n", cli.IDClaims.Issuer)
	fmt.Printf(" ID Token Audience : %s\n", cli.IDClaims.Audience)
	fmt.Printf(" ID Token Nonce Match: SUCCESS (%s)\n", cli.IDClaims.Nonce)
	fmt.Println()

	fmt.Println("[Step 5] Access Resource Server using Access Token...")
	sub, err := as.ValidateAccessToken(cli.AccessToken, "profile")
	if err != nil {
		fmt.Printf("Resource access failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(" Resource Server verified Access Token. Subject: %s\n", sub)
	fmt.Println()

	fmt.Println("[Step 6] Test Interception Attack (Attempt code exchange with wrong verifier)...")
	attackerCode, _ := as.Authorize(clientID, redirectURI, "openid", "user_alice_99", challenge, "S256", cli.Nonce)
	_, err = as.ExchangeCode(attackerCode.Code, clientID, redirectURI, "invalid_attacker_code_verifier_string_1234567890123")
	if err != nil {
		fmt.Printf(" Interception Attack Blocked! Error: %v\n", err)
	} else {
		fmt.Println(" ERROR: Interception attack succeeded unexpectedly!")
		os.Exit(1)
	}
	fmt.Println()

	fmt.Println("[Step 7] Test Refresh Token Rotation...")
	stolenRefresh := cli.RefreshToken
	newTokens, err := cli.RefreshTokens()
	if err != nil {
		fmt.Printf("Refresh failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf(" Rotated New Access Token : %s\n", newTokens.AccessToken)
	fmt.Printf(" Rotated New Refresh Token: %s\n", newTokens.RefreshToken)
	fmt.Println()

	fmt.Println("[Step 8] Test Refresh Token Replay Detection (Attacker reuses stolen Refresh Token)...")
	_, err = as.Refresh(stolenRefresh, clientID)
	if err != nil {
		fmt.Printf(" Refresh Token Replay Detected! Error: %v\n", err)
	} else {
		fmt.Println(" ERROR: Refresh token replay was not detected!")
		os.Exit(1)
	}

	fmt.Println(" Re-testing active refresh token after family revocation...")
	_, err = cli.RefreshTokens()
	if err != nil {
		fmt.Printf(" Active session revoked due to family revocation! Error: %v\n", err)
	} else {
		fmt.Println(" ERROR: Active session was not revoked after family reuse detection!")
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("=== Demo Completed Successfully ===")
}
