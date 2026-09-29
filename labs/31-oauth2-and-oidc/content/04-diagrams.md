# Diagrams

Berikut adalah diagram yang didasarkan pada implementasi dan dokumentasi lab.

## 1. Authorization Code Flow with PKCE

```text
Client                Authorization Server              Resource Server
  |                         |                               |
  |---(1) Auth Request------>|                               |
  |    response_type=code   |                               |
  |    client_id            |                               |
  |    redirect_uri         |                               |
  |    scope=openid...      |                               |
  |    code_challenge=S256  |                               |
  |    nonce                |                               |
  |                         |                               |
  |<--(2) Auth Code---------|                               |
  |    code=abc123          |                               |
  |                         |                               |
  |---(3) Token Exchange---->|                               |
  |    grant_type=auth_code |                               |
  |    code=abc123          |                               |
  |    code_verifier=xyz    |                               |
  |    redirect_uri         |                               |
  |                         |                               |
  |<--(4) Tokens------------|                               |
  |    access_token         |                               |
  |    refresh_token        |                               |
  |    id_token (JWT)       |                               |
  |                         |                               |
  |---(5) API Request------>|------------------------------>|
  |    Bearer access_token  |     (validates scope)         |
  |                         |                               |
  |<--(6) Response----------|<------------------------------|
  |    protected data       |                               |
```

## 2. PKCE Challenge-Response Verification

```text
Client Side:                          Server Side:
                                          
  code_verifier (random 32 bytes)         
        |                                  
        v                                  
  SHA-256(code_verifier)                   
        |                                  
        v                                  
  code_challenge (base64url)               
        |                                  
        +--- Authorization Request -------->+
        |    code_challenge=S256           |
        |                                  |
        |                                  |+ Store code_challenge
        |                                  |+ Bind to auth code
        |                                  |
        |<-- Authorization Code ------------+
        |    code=AUTH_CODE                
        |                                  
        +--- Token Exchange --------------->+
        |    code=AUTH_CODE                |
        |    code_verifier=VERIFIER        |
        |                                  |
        |                                  |+ SHA-256(verifier)
        |                                  |+ Compare with stored challenge
        |                                  |
        |<-- Access Token + ID Token -------+
```

## 3. ID Token Validation Flow

```text
Client receives ID Token (JWT):
  header.claims.signature
        |
        v
  [1] Decode & parse JWT structure
        |
        v
  [2] Verify HMAC-SHA256 signature using shared key
        |
        v
  [3] Validate iss (exact match with expected issuer)
        |
        v
  [4] Validate aud (must contain client_id)
        |
        v
  [5] Validate exp (must be in future)
        |
        v
  [6] Validate iat (must not be > now + 5 minutes)
        |
        v
  [7] Validate nonce (must match original request nonce)
        |
        v
  [8] Return validated claims (sub, email, etc.)
```

## 4. Refresh Token Rotation and Replay Detection

```text
Normal Refresh:                        Replay Detection:
                                        
  Client has RT_A                       Attacker has RT_A (stolen)
                                        
  Client POSTs RT_A                     Attacker POSTs RT_A
        |                                   |
        v                                   v
  Server checks RT_A exists               Server checks RT_A exists
        |                                   |
        v                                   v
  RT_A.Revoked = false                    RT_A.Revoked = false
        |                                   |
        v                                   v
  Server marks RT_A.Revoked = true        Server detects RT_A.Revoked = true
        |                                   |
        v                                   v
  Server creates RT_B (new)               Server flags FamilyID as compromised
        |                                   |
        v                                   v
  Returns new AT_B + RT_B                 Returns ErrTokenReplayDetected
                                        
  Next refresh uses RT_B                  Client must re-authenticate
```

## 5. Token Family Revocation

```text
FamilyID = fam_abc123
        |
        +-- RT_A (Revoked=true) ---+
        |                          |
        +-- AT_A (Active) -------->+--- Revoked when replay detected
        |                          |
        +-- RT_B (Revoked=false) -+
                                   |
                            Replay of RT_A triggers:
                            revokedFams["fam_abc123"] = true
                            All tokens in family invalidated
```
