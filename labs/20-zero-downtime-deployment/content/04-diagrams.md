# Architecture & Process Diagrams

## 1. System Architecture

```text
               +-----------------------------+
               |       Load Balancer         |
               +-----------------------------+
                       |              |
     /healthz/ready ok |              | In-flight HTTP Traffic
                       v              v
     +-------------------------------------------------+
     |                 internal/server                 |
     |                                                 |
     |  - /healthz/live   : Indicates process status   |
     |  - /healthz/ready  : Controls traffic routing   |
     |  - /work           : In-flight HTTP requests    |
     |  - preStop Hook    : Network update delay       |
     +-------------------------------------------------+
                       |              |
           User Writes |              | Fallback Reads
                       v              v
     +-------------------------------------------------+
     |                  internal/db                    |
     |                                                 |
     |  - UserStore                                    |
     |  - Legacy Field  : Name                         |
     |  - Modern Fields : FirstName, LastName          |
     |  - Expand / Contract Compatibility Logic        |
     +-------------------------------------------------+
                              ^
                              | Updates/Reads
     +-------------------------------------------------+
     |                 internal/worker                 |
     |                                                 |
     |  - jobChan         : Buffered in-memory queue   |
     |  - Goroutines      : Background processing      |
     |  - Stop()          : Finish active job cleanly  |
     +-------------------------------------------------+
```

---

## 2. Graceful Shutdown & PreStop Lifecycle

```text
Orchestrator               HTTP Server                  In-Flight Client
     |                          |                              |
     |--- Send SIGTERM -------->|                              |
     |                          |-- SetReady(false)            |
     |                          |   (detach from LB)           |
     |                          |                              |
     |                          |-- Sleep(preStop) ----------->| (Allows network routes
     |                          |   (wait routing tables)      |  to propagate to proxy)
     |                          |                              |
     |                          |-- srv.Shutdown()             |
     |                          |   (stop accepting new reqs)  |
     |                          |                              |
     |                          |                              |<-- Process ongoing req
     |                          |                              |--- Return 200 OK ---->
     |                          |-- wg.Wait()                  |
     |                          |   (all in-flight done)       |
     |<-- Process Exits (0) ----|                              |
```

---

## 3. Database Expand and Contract Pattern

```text
Phase 1: Old Schema
[ ID ] [ Name ]

Phase 2: Expand (Both Exist simultaneously)
[ ID ] [ Name ] [ FirstName ] [ LastName ]
   ^                 ^              ^
   |                 |              |
Writes legacy     Dual Write     Writes new
Reads legacy      Fallback Read  Reads new

Phase 3: Contract (Old Schema Retired)
[ ID ] [ FirstName ] [ LastName ]
```
