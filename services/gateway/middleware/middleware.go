package middleware

import (
	"github.com/go-chi/cors"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func WithCORS(handler http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Use this to allow specific origin hosts, * for matches any origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Signature"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
	})
	return c.Handler(handler)
}

func ChainCombine(middlewares ...Middleware) Middleware {
	return func(handler http.Handler) http.Handler {
		for _, m := range middlewares {
			handler = m(handler)
		}
		return handler
	}
}
