package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	// Wraps the ListenAndServe function inside a Goroutine (a lightweight, concurrent function execution in Go),
	// so it doesn't block our graceful shutdown logic below
	go func() {
		err := s.ListenAndServe()

		// Add a new check: err != http.ErrServerClosed
		// When s.Shutdown() is called, ListenAndServe() will return http.ErrServerClosed,
		// which is an expected error during a graceful shutdown.
		// We want to log fatal errors only for unexpected issues during server startup.
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// It's often a good practice to buffer the channel (make(chan os.Signal, 1))
	// to prevent potential blocking if the signal handler isn't immediately ready to receive.
	// When a signal is sent to the program, signal.Notify writes the signal to the channel.
	// If the channel is unbuffered and no goroutine is actively reading from it,
	// the signal will be dropped. This can cause our program to miss termination signals
	// like os.Kill or os.Interrupt.
	sigChan := make(chan os.Signal, 1)

	// This line configure the Go runtime to forward operating system signals
	// os.Interrupt (usually sent by pressing Ctrl+C) to the sigChan channel.
	// os.Kill cannot be trapped, so we replaced it with syscall.SIGTERM to handle termination signals properly.
	// This allows your program to be notified when it's asked to shut down.
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// We've created a new go func() { ... }() block.
	// This starts a separate goroutine specifically for handling the received operating system signals.
	// Inside this goroutine, we still block and wait for a signal on sigChan (sig := <-sigChan)
	go func() {
		// This line blocks the execution of the main goroutine until a
		// signal is received on the sigChan channel.
		// When a os.Kill or os.Interrupt signal is received,
		// the value of that signal will be assigned to the sig variable.
		sig := <-sigChan

		// Once a signal is received, this line logs a message indicating which
		// signal was received and that the server is going to shut down gracefully.
		log.Printf("Received signal: %v, going for graceful shutdown.\n", sig)

		// To simulate some cleanup work (closing database connections, etc)
		time.Sleep(5 * time.Second)

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
		if err := s.Shutdown(tc); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		} else {
			log.Println("HTTP server gracefully shut down.")
		}
	}()

	// We've added an empty select {} at the end of the main function.
	// This will cause the main goroutine to block indefinitely.
	// This is important because if the main function exits before the
	// signal handling goroutine finishes, our program might terminate prematurely,
	// and the graceful shutdown might not complete.
	select {}
}
