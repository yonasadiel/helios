package helios

import (
	"net/http"
)

// Middleware is a function that receive an HTTPHandler
// and return new HTTPHandler
type Middleware func(HTTPHandler) HTTPHandler

// WithMiddleware wrapped the http request to Request object,
// pass it to middleware, start from first middleware in m to the last.
func WithMiddleware(f HTTPHandler, m []Middleware) func(http.ResponseWriter, *http.Request) {
	wrapped := makeMiddleware(f, m)
	return handle(wrapped)
}

// handle the http request using the HTTPHandler, without middleware
func handle(f HTTPHandler) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req HTTPRequest = NewHTTPRequest(w, r)
		f(&req)
	})
}

// makeMiddleware chains multiple middleware into new one
func makeMiddleware(f HTTPHandler, m []Middleware) HTTPHandler {
	wrapped := f
	for i := range m {
		// start from the last element
		wrapped = m[len(m)-i-1](wrapped)
	}
	return wrapped
}
