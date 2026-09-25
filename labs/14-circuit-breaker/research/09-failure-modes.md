# 09 Failure Modes

1. **Premature Opening**: Threshold set too low; transient hiccups trip the circuit unnecessarily.
2. **Thundering Herd during Probe**: Allowing too many concurrent requests through during HALF_OPEN knocks down the recovering service.
3. **Infinite Open**: Cooldown is too long or health checks never pass.
4. **4xx False Positives**: Tripping the breaker on user-input validation errors (400, 404, 422) instead of system faults (500, 502, 503, 504, timeouts).
