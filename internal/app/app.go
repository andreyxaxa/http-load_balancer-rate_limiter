package app

import (
	"context"
	"loadbalancer/config"
	"loadbalancer/package/httpserver"
	"loadbalancer/package/loadbalancer"
	"loadbalancer/package/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Load Balancer
	lb := loadbalancer.NewLoadBalancer(cfg.LoadBalancer.Backends, l)
	log.Println(cfg.LoadBalancer.Backends)
	go loadbalancer.HealthCheck(lb)

	// HTTP Server
	log.Printf("Starting load balancer on :%s", cfg.HTTP.Port)
	httpServer := httpserver.New(lb, httpserver.Port(cfg.HTTP.Port))
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

	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Println("app - Run - httpServer.Shutdown: %w", err)
	}
}
