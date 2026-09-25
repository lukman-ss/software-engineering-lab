# 03 Core Concepts

## Circuit Breaker
A proxy wrapping a fragile dependency. Tracks failure rates. Trips OPEN upon reaching a threshold, rejecting future calls immediately without touching the network.

## Goal
To limit the blast radius of dependency failures, fail-fast, preserve caller resources, and allow downstream recovery.

## Why it's not a Retry
Retries repeatedly hammer the failing downstream. Circuit breakers actively block traffic from ever leaving the caller node.
