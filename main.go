package main

import (
	loadbalancer "load-balancer/load-balancer"
	"load-balancer/server"
	"log"
)

func main() {
	lb := loadbalancer.NewLoadBalancer()

	backends := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}

	for _, backend := range backends {
		if err := lb.AddBackend(backend); err != nil {
			log.Fatalf("Error adding backend: %v", err)
		}
	}

	server.StartServer(lb)
}
