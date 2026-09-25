# Content Brief

Topic: Deadlock in Concurrent Systems
Target Reader: Software Engineer
Problem: Deadlocks form permanent circular waits blocking transactions, requiring system-level aborts and application-level handling to resolve or prevent.
Core Mental Model: Deadlock is a circular dependency. The database aborts one as a "victim". Lock ordering prevents cycles from forming. Application retry recovers aborted victims.
Approved Research Status: APPROVED
Approved Engineering Status: APPROVED
Main Concepts: Circular Wait, Deadlock Monitor/Victim, Lock Ordering, Transaction Duration, Application-Level Retry.
Verified Behaviors: Naive bidirectional resource access causes deadlocks. Consistent lock ordering completely prevents deadlocks. Retry mechanisms can successfully recover from deadlock aborts. Longer transactions increase deadlock frequency.
Available Case Studies: N/A (Bank Transfer simulation used)
Warnings: The lab is an application-level simulation; actual databases use Wait-For Graphs (WFG). Self-transfers (`from == to`) can self-deadlock due to missing validation in the lab code.
