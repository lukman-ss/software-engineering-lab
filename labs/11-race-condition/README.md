# Go Memory Race Lab

Sub-lab for investigating Go memory data races, improper shared state access, and concurrency primitives.

## Files
- `datarace.go`: Contains `UnsafeCounter`, `MutexCounter`, `AtomicCounter`, and `ChannelCounter`.
- `datarace_test.go`: Tests concurrent access and race detection.

## Synchronization Primitives Comparison

1. **`sync.Mutex` / `sync.RWMutex`**
   - **When to use**: Protecting critical sections spanning multiple operations or complex composite state structures (e.g., maps, structs with multiple fields).
   - **Trade-offs**: Can introduce lock contention, blocking goroutines. Risk of deadlocks if nested improperly.

2. **`sync/atomic`**
   - **When to use**: Extremely fast, lock-free synchronization for single primitive variables (counters, flags, pointers).
   - **Trade-offs**: Limited to primitive types and simple operations. Harder to coordinate complex invariants across multiple variables.

3. **Channels (CSP / Ownership Transfer)**
   - **When to use**: Passing ownership of data between goroutines, event signaling, or maintaining isolated state where only one goroutine owns the variable (actor model).
   - **Trade-offs**: Higher overhead due to channel allocation and potential goroutine blocking/deadlocking if unbuffered or improperly coordinated. Can make simple counting slower than mutex/atomic.
