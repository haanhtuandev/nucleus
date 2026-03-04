package api

import (
	"boilerplate/internal/auth"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type ctxKey int

const userIDKey ctxKey = iota

type RateLimit struct {
	Mu           sync.RWMutex
	RateMap      map[string]map[int]int
	LimitPerHour int
}

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

func DefaultCORSConfig() CORSConfig {
	allowedOrigins := []string{"http://localhost:3000"}
	if envOrigins := os.Getenv("ALLOWED_ORIGINS"); envOrigins != "" {
		allowedOrigins = append(allowedOrigins, strings.Split(envOrigins, ",")...)
	}

	return CORSConfig{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}
}

func (a *ApiConfig) authorizeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "no jwt token found", err)
			return
		}
		user_id, err := auth.ValidateJWT(token, a.Secret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "error validating jwt", err)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, user_id)
		next(w, r.WithContext(ctx))
	})
}

func (rl *RateLimit) rateLimitMiddleware(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIPAddr(r)
		hour := time.Now().Hour()

		rl.Mu.Lock()
		// 1. MUST defer here so every 'return' path unlocks the Mutex
		defer rl.Mu.Unlock()

		// 2. Initialize if IP is new
		if rl.RateMap[ip] == nil {
			rl.RateMap[ip] = make(map[int]int)
		}

		// 3. Increment
		rl.RateMap[ip][hour]++

		// 4. Check limit
		if rl.RateMap[ip][hour] > rl.LimitPerHour {
			// We removed the 'go resetRate' because the Nuke handles it now
			respondWithError(w, http.StatusTooManyRequests, "Too many requests", nil)
			return
		}

		next(w, r)
	})
}

func getClientIPAddr(r *http.Request) string {
	// Check common headers for the client IP when behind a proxy
	for _, header := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		ips := r.Header.Get(header)
		if ips != "" {
			// The first IP in X-Forwarded-For is generally the original client IP
			ip := strings.Split(ips, ",")[0]
			ip = strings.TrimSpace(ip)
			// Optional: validate the IP address
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	// Fallback to RemoteAddr if no proxy headers are found or valid
	// RemoteAddr includes both IP and port (e.g., "127.0.0.1:54321")
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}

	return r.RemoteAddr // Fallback to raw RemoteAddr if SplitHostPort fails
}

func (rl *RateLimit) StartGlobalNuke() {
	// Ticker fires every hour
	ticker := time.NewTicker(1 * time.Hour)

	go func() {
		for range ticker.C {
			rl.Mu.Lock()
			// The Nuke: completely re-initialize the map
			// This clears all IPs and all counts instantly
			rl.RateMap = make(map[string]map[int]int)
			rl.Mu.Unlock()

			fmt.Println("Global RateLimit Nuke performed: Map cleared.")
		}
	}()
}

func CORS(config CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, o := range config.AllowedOrigins {
				if o == "*" || o == origin {
					allowed = true
					break
				}
			}

			if allowed {
				// Set CORS headers
				if config.AllowedOrigins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}

				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))

				// Handle preflight OPTIONS request
				if r.Method == "OPTIONS" {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
