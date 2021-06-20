package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/AkashGit21/hostelites/api"
)

type Server struct {
	*http.Server
}

func NewServer() (*http.Server, error) {
	log.Println("Configuring Server...")

	api, err := api.New()
	if err != nil {
		return nil, err
	}

	addr := ":8080"

	srv := http.Server{
		Addr:    addr,
		Handler: api,
		// Read and Write will Timeout after 2s if nothing goes correctly
		ReadTimeout:  time.Duration(2 * time.Second),
		WriteTimeout: time.Duration(2 * time.Second),
		// Connection will not stay Idle for more than 30 minutes
		IdleTimeout:    time.Duration(30 * time.Minute),
		MaxHeaderBytes: 1 << 20,
	}

	return &srv, nil
}

func StartServer(srv *http.Server) {
	log.Println("Starting Server...")

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
			return
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
		return
	}

	<-idleConnsClosed
}

func main() {

	srv, err := NewServer()
	if err != nil {
		log.Panic(err)
	}

	StartServer(srv)
}
