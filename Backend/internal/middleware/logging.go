// logging middleware code here

package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// function for header writing

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// main logging function , taking a handler and returns an handler (the main middleware part) here

// using pointer to modify the status code

func Logging(next http.Handler) http.Handler {

	// converting function into handler

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// storing time of request start

		start := time.Now()

		// creating a custom response writer
		// wrapping the struct

		wrapped := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		// at the end logging all(method , url , status code , time duration , remote address - client ip)

		log.Printf(
			"%s %s %d %s %s",
			r.Method,
			r.RequestURI,
			wrapped.status,
			time.Since(start),
			r.RemoteAddr,
		)
	})
}
