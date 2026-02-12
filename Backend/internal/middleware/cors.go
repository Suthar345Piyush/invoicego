// CORS middleware code here
// http methods , headers and important things for auth

package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

// this function returns http handlers

func CORS(allowedOrigins []string) func(http.Handler) http.Handler {

	return cors.Handler(cors.Options{

		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

}
