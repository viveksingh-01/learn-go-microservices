package main

import (
	"log"
	"net/http"
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

		// This specifies the network address the server should listen on.
		// ":9090" means it will listen on all available network interfaces on port 9090.
		Addr: ":9090",

		// Tells the server to use our custom router to determine which handler
		// should handle each incoming request.
		Handler: sm,

		// This sets the maximum amount of time an idle (keep-alive) connection will
		// remain open before the server closes it.
		// In this case, it's set to 120 seconds (2 minutes).
		// This helps prevent resource exhaustion from inactive connections
		IdleTimeout: 120 * time.Second,

		// This sets the maximum duration for reading the entire request, including the body.
		// If the client takes longer than 1 second to send the full request,
		// the server will time out and close the connection.
		// This helps protect against slow or malicious clients.
		ReadTimeout: 1 * time.Second,

		// This sets the maximum duration for writing the response back to the client.
		// If the server takes longer than 1 second to send the complete response,
		// it will time out and close the connection.
		// This also helps prevent issues with slow clients.
		WriteTimeout: 1 * time.Second,
	}

	// Starts the HTTP server using the configuration we defined in our custom server 's'.
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
