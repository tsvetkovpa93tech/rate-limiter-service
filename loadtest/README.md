# Load Testing

This directory contains load testing scripts for the Rate Limiter Service.

## Prerequisites

Install k6:
```bash
# macOS
brew install k6

# Linux
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# Windows
choco install k6
```

## Running Tests

### Token Bucket Algorithm

```bash
k6 run loadtest/k6-token-bucket.js
```

### Sliding Window Algorithm

```bash
k6 run loadtest/k6-sliding-window.js
```

### With Custom Base URL

```bash
BASE_URL=http://localhost:8080 k6 run loadtest/k6-token-bucket.js
```

## Generating Reports

### HTML Report

```bash
k6 run --out json=results.json loadtest/k6-token-bucket.js
# Then use k6-reporter or other tools to generate HTML
```

### Cloud Results

```bash
k6 cloud loadtest/k6-token-bucket.js
```

## Test Scenarios

Both tests use the same load pattern:
- Ramp up to 50 users over 30 seconds
- Stay at 50 users for 1 minute
- Ramp up to 100 users over 30 seconds
- Stay at 100 users for 1 minute
- Ramp up to 200 users over 30 seconds
- Stay at 200 users for 1 minute
- Ramp down to 0 users over 30 seconds

## Metrics

The tests measure:
- Request duration (latency)
- Error rate
- Throughput (requests per second)
- Response status codes

## Expected Results

- **Token Bucket**: Better for burst traffic, allows short bursts above the limit
- **Sliding Window**: More precise rate limiting, stricter enforcement

