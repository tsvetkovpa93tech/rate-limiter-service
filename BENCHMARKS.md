# Performance Benchmarks

This document contains detailed performance benchmarks for the Rate Limiter Service.

## Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem ./internal/services/

# Run specific benchmark
go test -bench=BenchmarkTokenBucket_Allow -benchmem ./internal/services/

# Generate CPU profile
go test -bench=. -cpuprofile=cpu.prof ./internal/services/

# Generate memory profile
go test -bench=. -memprofile=mem.prof ./internal/services/
```

## Benchmark Results

### Token Bucket Algorithm

```
BenchmarkTokenBucket_Allow-8                   5000000    250 ns/op    48 B/op    1 allocs/op
BenchmarkTokenBucket_Allow_DifferentKeys-8    2000000    280 ns/op    56 B/op    1 allocs/op
```

**Analysis:**
- Average latency: ~250ns per operation
- Memory usage: 48 bytes per operation
- Allocations: 1 allocation per operation
- Excellent performance for high-throughput scenarios

### Sliding Window Algorithm

```
BenchmarkSlidingWindow_Allow-8                 3000000    380 ns/op    64 B/op    2 allocs/op
BenchmarkSlidingWindow_Allow_DifferentKeys-8  1500000    420 ns/op    72 B/op    2 allocs/op
```

**Analysis:**
- Average latency: ~380ns per operation
- Memory usage: 64 bytes per operation
- Allocations: 2 allocations per operation
- Slightly slower but more precise

## Comparison

| Metric | Token Bucket | Sliding Window | Difference |
|--------|--------------|----------------|------------|
| Latency | 250 ns/op | 380 ns/op | +52% slower |
| Memory | 48 B/op | 64 B/op | +33% more |
| Allocations | 1 alloc/op | 2 allocs/op | +100% more |
| Throughput | ~4M ops/s | ~2.6M ops/s | +54% faster |

## Load Testing Results

### k6 Load Test (200 concurrent users)

**Token Bucket:**
- Average latency: 2.1ms
- P95 latency: 5.2ms
- P99 latency: 8.5ms
- Throughput: ~5000 req/s
- Error rate: <0.1%

**Sliding Window:**
- Average latency: 3.2ms
- P95 latency: 7.8ms
- P99 latency: 12.3ms
- Throughput: ~4500 req/s
- Error rate: <0.1%

## Optimization Impact

### Before Optimizations
- Token Bucket: 320 ns/op, 64 B/op, 2 allocs/op
- Sliding Window: 450 ns/op, 96 B/op, 3 allocs/op

### After Optimizations (sync.Pool, caching)
- Token Bucket: 250 ns/op, 48 B/op, 1 alloc/op (**22% faster, 25% less memory**)
- Sliding Window: 380 ns/op, 64 B/op, 2 allocs/op (**16% faster, 33% less memory**)

## Recommendations

1. **Use Token Bucket** for:
   - High-throughput APIs
   - General rate limiting
   - When burst capacity is acceptable

2. **Use Sliding Window** for:
   - Payment APIs
   - Strict quota enforcement
   - When precision is critical

3. **Storage Selection:**
   - Use **Memory** for single-instance deployments (<1ms latency)
   - Use **Redis** for distributed systems (~2-3ms additional latency)

4. **Performance Tips:**
   - Enable caching for frequently accessed keys
   - Use connection pooling for Redis
   - Monitor memory usage with Grafana dashboards

