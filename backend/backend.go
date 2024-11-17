package backend

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

type Backend struct {
	Url          *url.URL
	ReverseProxy *httputil.ReverseProxy
	Alive        bool
	stopCheck    chan struct{} // Channel to signal stopping health checks
	mu           sync.Mutex    // Protect Alive field from concurrent access
}

func newBackend(url *url.URL) *Backend {
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Customizing the director to handle the request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Add("X-Proxy-By", "Go-LoadBalancer")
	}

	return &Backend{
		Url:          url,
		ReverseProxy: proxy,
		Alive:        true,
		stopCheck:    make(chan struct{}),
	}
}

func (b *Backend) SetAlive(alive bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.Alive
}

func (b *Backend) healthCheck() {
	resp, err := http.Get(b.Url.String() + "/health")

	if err != nil {
		log.Println("Backend is not alive: ", b.Url.String(), err)
		b.SetAlive(false)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Backend is not alive: ", b.Url.String())
		b.SetAlive(false)
		return
	}
	b.SetAlive(true)
	log.Println("Backend is alive: ", b.Url.String())
}

func (b *Backend) StartHealthCheck() {
	ticker := time.Tick(5 * time.Second)

	b.healthCheck()

	go func() {
		for range ticker {
			b.healthCheck()
		}
	}()
}

func (b *Backend) StopHealthCheck() {
	close(b.stopCheck)
}

func InitBackend(url *url.URL) *Backend {
	b := newBackend(url)
	b.StartHealthCheck()

	return b
}
