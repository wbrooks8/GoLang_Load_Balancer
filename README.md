# Go Load Balancer

This project is a small Go-based load balancer built as a learning exercise. It demonstrates a round-robin proxy pattern, upstream server selection, and simple health checking.

## What the app does

- Starts a local HTTP server on a configured port
- Accepts incoming requests
- Selects an upstream server using round-robin
- Forwards the request to that server through a reverse proxy
- Skips unhealthy servers when possible
- Returns a 503 when no healthy servers are available

## Project layout

- `src/main.go` – application startup and routing
- `src/interface/server.go` – server interface contract
- `src/models/server.go` – upstream server implementation and health checks
- `src/models/load_balancer.go` – round-robin routing and request forwarding
- `src/models/load_balancer_test.go` – unit tests for core behavior

## How to run

```bash
go run ./src
```

Then open:

```text
http://localhost:8000
```

## Why this exists

This is intentionally a learning-focused project rather than a production-ready load balancer. It is designed to help you understand:

- reverse proxies
- round robin selection
- health checking
- concurrency concerns in Go
- how routing logic is separated from concrete implementations

## Important note

This code is useful for education, but it is not a production load balancer. It does not include:

- advanced retry policies
- weighted balancing
- circuit breakers
- production observability
- TLS termination and secure configuration
- real distributed health monitoring

## Next learning steps

Try improving the project in this order:

1. Make health checks configurable and testable
2. Add additional server strategies besides round robin
3. Add metrics and request logging
4. Refactor the server abstraction to support more realistic health checks
5. Add end-to-end tests with local mock upstream servers

## Useful commands

```bash
go test ./...
go build ./...
```
