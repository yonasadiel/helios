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
	return Handle(wrapped)
}

// CreateCORSMiddleware add Access-Control-Allow-Origin headers
// to the response.
func CreateCORSMiddleware(allowedOrigins []string) Middleware {
	var allAllowed bool = false
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == "*" {
			allAllowed = true
		}
	}

	return func(f HTTPHandler) HTTPHandler {
		return func(req Request) {
			var origin string = req.GetHeader("Origin")
			var allowed bool = false

			if allAllowed {
				allowed = true
			} else {
				for _, allowedOrigin := range allowedOrigins {
					if allowedOrigin == origin {
						allowed = true
					}
				}
			}

			if allowed {
				req.SetHeader("Access-Control-Allow-Origin", origin)
			} else {
				req.SetHeader("Access-Control-Allow-Origin", "")
			}

			f(req)
		}
	}
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
