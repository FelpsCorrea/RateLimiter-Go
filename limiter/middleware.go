package limiter

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/FelpsCorrea/RateLimiter-Go/config"
)

func RateLimitMiddleware(storage LimiterStorage, config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ip := strings.Split(r.RemoteAddr, ":")[0]
			apiKey := r.Header.Get("API_KEY")

			var key string
			var limit int

			if apiKey != "" {
				key = "token:" + apiKey
				limit = config.RateLimitPerToken
			} else {
				key = "ip:" + ip
				limit = config.RateLimitPerIP
			}

			allowed, err := storage.Allow(ctx, key, limit, config.BlockDuration)
			if err != nil {
				log.Printf("Error checking rate limit: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				log.Printf("Request blocked for key: %s", key)
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			log.Printf("Request allowed for key: %s", key)
			next.ServeHTTP(w, r)
		})
	}
}
