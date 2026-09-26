# Research Plan

## Research Topic
Rate Limiting & Backpressure - Sistem Sehat Harus Bisa Bilang "Cukup"

## Objective
Investigate rate limiting and backpressure mechanisms in distributed systems, focusing on practical engineering approaches, algorithms (Token Bucket, Exponential Backoff with Jitter), queue management, and their role in system scalability and fault isolation.

## Research Questions
1. What are the primary rate limiting algorithms and their practical implementations?
2. How does backpressure differ from rate limiting, and how is it implemented in queue-based systems?
3. What are the best practices for implementing fair resource allocation in multi-tenant systems?
4. How do Little's Law and queue metrics inform capacity planning and operational alerts?
5. What are common anti-patterns in rate limiting and backpressure, and how can they be avoided?

## Search Strategy
- Search for official documentation on rate limiting algorithms (Token Bucket, Leaky Bucket, Sliding Window)
- Look for academic papers on backpressure in distributed systems
- Find industry best practices from cloud providers (AWS, GCP, Azure)
- Search for practical implementations in popular open-source projects (nginx, Envoy, Redis, Kafka)
- Investigate queue management patterns in message brokers

## Expected Primary Sources
- API gateway documentation (nginx, Envoy Gateway)
- Message broker documentation (Kafka, RabbitMQ, Redis Streams)
- Academic papers on queueing theory and backpressure
- Cloud provider best practices for rate limiting
- Open-source implementations of rate limiting libraries

## Risks / Unknowns
- Some implementation details may be proprietary
- Language barrier: Many authoritative sources in English
- Some claims in the topic spec need verification (e.g., exact calculations, specific algorithms)
