package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andreyxaxa/http-load_balancer-rate_limiter/config"
	"github.com/andreyxaxa/http-load_balancer-rate_limiter/internal/repo/persistent"
	"github.com/andreyxaxa/http-load_balancer-rate_limiter/pkg/httpserver"
	"github.com/andreyxaxa/http-load_balancer-rate_limiter/pkg/loadbalancer"
	"github.com/andreyxaxa/http-load_balancer-rate_limiter/pkg/logger"
	"github.com/andreyxaxa/http-load_balancer-rate_limiter/pkg/postgres"
	"github.com/andreyxaxa/http-load_balancer-rate_limiter/pkg/ratelimiter"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Rate Limiter
	rl := ratelimiter.NewRateLimiter(persistent.New(pg))

	// Load Balancer
	lb := loadbalancer.NewLoadBalancer(cfg.LoadBalancer.Backends)
	log.Println(cfg.LoadBalancer.Backends)
	go loadbalancer.HealthCheck(lb)

	// HTTP Server
	handler := rl.Middleware(lb)
	log.Printf("Starting load balancer on :%s", cfg.HTTP.Port)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	go func() {
		if err := httpServer.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	s := <-interrupt
	log.Printf("app - Run - signal: %s", s.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Println("app - Run - httpServer.Shutdown: %w", err)
	}
}
