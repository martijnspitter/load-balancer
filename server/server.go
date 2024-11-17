package server

import (
	"context"
	loadbalancer "load-balancer/load-balancer"
	"log"
	"net/http"
	"time"

	"os"
	"os/signal"
)

func NewServer(lb *loadbalancer.LoadBalancer) *http.Server {
	return &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      lb,
	}
}

func StartServer(lb *loadbalancer.LoadBalancer) {
	s := NewServer(lb)

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := s.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Starting server on %s", s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
