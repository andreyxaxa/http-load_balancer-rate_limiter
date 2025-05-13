package loadbalancer

import (
	"net/http"
)

// --- Backend ---

type Backend interface {
	CheckAlive() bool
}

// --- LoadBalancer ---

type LoadBalancer interface {
	getNextBackend() *Backend
	ServeHTTP(http.ResponseWriter, *http.Request)
}
