# Diagrams

## 1. Constructor Injection Pattern

```text
+-------------+         constructs         +---------------+
|   main.go   | -------------------------> |  RealGateway  |
+-------------+                            +---------------+
       |                                           |
       | injects via NewProcessor                  | implements PaymentGateway
       v                                           v
+-------------+       calls Charge()       +----------------+
|  Processor  | -------------------------> | PaymentGateway |
+-------------+                            +----------------+
```

## 2. Service Locator Anti-Pattern

```text
+----------------+        constructs         +---------------+
|    main.go     | ------------------------> |  RealGateway  |
+----------------+                           +---------------+
        |                                            |
        | wraps into                                 |
        v                                            v
+----------------+     injects via           +---------------+
| SimpleContainer| ------------------------> | BadProcessor  |
+----------------+    NewBadProcessor        +---------------+
                                                     |
                                                     | queries GetPaymentGateway()
                                                     v
                                             +----------------+
                                             | PaymentGateway |
                                             +----------------+
```

## 3. Test Isolation via Mock Injection

```text
+---------------------+        constructs         +---------------+
| tests/processor_test| ------------------------> |  MockGateway  |
+---------------------+                           +---------------+
           |                                              |
           | injects via NewProcessor                     | implements PaymentGateway
           v                                              v
+---------------------+       calls Charge()      +----------------+
|      Processor      | ------------------------> | PaymentGateway |
+---------------------+                           +----------------+
           |
           | verifies ChargedMoney state
           v
+---------------------+
| Assertions Pass/Fail|
+---------------------+
```
