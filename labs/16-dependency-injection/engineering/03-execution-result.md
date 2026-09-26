# Execution Result

## Build
Command: `go build ./...`
Result: Success

## Tests
Command: `go test -v -count=1 ./...`
Result: 
```
?   	lab16/cmd/demo	[no test files]
?   	lab16/internal/di	[no test files]
=== RUN   TestProcessor_Success
--- PASS: TestProcessor_Success (0.00s)
=== RUN   TestProcessor_GatewayError
--- PASS: TestProcessor_GatewayError (0.00s)
=== RUN   TestProcessor_InvalidAmount
--- PASS: TestProcessor_InvalidAmount (0.00s)
=== RUN   TestBadProcessor_Success
--- PASS: TestBadProcessor_Success (0.00s)
=== RUN   TestBadProcessor_GatewayError
--- PASS: TestBadProcessor_GatewayError (0.00s)
=== RUN   TestBadProcessor_InvalidAmount
--- PASS: TestBadProcessor_InvalidAmount (0.00s)
PASS
ok  	lab16/tests	0.387s
```

## Race Detector
Command: `go test -race ./...`
Result:
```
?   	lab16/cmd/demo	[no test files]
?   	lab16/internal/di	[no test files]
ok  	lab16/tests	1.356s
```

## Demo
Command: `go run ./cmd/demo`
Result:
```
--- Running Constructor Injection ---
RealGateway charging 100 USD
--- Running Service Locator ---
RealGateway charging 200 USD
```

## Final Engineering Status
READY_FOR_ENGINEERING_AUDIT
