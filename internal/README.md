# Rate Limiter Core

Ядро rate limiter сервиса с production-ready реализацией алгоритмов и хранилищ.

## Структура

```
internal/
├── limiter.go              # Интерфейс RateLimiter
├── factory.go              # Фабрика для создания лимитеров
├── services/               # Реализации алгоритмов
│   ├── token_bucket.go     # Token Bucket алгоритм
│   └── sliding_window.go   # Sliding Window Log алгоритм
└── storage/                # Реализации хранилищ
    ├── storage.go          # Интерфейс Storage
    ├── memory.go           # In-memory хранилище (sync.Map)
    └── redis.go            # Redis хранилище
```

## Использование

### Создание лимитера через фабрику

```go
import (
    "context"
    "time"
    "log/slog"
    
    "github.com/tsvetkovpa93tech/rate-limiter-service/internal"
    "github.com/tsvetkovpa93tech/rate-limiter-service/internal/storage"
)

// Создание хранилища
logger := slog.Default()
memStorage := storage.NewMemoryStorage(logger)

// Создание Token Bucket лимитера
limiter, err := internal.NewRateLimiter(internal.LimiterConfig{
    Algorithm: internal.AlgorithmTokenBucket,
    Limit:     100,
    Window:    time.Minute,
    Storage:   memStorage,
    Logger:    logger,
})

// Использование
ctx := context.Background()
allowed, err := limiter.Allow(ctx, "user:123")
```

### Алгоритмы

#### Token Bucket
- Плавное ограничение с возможностью "burst"
- Токены пополняются с постоянной скоростью
- Подходит для общего использования

#### Sliding Window Log
- Точное ограничение с отслеживанием временных меток
- Более точный контроль над количеством запросов
- Подходит для строгих требований

### Хранилища

#### Memory Storage
- Быстрое in-memory хранилище на основе sync.Map
- Подходит для single-instance развертываний
- Автоматическое удаление истекших записей

#### Redis Storage
- Распределенное хранилище для multi-instance развертываний
- Поддержка TTL для автоматического истечения
- Подходит для production окружений

## Особенности

- ✅ Поддержка `context.Context` для отмены операций
- ✅ Структурированное логирование через `log/slog`
- ✅ Правильная обработка ошибок
- ✅ Полное покрытие unit-тестами
- ✅ Production-ready код

## Тестирование

Запуск всех тестов:

```bash
go test ./internal/...
```

Запуск тестов с покрытием:

```bash
go test -cover ./internal/...
```

