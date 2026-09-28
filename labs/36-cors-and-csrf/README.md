# Lab 36: CORS & CSRF — Securing Backend from Cross-Origin Attacks

This lab provides an implementation and empirical proof demonstrating why **CORS is not a backend firewall or CSRF protection mechanism**, and how robust defense-in-depth measures against CSRF operate.

## Architecture

- **`internal/cors`**: Spec-compliant CORS middleware (`OPTIONS` preflight, allowed origins, method/header safelists, credential checks).
- **`internal/csrf`**: Anti-CSRF mechanisms including HMAC-SHA256 signed session-bound tokens, Fetch Metadata (`Sec-Fetch-Site`), and API custom header middleware.
- **`internal/bank`**: Bank application service simulating cookie-authenticated balance inquiries, vulnerable transfer endpoints, and protected transfer endpoints.
- **`cmd/demo`**: Runnable CLI program showcasing attacks against vulnerable vs. protected configurations.
- **`tests`**: Integration test suite verifying cross-origin requests, preflight, race safety, and attack mitigation.

## Running Tests

Execute all unit and integration tests:

```bash
go test -v ./...
```

Execute tests with the Go race detector:

```bash
go test -race ./...
```

## Running Demo

Run the end-to-end demonstration:

```bash
go run ./cmd/demo
```
