package interfaces

import (
	"context"
)

// RateLimiter описывает контракт для алгоритмов ограничения скорости.
type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}
