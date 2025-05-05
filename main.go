package main

import (
	"net/http"

	"github.com/viveksingh-01/learn-go-microservices/handlers"
)

// This is the main function, the entry point of our Go program. Execution always begins here.
func main() {

	// Creates an instance of the Hello handler.
	// &handlers.Hello{} creates a pointer (&) to a new instance of the Hello struct.
	// This pointer 'hh' will hold the memory address of the Hello struct object.
	hh := &handlers.Hello{}

	// Creates a new instance of a ServeMux
	sm := http.NewServeMux()

	// sm.Handle("/", hh) registers the HTTP handler hh (our Hello handler) to
	// handle requests that match the pattern '/'
	// The second argument is the handler itself.
	// Since the Hello struct has a ServeHTTP method, its pointer hh satisfies the http.Handler interface.
	sm.Handle("/", hh)

	// The second argument sm is the ServeMux, which tells the HTTP server to use our
	// custom router sm to handle incoming requests and dispatch them to the appropriate handlers.
	http.ListenAndServe(":9090", sm)
}

// ServeMux:
// - An HTTP request multiplexer.
// - It matches the URL of each incoming request against a list of registered patterns and calls
//   the handler for the pattern that most closely matches the URL.

// Handler:
// - An Interface which responds to HTTP requests.
// - Any struct with method ServeHTTP(ResponseWriter, *Request) will implement Handler interface.
