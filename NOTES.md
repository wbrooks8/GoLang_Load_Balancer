# Learning Notes

## 1. What this project teaches

This project demonstrates the fundamentals of a reverse proxy and simple load balancer.

The core idea is:

- a client calls the load balancer
- the balancer chooses an upstream server
- the balancer forwards the request
- the upstream server responds to the client as if the balancer were the original target

## 2. Round robin

Round robin means: take turns choosing servers in a cycle.

If there are three servers:

- request 1 -> server A
- request 2 -> server B
- request 3 -> server C
- request 4 -> server A again

This is a simple and fair way to distribute traffic.

## 3. Health checks

A load balancer should not keep sending traffic to a server that is down.

The pattern here is:

- periodically check server health
- track whether a target is alive
- skip unhealthy servers when choosing a backend

This is a key part of real load balancing.

## 4. Concurrency in Go

The `LoadBalancer` uses a mutex to avoid races when multiple requests access the round-robin counter.

This matters because Go programs often serve many requests at the same time. Without locking, shared state could be read and written concurrently in unsafe ways.

## 5. Reverse proxy

The `httputil.ReverseProxy` in Go is the mechanism that forwards incoming HTTP traffic to a backend server while preserving request semantics.

This is one of the main building blocks of a simple HTTP load balancer.

## 6. Why this is a learning app

This project is intentionally simple. That makes it easy to reason about, but it also means some production concerns are deliberately left out.

Examples of missing production features:

- weighted routing
- sticky sessions
- circuit breakers
- queueing and retries
- metrics and tracing
- TLS configuration
- large-scale health orchestration

## 7. Good next improvements

If you continue developing this project, good next steps are:

- add a health-check configuration value instead of hardcoded timing
- add tests around failure scenarios
- add logging for upstream selection
- improve server status tracking
- add a second balancing strategy such as least connections or weighted round robin

## 8. Why tests matter here

Without tests, it is easy to miss edge cases such as:

- all servers are unhealthy
- a server goes offline after startup
- round robin unexpectedly repeats the same server
- the balancer deadlocks or miscounts requests

Tests help turn a demonstration app into a reliable system you can reason about.
