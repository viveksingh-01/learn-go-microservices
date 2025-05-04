package main

import (
	"fmt"
	"io"
	"net/http"
)

// This is the main function, the entry point of our Go program. Execution always begins here.
func main() {
	// http.HandleFunc registers a handler function for a specific URL path on DefaultServeMux (an HTTP handler).

	// In this case, the path is "/", which represents the root of our website
	// The second argument is an anonymous function (a function without a name) that will be
	// executed whenever a request comes in for the "/" path.
	// It takes takes two arguments:
	// 'w http.ResponseWriter': This is an interface that allows us to write the HTTP response
	// back to the client (the browser or whatever made the request).
	// 'r *http.Request': This is a pointer to a Request struct, which contains all the information
	// about the incoming HTTP request (like headers, body, method, etc.).
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Reads the entire body of the incoming HTTP request.
		// r.Body is an io.ReadCloser, which represents the stream of data in the request body.
		// io.ReadAll reads all the data from the io.ReadCloser until it encounters an error or the end of the stream.
		// The reason we can pass r.Body (an io.ReadCloser) to io.ReadAll (which expects an io.Reader) is
		// because io.ReadCloser implements the io.Reader interface.
		d, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Oops..", http.StatusBadRequest)
			return
		}

		// Writes the HTTP response back to the client.
		// fmt.Fprintf formats and writes a string to a writer. In this case, the writer is w (http.ResponseWriter).
		// "Hello %s" is the format string. The %s is a placeholder for a string value.
		// d (the byte slice we read from the request body) is used as the value to replace %s.
		// Go will automatically convert the byte slice to a string for the output.
		fmt.Fprintf(w, "Hello %s", d)
	})

	// Starts the HTTP server and makes it listen for incoming connections on the network address :9090.
	// The second argument is a handler to use for incoming connections that don't match any specific HandleFunc registrations.
	// Passing nil here tells the server to use the default ServeMux, which already contains the handler we registered for "/".
	http.ListenAndServe(":9090", nil)
}

// ServeMux:
// - An HTTP request multiplexer.
// - It matches the URL of each incoming request against a list of registered patterns and calls
//   the handler for the pattern that most closely matches the URL.

// Handler:
// - An Interface which responds to HTTP requests.
// - Any struct with method ServeHTTP(ResponseWriter, *Request) will implement Handler interface.

// NOTE:

// The reason we can pass w (an http.ResponseWriter) to fmt.Fprintf (which expects an io.Writer) is because
// http.ResponseWriter implements the io.Writer interface.
//
// Let's break down what that means:
//
// Interfaces in Go:
// In Go, an interface is a type that defines a set of method signatures.
// Any type that provides implementations for all the methods defined in an interface is said to implement that interface.
// This implementation is implicit; we don't need to explicitly declare that a type implements an interface.

// The 'io' package in Go provides basic interfaces for I/O primitives. The io.Writer interface is defined as follows:
// type Writer interface {
//     Write(p []byte) (n int, err error)
// }
// Any type that has a Write method that takes a byte slice ([]byte) and returns the number of bytes written (int) and
// an optional error (error) satisfies the io.Writer interface.

// http.ResponseWriter Interface
// The http.ResponseWriter interface, from the net/http package, is used by an HTTP handler to construct the HTTP response.
//
// While it has other methods specific to HTTP responses (like WriteHeader, Header), it also has a Write method that
// matches the signature of the io.Writer interface:
// type ResponseWriter interface {
//     Header() Header
//     Write([]byte) (int, error)
//     WriteHeader(statusCode int)
// }

// How fmt.Fprintf Uses io.Writer:
//
// The fmt.Fprintf function has the following signature:
// func Fprintf(w io.Writer, format string, a ...any) (n int, err error)
//
// As we can see, the first argument w is of type io.Writer.
// fmt.Fprintf doesn't need to know the specific underlying type of w.
// All it cares about is that w has a Write method that allows it to write a sequence of bytes.
// fmt.Fprintf formats the given string according to the format specifier and then uses the Write method of the provided
// io.Writer to send the resulting bytes somewhere (in this case, back to the HTTP client).
