package app

import (
	"context"
	"loadbalancer/config"
	"loadbalancer/package/httpserver"
	"loadbalancer/package/loadbalancer"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *config.Config) {
	// logger

	// Load Balancer
	lb := loadbalancer.NewLoadBalancer(cfg.LoadBalancer.Backends)
	go loadbalancer.HealthCheck(lb)

	// HTTP Server
	log.Printf("Starting load balancer on :8080")
	httpServer := httpserver.New(lb, httpserver.Port("8080")) // TODO: через конфиг
	go httpServer.Start()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM)

	s := <-interrupt
	log.Printf("app - Run - signal: %s", s.String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Println("app - Run - httpServer.Shutdown: %w", err)
	}
}
