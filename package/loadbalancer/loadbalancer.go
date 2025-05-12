package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"
)

// --- Backend ---

// TODO: сделать интерфейс

type Backend struct {
	URL          *url.URL
	Alive        bool
	ReverseProxy *httputil.ReverseProxy
}

func (b *Backend) CheckAlive() bool {
	req, err := http.NewRequest("GET", b.URL.String(), nil)
	if err != nil {
		return false
	}

	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 500 {
		return false
	}

	return true
}

// --- LoadBalancer ---

// TODO: сделать интерфейс

type LoadBalancer struct {
	backends []*Backend
	index    int32
}

func (lb *LoadBalancer) getNextBackend() *Backend {
	length := len(lb.backends)

	for i := 0; i < length; i++ {
		idx := atomic.AddInt32(&lb.index, 1) % int32(length)
		if lb.backends[idx].Alive {
			return lb.backends[idx]
		}
	}

	return nil
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO: lb.l.Info(...). Ну логгер заинжектить в балансер и туда писать всё.
	log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

	backend := lb.getNextBackend()
	if backend == nil {
		log.Println("No backend available")
		http.Error(w, "No backend available", http.StatusServiceUnavailable)
		return
	}

	// TODO: lb.l.Info(...)
	log.Printf("Forwarding to backend : %s", backend.URL)
	backend.ReverseProxy.ServeHTTP(w, r)
}

// TODO: также будет принимать интерфейс логгера.
func NewLoadBalancer(urls []string) *LoadBalancer {
	var backends []*Backend

	for _, u := range urls {
		parsedURL, err := url.Parse(u)
		if err != nil {
			log.Fatalf("Invalid backend URL: %s", u)
		}

		proxy := httputil.NewSingleHostReverseProxy(parsedURL)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error for %s: %v", parsedURL, err)
			http.Error(w, "Backend error", http.StatusBadGateway)
		}

		backends = append(backends, &Backend{
			URL:          parsedURL,
			Alive:        true,
			ReverseProxy: proxy,
		})
	}

	return &LoadBalancer{backends: backends}
}

// ---

func HealthCheck(lb *LoadBalancer) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		for _, b := range lb.backends {
			b.Alive = b.CheckAlive()
		}
	}
}
