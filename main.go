package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/viveksingh-01/learn-go-microservices/handlers"
)

func main() {

	// Creates an instance of the Hello handler.
	hh := &handlers.Hello{}

	// Creates a new instance of a ServeMux
	sm := http.NewServeMux()

	// Registers the Hello handler to handle all incoming requests to the root path (/).
	sm.Handle("/", hh)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// Wrapped the ListenAndServe function inside a Goroutine (a lightweight, concurrent function execution in Go),
	// so it doesn't block our graceful shutdown logic below
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// This creates a new channel named sigChan that can receive os.Signal values.
	// Channels are a way for goroutines to communicate safely.
	sigChan := make(chan os.Signal)

	// These lines configure the Go runtime to forward operating system signals
	// os.Kill (usually sent by the kill command) and
	// os.Interrupt (usually sent by pressing Ctrl+C) to the sigChan channel.
	// This allows your program to be notified when it's asked to shut down.
	signal.Notify(sigChan, os.Kill)
	signal.Notify(sigChan, os.Interrupt)

	// This line blocks the execution of the main goroutine until a
	// signal is received on the sigChan channel.
	// When a os.Kill or os.Interrupt signal is received,
	// the value of that signal will be assigned to the sig variable.
	sig := <-sigChan

	// Once a signal is received, this line logs a message indicating which
	// signal was received and that the server is going to shut down gracefully.
	log.Printf("Received signal: %v, going for graceful shutdown.\n", sig)

	// This creates a new context.Context with a timeout of 30 seconds.
	// context.Background() creates an empty root context.
	// context.WithTimeout() derives a new context from the parent context (context.Background())
	// that will be automatically canceled after the specified duration (30 seconds).
	// It returns the new context (tc) and a cancel function.
	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// 'defer' schedules the cancel function to be called when the main function exits.
	// This is important to release resources associated with the context,
	// even if the shutdown completes successfully before the timeout.
	defer cancel()

	// This is the crucial part of the graceful shutdown.
	// s.Shutdown(tc) attempts to gracefully shut down the HTTP server.
	// It stops accepting new connections and tries to close all idle connections.
	// It then waits for all active requests to complete (up to the timeout specified in the tc context).
	// If the timeout is reached before all requests complete,
	// the server will forcibly close any remaining active connections.
	s.Shutdown(tc)
}
