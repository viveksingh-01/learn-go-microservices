package handlers

import (
	"fmt"
	"io"
	"net/http"
)

// the Hello struct is an empty struct.
// Its purpose here is primarily to satisfy the interface requirement for an HTTP handler.
type Hello struct{}

// ServeHTTP is the crucial method that makes the Hello type an HTTP handler.
// It's part of the http.Handler interface from the net/http package.
// Any type that implements this ServeHTTP(ResponseWriter, *Request) method can
// be used to handle HTTP requests.
func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Oops..", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Hello %s", d)
}
