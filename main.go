package main

import (
	"log"
	"net/http"

	"github.com/FelpsCorrea/RateLimiter-Go/config"
	"github.com/FelpsCorrea/RateLimiter-Go/handlers"
	"github.com/FelpsCorrea/RateLimiter-Go/limiter"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()

	client := limiter.NewRedisClient(cfg.RedisAddr)
	defer client.Close()

	r := chi.NewRouter()
	r.Use(limiter.RateLimitMiddleware(client, cfg))

	r.Get("/", handlers.HelloHandler)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
