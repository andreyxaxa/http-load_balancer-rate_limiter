package main

import (
	"loadbalancer/config"
	"loadbalancer/internal/app"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./load_balancer <config file>")
	}

	configPath := os.Args[1]
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app.Run(cfg)
}
