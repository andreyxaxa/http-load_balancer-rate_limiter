package repo

import (
	"context"
	"loadbalancer/package/ratelimiter"
)

type RateLimitRepo interface {
	GetConfig(ctx context.Context, clientID string) (*ratelimiter.ClientConfig, error)
	SetConfig(ctx context.Context, clientID string, cfg *ratelimiter.ClientConfig) error
	DeleteConfig(ctx context.Context, clientID string) error
}
