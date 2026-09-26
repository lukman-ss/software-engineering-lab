# 08 Observability

Crucial metrics to track circuit state (architectural recommendation for production; omitted in this minimal reference codebase):
- `circuit_state`: Enum (CLOSED=0, OPEN=1, HALF_OPEN=2).
- `circuit_open_count`: Counter indicating how often circuit flipped OPEN.
- `failure_count`: Recent consecutive failures or rolling failure sum.
- `rejected_call_count`: Count of fast-failed requests while OPEN.
- `dependency_latency`: Request durations.
