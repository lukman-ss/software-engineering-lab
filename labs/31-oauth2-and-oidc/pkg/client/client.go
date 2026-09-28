package client

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"labs/31-oauth2-and-oidc/pkg/oidc"
	"labs/31-oauth2-and-oidc/pkg/pkce"
	"labs/31-oauth2-and-oidc/pkg/server"
)

type Client struct {
	ClientID     string
	RedirectURI  string
	Server       *server.AuthorizationServer
	SigningKey   []byte
	Verifier     string
	State        string
	Nonce        string
	AccessToken  string
	RefreshToken string
	IDClaims     *oidc.IDTokenClaims
}

func NewClient(clientID, redirectURI string, as *server.AuthorizationServer, signingKey []byte) *Client {
	return &Client{
		ClientID:    clientID,
		RedirectURI: redirectURI,
		Server:      as,
		SigningKey:  signingKey,
	}
}

func (c *Client) BuildAuthorizationRequest(scope string) (string, error) {
	pair, err := pkce.GeneratePKCEPair("S256")
	if err != nil {
		return "", err
	}
	c.Verifier = pair.CodeVerifier

	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", err
	}
	c.State = hex.EncodeToString(stateBytes)

	nonceBytes := make([]byte, 16)
	if _, err := rand.Read(nonceBytes); err != nil {
		return "", err
	}
	c.Nonce = hex.EncodeToString(nonceBytes)

	return pair.CodeChallenge, nil
}

func (c *Client) Exchange(code string) (*server.TokenResponse, error) {
	resp, err := c.Server.ExchangeCode(code, c.ClientID, c.RedirectURI, c.Verifier)
	if err != nil {
		return nil, err
	}

	c.AccessToken = resp.AccessToken
	c.RefreshToken = resp.RefreshToken

	if resp.IDToken != "" {
		claims, err := oidc.ParseAndVerifyIDToken(resp.IDToken, c.SigningKey, c.Server.Issuer, c.ClientID, c.Nonce, time.Now())
		if err != nil {
			return nil, fmt.Errorf("id_token validation failed: %w", err)
		}
		c.IDClaims = claims
	}

	return resp, nil
}

func (c *Client) RefreshTokens() (*server.TokenResponse, error) {
	resp, err := c.Server.Refresh(c.RefreshToken, c.ClientID)
	if err != nil {
		return nil, err
	}
	c.AccessToken = resp.AccessToken
	c.RefreshToken = resp.RefreshToken
	return resp, nil
}
