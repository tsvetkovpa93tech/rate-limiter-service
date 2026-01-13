# Rate Limiter Service

A production-ready, high-performance rate limiting service built with Go. This service provides flexible rate limiting capabilities with multiple algorithms and storage backends, designed for distributed systems and microservices architectures.

## Table of Contents

- [What is Rate Limiting?](#what-is-rate-limiting)
- [Features](#features)
- [Architecture](#architecture)
- [Rate Limiting Algorithms](#rate-limiting-algorithms)
- [Storage Backends](#storage-backends)
- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Configuration](#configuration)
- [Load Testing](#load-testing)
- [Monitoring & Metrics](#monitoring--metrics)
- [Docker Deployment](#docker-deployment)
- [Development](#development)
- [Extending the Project](#extending-the-project)
- [License](#license)

## What is Rate Limiting?

Rate limiting is a technique used to control the rate of requests sent or received by a system. It helps:

- **Prevent abuse**: Protect APIs from being overwhelmed by malicious or buggy clients
- **Ensure fairness**: Distribute resources evenly among users
- **Maintain stability**: Prevent system overload and ensure consistent performance
- **Control costs**: Limit resource consumption in cloud environments

This service provides a standalone rate limiting solution that can be integrated into any application or used as a microservice in distributed architectures.

## Features

- ✅ **Multiple Algorithms**: Token Bucket and Sliding Window Log implementations
- ✅ **Flexible Storage**: In-memory (sync.Map) and Redis support
- ✅ **REST API**: Clean HTTP API with comprehensive error handling
- ✅ **Production Ready**: Graceful shutdown, structured logging, and comprehensive metrics
- ✅ **Observability**: Prometheus metrics and health checks
- ✅ **Configurable**: Environment variables and YAML configuration
- ✅ **Docker Support**: Ready-to-use Docker images and Docker Compose setup
- ✅ **Load Tested**: Performance benchmarks and comparison results
- ✅ **Enterprise Features**: Multi-tenancy, dynamic limits, webhooks
- ✅ **Optimized**: Object pooling, caching, connection pooling
- ✅ **Monitoring**: Grafana dashboards and Prometheus alerts

## Architecture

The project follows clean architecture principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer (Chi Router)                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Handlers   │  │  Middleware   │  │   Metrics     │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                    Service Layer                             │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         Rate Limiter Service                          │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
┌─────────────────────────────────────────────────────────────┐
│                    Core Layer                                │
│  ┌──────────────┐              ┌──────────────┐            │
│  │  Algorithms  │              │   Storage    │            │
│  │              │              │              │            │
│  │ Token Bucket │              │   Memory     │            │
│  │ Sliding Win. │              │    Redis    │            │
│  └──────────────┘              └──────────────┘            │
└─────────────────────────────────────────────────────────────┘
```

### Project Structure

```
rate-limiter-service/
├── cmd/
│   └── server/              # Application entry point
├── internal/
│   ├── handlers/           # HTTP handlers (limit, health, metrics)
│   ├── middleware/         # HTTP middleware (logging, recovery, CORS)
│   ├── service/            # Business logic layer
│   ├── metrics/           # Prometheus metrics
│   ├── services/          # Rate limiting algorithms
│   └── storage/          # Storage implementations
├── pkg/
│   └── config/            # Configuration management
├── api/
│   └── openapi.yaml       # OpenAPI/Swagger specification
├── loadtest/              # Load testing scripts (k6)
├── examples/              # Example clients and usage
├── configs/               # Configuration files
└── prometheus/            # Prometheus configuration
```

## Rate Limiting Algorithms

### Algorithm Comparison

| Feature | Token Bucket | Sliding Window Log |
|---------|--------------|-------------------|
| **Complexity** | Simple | More complex |
| **Memory Usage** | Low (2 integers per key) | Medium (array of timestamps) |
| **Burst Capacity** | ✅ Yes (up to bucket size) | ❌ No |
| **Precision** | Good | Excellent |
| **Performance** | Faster (~2ms avg) | Slightly slower (~3ms avg) |
| **Use Case** | General purpose | Strict requirements |
| **Best For** | APIs, general rate limiting | Payment APIs, strict quotas |

### Token Bucket

The Token Bucket algorithm maintains a bucket with a fixed number of tokens. Tokens are added to the bucket at a constant rate. When a request arrives, a token is consumed. If no tokens are available, the request is denied.

**Characteristics:**
- ✅ Allows burst traffic (up to bucket capacity)
- ✅ Smooth rate limiting
- ✅ Lower memory overhead (stores only 2 integers: tokens count and last refill time)
- ✅ Simpler implementation
- ✅ Best for: General purpose rate limiting

**Example:**
```json
{
  "key": "user:123",
  "algorithm": "token_bucket",
  "limit": 100,
  "window": "1m"
}
```
This allows 100 requests per minute with burst capacity.

### Sliding Window Log

The Sliding Window Log algorithm maintains a log of request timestamps within a time window. When a request arrives, old timestamps outside the window are removed, and the new timestamp is added. If the count exceeds the limit, the request is denied.

**Characteristics:**
- ✅ More precise rate limiting (exact count)
- ✅ No burst capacity (strict enforcement)
- ✅ Higher memory overhead (stores array of timestamps)
- ✅ More accurate for strict quotas
- ✅ Best for: Strict rate limiting requirements

**Example:**
```json
{
  "key": "user:123",
  "algorithm": "sliding_window",
  "limit": 100,
  "window": "1m"
}
```
This strictly enforces 100 requests per minute with no bursts.

### Performance Benchmarks

Benchmark results from Go benchmark tests (run with `go test -bench=. ./internal/services/`):

```
BenchmarkTokenBucket_Allow-8                   5000000    250 ns/op    48 B/op    1 allocs/op
BenchmarkSlidingWindow_Allow-8                 3000000    380 ns/op    64 B/op    2 allocs/op
BenchmarkTokenBucket_Allow_DifferentKeys-8     2000000    280 ns/op    56 B/op    1 allocs/op
BenchmarkSlidingWindow_Allow_DifferentKeys-8   1500000    420 ns/op    72 B/op    2 allocs/op
```

**Key Findings:**
- Token Bucket is ~35% faster than Sliding Window
- Token Bucket uses ~25% less memory
- Both algorithms scale well with different keys
- Performance difference is minimal for most use cases

## Storage Backends

### In-Memory (sync.Map)

Fast, lightweight storage using Go's `sync.Map`. Suitable for single-instance deployments.

**Pros:**
- ✅ Very fast (no network overhead)
- ✅ No external dependencies
- ✅ Simple deployment

**Cons:**
- ❌ Not suitable for distributed systems
- ❌ Data lost on restart

### Redis

Distributed storage using Redis. Suitable for multi-instance deployments and production environments.

**Pros:**
- ✅ Distributed rate limiting
- ✅ Persistent storage
- ✅ High performance
- ✅ Scalable

**Cons:**
- ❌ Requires Redis infrastructure
- ❌ Network latency

## Quick Start

### Prerequisites

- Go 1.21 or higher
- (Optional) Docker and Docker Compose for containerized deployment
- (Optional) Redis for distributed storage

### Local Development

1. **Clone the repository:**
```bash
git clone https://github.com/yourusername/rate-limiter-service.git
cd rate-limiter-service
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Run the server:**
```bash
go run ./cmd/server
```

4. **Test the API:**
```bash
curl -X POST http://localhost:8080/api/v1/limit-check \
  -H "Content-Type: application/json" \
  -d '{"key": "user:123", "limit": 5, "window": "1m"}'
```

### Docker Deployment

1. **Start all services (including Redis, Prometheus, Grafana):**
```bash
docker-compose up -d
```

2. **Check service health:**
```bash
curl http://localhost:8080/health
```

3. **Access services:**
- API: http://localhost:8080
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000 (admin/admin)

## API Documentation

### POST /api/v1/limit-check

Check if a request should be allowed based on rate limiting rules.

**Request:**
```bash
curl -X POST http://localhost:8080/api/v1/limit-check \
  -H "Content-Type: application/json" \
  -d '{
    "key": "user:123",
    "algorithm": "token_bucket",
    "limit": 100,
    "window": "1m"
  }'
```

**Response (200 OK - Allowed):**
```json
{
  "allowed": true,
  "reset_at": 1704067200
}
```

**Response (429 Too Many Requests - Denied):**
```json
{
  "allowed": false,
  "reset_at": 1704067200,
  "message": "Rate limit exceeded"
}
```

### GET /health

Health check endpoint.

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok",
  "service": "rate-limiter-service",
  "version": "1.0.0"
}
```

### GET /metrics

Prometheus metrics endpoint.

```bash
curl http://localhost:8080/metrics
```

## Configuration

### Environment Variables

Copy `env.example` to `.env` and configure:

```bash
# Server
RL_SERVER_PORT=8080
RL_SERVER_READ_TIMEOUT=15s
RL_SERVER_WRITE_TIMEOUT=15s
RL_SERVER_IDLE_TIMEOUT=60s

# Storage
RL_STORAGE_TYPE=memory  # or "redis"
RL_REDIS_ADDRESS=localhost:6379
RL_REDIS_DB=0
RL_REDIS_PASSWORD=

# Rate Limiter
RL_DEFAULT_ALGORITHM=token_bucket
RL_DEFAULT_LIMIT=100
RL_DEFAULT_WINDOW=1m

# CORS
RL_CORS_ALLOWED_ORIGINS=*
```

### Configuration File

Create `configs/config.yml`:

```yaml
server:
  port: 8080
  read_timeout: 15s
  write_timeout: 15s
  idle_timeout: 60s

storage:
  type: memory  # or "redis"
  redis_address: localhost:6379
  redis_db: 0
  redis_password: ""

limiter:
  default_algorithm: token_bucket
  default_limit: 100
  default_window: 1m

cors:
  allowed_origins:
    - "*"
```

## Load Testing

Load testing scripts are available in the `loadtest/` directory using k6.

### Running Load Tests

1. **Install k6:**
```bash
# macOS
brew install k6

# Linux
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6
```

2. **Run Token Bucket test:**
```bash
k6 run loadtest/k6-token-bucket.js
```

3. **Run Sliding Window test:**
```bash
k6 run loadtest/k6-sliding-window.js
```

### Performance Comparison

Based on load testing with 200 concurrent users and Go benchmarks:

| Algorithm | Avg Latency | P95 Latency | Throughput | Memory Usage | Allocations |
|-----------|-------------|-------------|------------|--------------|-------------|
| Token Bucket | ~2ms | ~5ms | ~5000 req/s | Low (48 B/op) | 1 alloc/op |
| Sliding Window | ~3ms | ~8ms | ~4500 req/s | Medium (64 B/op) | 2 allocs/op |

**Go Benchmark Results:**
```
BenchmarkTokenBucket_Allow-8                   5000000    250 ns/op    48 B/op    1 allocs/op
BenchmarkSlidingWindow_Allow-8                 3000000    380 ns/op    64 B/op    2 allocs/op
```

**Key Findings:**
- Token Bucket performs ~35% faster (250ns vs 380ns per operation)
- Token Bucket uses ~25% less memory (48B vs 64B per operation)
- Token Bucket has lower GC pressure (1 vs 2 allocations per operation)
- Sliding Window provides more precise rate limiting
- Both algorithms handle 200+ concurrent users efficiently
- Memory usage is acceptable for both algorithms

See [BENCHMARKS.md](BENCHMARKS.md) for detailed benchmark results and optimization analysis.

## Monitoring & Metrics

### Prometheus Metrics

The service exposes the following Prometheus metrics:

- `rate_limiter_total_requests` - Total requests (by method, endpoint, status)
- `rate_limiter_blocked_requests` - Blocked requests (by algorithm, key)
- `rate_limiter_allowed_requests_total` - Allowed requests (by algorithm)
- `rate_limiter_denied_requests_total` - Denied requests (by algorithm)
- `rate_limiter_check_errors_total` - Check errors (by algorithm)
- `rate_limiter_request_duration_seconds` - Request duration histogram

### Grafana Dashboards

The Docker Compose setup includes Grafana with pre-configured dashboards for:
- Request rate and latency
- Error rates
- Algorithm performance comparison
- Storage backend metrics

## Docker Deployment

### Build Docker Image

```bash
docker build -t rate-limiter-service:latest .
```

### Run with Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f rate-limiter

# Stop services
docker-compose down
```

### Production Considerations

1. **Use Redis** for distributed rate limiting
2. **Configure CORS** properly for your domain
3. **Set appropriate timeouts** based on your load
4. **Monitor metrics** using Prometheus/Grafana
5. **Use HTTPS** in production (add reverse proxy)

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Building

```bash
# Development build
go build -o bin/server ./cmd/server

# Production build (optimized)
go build -ldflags="-s -w" -o bin/server ./cmd/server
```

### Using Makefile

```bash
make build      # Build application
make run        # Run application
make test       # Run tests
make docker-up  # Start Docker services
```

## Enterprise Features

### Multi-Tenancy

Support for different rate limits per API key/tenant:

```go
// Configure tenant limits
tenantManager.SetConfig(&tenant.TenantConfig{
    APIKey:    "api-key-123",
    Algorithm: "token_bucket",
    Limit:     1000,
    Window:    time.Minute,
    Enabled:   true,
})
```

### Dynamic Limit Updates

Update rate limits without restarting the service:

```go
// Update tenant limits dynamically
tenantManager.UpdateConfig("api-key-123", 2000, time.Minute)
```

### Webhook Notifications

Get notified when rate limits are exceeded:

```yaml
webhook:
  enabled: true
  url: "https://your-service.com/webhooks/rate-limit"
  timeout: "5s"
```

Webhook payload:
```json
{
  "type": "blocked",
  "key": "user:123",
  "algorithm": "token_bucket",
  "limit": 100,
  "window": "1m",
  "timestamp": "2024-01-01T12:00:00Z",
  "message": "Rate limit exceeded"
}
```

## Performance Optimizations

### Object Pooling

The service uses `sync.Pool` for frequently allocated objects:
- JSON buffers
- HTTP request buffers
- Reduces GC pressure

### Caching

In-memory cache for frequent rate limit checks:
- Configurable TTL
- Automatic cleanup
- Reduces storage lookups

### Redis Connection Pooling

Optimized Redis connection management:
- Configurable pool size (default: 10)
- Minimum idle connections (default: 5)
- Connection timeouts and retries

## Extending the Project

### Adding a New Algorithm

1. Implement the `RateLimiter` interface in `internal/services/`:
```go
type MyAlgorithm struct {
    storage storage.Storage
    // ... fields
}

func (m *MyAlgorithm) Allow(ctx context.Context, key string) (bool, error) {
    // Implementation
}
```

2. Add to factory in `internal/factory.go`:
```go
case AlgorithmMyAlgorithm:
    return services.NewMyAlgorithm(...), nil
```

### Adding a New Storage Backend

1. Implement the `Storage` interface in `internal/storage/`:
```go
type MyStorage struct {
    // ... fields
}

func (m *MyStorage) Get(ctx context.Context, key string) (interface{}, error) {
    // Implementation
}
```

2. Add to storage factory in `internal/storage/storage_factory.go`

### Adding Middleware

Create new middleware in `internal/middleware/` and add to router in `cmd/server/main.go`.

## Example Client

See `examples/client.go` for a complete example of how to use the service:

```bash
go run ./examples/client.go
```

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## Support

For issues and questions:
- Open an issue on GitHub
- Check the documentation in `docs/`
- Review the OpenAPI spec in `api/openapi.yaml`
