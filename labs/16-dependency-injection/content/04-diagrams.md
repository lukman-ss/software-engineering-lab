# Diagrams

## 1. Constructor Injection Pattern

```text
main.go (composition root)
       |
       | constructs
       v
  +-------------+         implements        +----------------+
  | RealGateway | <----------------------- | PaymentGateway |
  +-------------+                          +----------------+
        |                                           |
        | injects via NewProcessor                  | used by
        v                                           v
  +-------------+        ProcessPayment()      +----------------+
  |  Processor  | --------------------------> | PaymentGateway |
  +-------------+        (delegates)          +----------------+
        |
        | creates
        v
    Money{...}
```

## 2. Service Locator Anti-Pattern

```text
main.go (composition root)
       |
       | constructs
       v
  +-------------+         implements        +----------------+
  | RealGateway | <----------------------- | PaymentGateway |
  +-------------+                          +----------------+
        |                                           |
        | wraps into                                |
        v                                           v
  +----------------+    GetPaymentGateway()    +----------------+
  | SimpleContainer| ------------------------> | Container      |
  +----------------+    (interface)            +----------------+
        |                                            |
        | injects via NewBadProcessor                | queries
        v                                            v
  +----------------+      ProcessPayment()      +----------------+
  | BadProcessor   | ------------------------> | Container      |
  +----------------+      (delegates)          +----------------+
                                                       |
                                                       | gets
                                                       v
                                                  +----------------+
                                                  | PaymentGateway |
                                                  +----------------+
```

## 3. Test Isolation via Mock Injection

```text
tests/processor_test.go
       |
       | constructs
       v
  +-------------+         implements        +----------------+
  |  MockGateway | <---------------------- | PaymentGateway |
  +-------------+                          +----------------+
        |                                           |
        | injects via NewProcessor                  | used by
        v                                           v
  +-------------+        ProcessPayment()      +----------------+
  |  Processor  | --------------------------> | PaymentGateway |
  +-------------+        (delegates)          +----------------+
        |
        | verifies state
        v
  Assertions (PASS/FAIL)
```

## 4. Dependency Flow Summary

```text
┌─────────────────────────────────────────────────────────────────────┐
│ Composition Root (main.go)                                          │
│  - instantiate concrete implementations (RealGateway)               │
│  - wire dependencies (NewProcessor(realGateway))                    │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              │ injects
                              v
┌─────────────────────────────────────────────────────────────────────┐
│ Consumer (Processor / BadProcessor)                                 │
│  - depends only on interface (PaymentGateway)                       │
│  - no knowledge of concrete implementation                          │
│  - validates input BEFORE calling gateway                           │
└─────────────────────────────────────────────────────────────────────┘
                              │
                              │ calls
                              v
┌─────────────────────────────────────────────────────────────────────┐
│ Infrastructure (RealGateway / MockGateway)                          │
│  - implements concrete behavior (HTTP, mock, etc.)                  │
└─────────────────────────────────────────────────────────────────────┘
```

## 5. Test Coverage Map

```text
Test Cases                              | Constructor Injection | Service Locator
----------------------------------------|----------------------|------------------
TestProcessor_Success /                 | ✅ PASS               | ✅ PASS
TestBadProcessor_Success                | (amount=50/75)        | (amount=50/75)
                                        |                       |
TestProcessor_GatewayError /            | ✅ PASS               | ✅ PASS
TestBadProcessor_GatewayError           | (ShouldFail=true)     | (ShouldFail=true)
                                        |                       |
TestProcessor_InvalidAmount /           | ✅ PASS               | ✅ PASS
TestBadProcessor_InvalidAmount          | (amount<0)            | (amount<0)
                                        | gateway NOT called    | gateway NOT called
                                        | (ChargedMoney=0)      | (ChargedMoney=0)
```
