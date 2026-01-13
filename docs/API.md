# API Documentation

## Endpoints

### POST /api/v1/limit-check

Проверяет, должен ли запрос быть разрешен на основе правил rate limiting.

**Request Body:**
```json
{
  "key": "user:123",
  "algorithm": "token_bucket",  // Optional: "token_bucket" or "sliding_window"
  "limit": 100,                 // Optional: количество запросов
  "window": "1m"                // Optional: временное окно (e.g., "1m", "30s")
}
```

**Response (200 OK):**
```json
{
  "allowed": true,
  "reset_at": 1704067200
}
```

**Response (429 Too Many Requests):**
```json
{
  "allowed": false,
  "reset_at": 1704067200,
  "message": "Rate limit exceeded"
}
```

### GET /health

Health check endpoint.

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

**Response:**
Prometheus metrics format

## Prometheus Metrics

- `rate_limiter_total_requests` - Общее количество запросов (по method, endpoint, status)
- `rate_limiter_blocked_requests` - Количество заблокированных запросов (по algorithm, key)
- `rate_limiter_allowed_requests_total` - Количество разрешенных запросов (по algorithm)
- `rate_limiter_denied_requests_total` - Количество отклоненных запросов (по algorithm)
- `rate_limiter_check_errors_total` - Количество ошибок проверки лимита (по algorithm)
- `rate_limiter_request_duration_seconds` - Длительность запросов (по method, endpoint, status)

## Middleware

- **Request Logging**: Логирует все HTTP запросы с деталями
- **Recovery**: Восстанавливается от паник
- **CORS**: Поддержка CORS headers
- **Request ID**: Добавляет уникальный ID к каждому запросу
- **Real IP**: Определяет реальный IP клиента

