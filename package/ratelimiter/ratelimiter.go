package ratelimiter

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	_defaultCapacity = 5
	_defaultFillRate = 0.5
)

// --- ClientConfig ---

type ClientConfig struct {
	Capacity int
	FillRate float64
}

// --- TokenBucket ---

type TokenBucket struct {
	mu         sync.Mutex
	tokens     float64
	lastRefill time.Time
	config     ClientConfig
}

func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()

	b.tokens += elapsed * b.config.FillRate

	if b.tokens > float64(b.config.Capacity) {
		b.tokens = float64(b.config.Capacity)
	}

	b.lastRefill = now

	if b.tokens >= 1 {
		b.tokens -= 1
		return true
	}

	return false
}

// --- RateLimiter ---

// Интерфейс в месте использования
type RateLimitRepo interface {
	GetConfig(ctx context.Context, clientIP string) (*ClientConfig, error)
	SetConfig(ctx context.Context, clientIP string, cfg *ClientConfig) error
}

type RateLimiter struct {
	repo    RateLimitRepo // Dependency Injection
	buckets map[string]*TokenBucket
	mu      sync.RWMutex
}

func NewRateLimiter(repo RateLimitRepo) *RateLimiter {
	return &RateLimiter{
		repo:    repo,
		buckets: make(map[string]*TokenBucket),
	}
}

func (l *RateLimiter) SetLimit(ctx context.Context, key string, cfg ClientConfig) error {
	err := l.repo.SetConfig(ctx, key, &cfg)
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if bucket, exists := l.buckets[key]; exists {
		bucket.mu.Lock()
		bucket.config = cfg
		bucket.tokens = float64(cfg.Capacity)
		bucket.lastRefill = time.Now()
		bucket.mu.Unlock()
	}

	return nil
}

func (l *RateLimiter) getBucketForKey(ctx context.Context, key string) (*TokenBucket, error) {
	l.mu.RLock()
	bucket, exists := l.buckets[key]
	l.mu.RUnlock()

	if exists {
		return bucket, nil
	}

	cfg, err := l.repo.GetConfig(ctx, key)
	if err != nil {
		// Если не нашли юзера - создаем для него дефолтный конфиг
		if errors.Is(err, sql.ErrNoRows) {
			cfg = &ClientConfig{
				Capacity: _defaultCapacity,
				FillRate: _defaultFillRate,
			}
			if err := l.repo.SetConfig(ctx, key, cfg); err != nil {
				log.Println("RateLimiter - getBucketForKey - repo.SetConfig: %w", err)
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	bucket = &TokenBucket{
		tokens:     float64(cfg.Capacity),
		lastRefill: time.Now(),
		config:     *cfg,
	}

	l.mu.Lock()
	l.buckets[key] = bucket
	l.mu.Unlock()

	return bucket, nil
}

func (l *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Unable to determine client IP", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		bucket, err := l.getBucketForKey(ctx, clientIP)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Rate limiter config error", http.StatusInternalServerError)
			return
		}

		if !bucket.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
