package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ww24/lawn"
)

const (
	defaultMaxAge   = "24h"
	defaultPort     = "8000"
	shutdownTimeout = 15 * time.Second
)

var (
	username = os.Getenv("GITHUB_USERNAME")
	maxAge   = os.Getenv("MAX_AGE")
	port     = os.Getenv("PORT")
)

func main() {
	if username == "" {
		log.Fatalln("GITHUB_USERNAME env is required")
	}

	if maxAge == "" {
		maxAge = defaultMaxAge
	}
	if port == "" {
		port = defaultPort
	}
	d, err := time.ParseDuration(maxAge)
	if err != nil {
		log.Fatalln("Failed to parse maxAge:", err)
	}

	cli := lawn.NewClient()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: newHandler(cli, d),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalln("ListenAndServe error:", err)
		}
	}()
	log.Println("Listen at", port)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("Shutdown error:", err)
	}
}
