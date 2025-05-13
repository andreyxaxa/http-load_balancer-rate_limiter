package persistent

import (
	"context"
	"fmt"
	"loadbalancer/package/postgres"
	"loadbalancer/package/ratelimiter"

	"github.com/Masterminds/squirrel"
)

type RateLimitRepo struct {
	*postgres.Postgres
}

func New(pg *postgres.Postgres) *RateLimitRepo {
	return &RateLimitRepo{pg}
}

func (r *RateLimitRepo) GetConfig(ctx context.Context, clientIP string) (*ratelimiter.ClientConfig, error) {
	sql, args, err := r.Builder.
		Select("capacity", "fill_rate").
		From("rate_limits").
		Where(squirrel.Eq{"client_ip": clientIP}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("RateLimitRepo - GetConfig - r.Builder: %w", err)
	}

	row := r.Pool.QueryRow(ctx, sql, args...)

	var cfg ratelimiter.ClientConfig
	if err := row.Scan(
		&cfg.Capacity,
		&cfg.FillRate,
	); err != nil {
		return nil, fmt.Errorf("RateLimitRepo - GetConfig - row.Scan: %w", err)
	}

	return &cfg, nil
}

func (r *RateLimitRepo) SetConfig(ctx context.Context, clientIP string, cfg *ratelimiter.ClientConfig) error {
	sql, args, err := r.Builder.
		Insert("rate_limits").
		Columns("client_ip, capacity, fill_rate").
		Values(clientIP, cfg.Capacity, cfg.FillRate).
		Suffix("ON CONFLICT (client_ip) DO UPDATE SET capacity = EXCLUDED.capacity, fill_rate = EXCLUDED.fill_rate").
		ToSql()
	if err != nil {
		return fmt.Errorf("RateLimitRepo - SetConfig - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("RateLimitRepo - SetConfig - r.Pool.Exec: %w", err)
	}

	return nil
}
