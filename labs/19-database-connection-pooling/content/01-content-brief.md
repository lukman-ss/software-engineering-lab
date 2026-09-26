# Content Brief

Topic: Database Connection Pooling, Sizing, and Exhaustion
Target Reader: Backend Engineers
Problem: Misconception that larger connection pools increase database throughput, leading to connection exhaustion, latency spikes, and server crashes.
Core Mental Model: Connection pools should match hardware capacity (CPU cores/disk), not application concurrency.
Approved Research Status: APPROVED
Approved Engineering Status: APPROVED
Main Concepts: Connection overhead, pool sizing limits (core_count * 2), connection leaks, pool starvation.
Verified Behaviors: Direct connection penalty, oversized pool rejection, connection leak/starvation timeouts.
Available Case Studies: Simulated latency overhead, server rejection on oversized pool, pool starvation via long network call leak.
Warnings: Hardware saturation knee; real-world context switching/memory allocation not fully simulated in lab.