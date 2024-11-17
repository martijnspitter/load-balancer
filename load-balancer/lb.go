package loadbalancer

import (
	"load-balancer/backend"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type LoadBalancer struct {
	backends []*backend.Backend
	mutex    sync.Mutex
}

func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{}
}

func (lb *LoadBalancer) AddBackend(urlStr string) error {
	url, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	backend := backend.InitBackend(url)

	lb.mutex.Lock()
	lb.backends = append(lb.backends, backend)
	lb.mutex.Unlock()

	return nil
}

func (lb *LoadBalancer) NextBackend() *backend.Backend {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	return lb.roundRobin()
}

func (lb *LoadBalancer) roundRobin() *backend.Backend {
	if len(lb.backends) == 0 {
		return nil
	}

	backend := lb.backends[0]
	lb.backends = append(lb.backends[1:], backend)

	if !backend.Alive {
		return nil
	}

	return backend
}

func (lb *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend := lb.NextBackend()
	if backend == nil {
		http.Error(w, "No backends available", http.StatusServiceUnavailable)
		return
	}

	log.Println("Received request from ", r.RemoteAddr)
	log.Println(r.Method, r.URL.Path, r.Proto)
	log.Println("Host: ", r.Host)
	log.Println("User-Agent: ", r.UserAgent())
	log.Println("Accept: ", r.Header.Get("Accept"))

	backend.ReverseProxy.ServeHTTP(w, r)
}
