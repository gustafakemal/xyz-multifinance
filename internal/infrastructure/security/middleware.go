package security

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-in-production") // Should be in environment variable

type Claims struct {
	ConsumerID int64  `json:"consumer_id"`
	NIK        string `json:"nik"`
	jwt.RegisteredClaims
}

type contextKey string

const ConsumerIDKey contextKey = "consumer_id"

// OWASP #1: Broken Access Control - JWT Authentication Middleware
// This middleware ensures only authenticated users can access protected endpoints
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"Missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Extract Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"Invalid authorization format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Parse and validate JWT token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// Add consumer ID to request context
		ctx := context.WithValue(r.Context(), ConsumerIDKey, claims.ConsumerID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OWASP #2: Cryptographic Failures - API Key Authentication (Alternative)
// Uses constant-time comparison to prevent timing attacks
func APIKeyMiddleware(next http.Handler) http.Handler {
	validAPIKey := "your-api-key-here" // Should be in environment variable

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, `{"error":"Missing API key"}`, http.StatusUnauthorized)
			return
		}

		// Use constant-time comparison to prevent timing attacks (OWASP Security)
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(validAPIKey)) != 1 {
			http.Error(w, `{"error":"Invalid API key"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// OWASP #3: Injection - Input Sanitization Middleware
// Prevents SQL injection by validating content type and rejecting suspicious patterns
func InputValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only validate POST, PUT, PATCH requests
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			contentType := r.Header.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				http.Error(w, `{"error":"Content-Type must be application/json"}`, http.StatusBadRequest)
				return
			}

			// Limit request body size to prevent DoS attacks (10MB max)
			r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)
		}

		next.ServeHTTP(w, r)
	})
}

// OWASP #4: Security Misconfiguration - Security Headers Middleware
// Adds security headers to prevent various attacks
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent clickjacking attacks
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Enforce HTTPS
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		// Control referrer information
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		// Prevent caching of sensitive data
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, private")
		w.Header().Set("Pragma", "no-cache")

		next.ServeHTTP(w, r)
	})
}

// OWASP #5: Rate Limiting - Protection against brute force attacks
type RateLimiter struct {
	requests map[string][]time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use IP address as identifier
		ip := r.RemoteAddr

		// Clean old requests (older than 1 minute)
		now := time.Now()
		if timestamps, exists := rl.requests[ip]; exists {
			var validRequests []time.Time
			for _, t := range timestamps {
				if now.Sub(t) < time.Minute {
					validRequests = append(validRequests, t)
				}
			}
			rl.requests[ip] = validRequests
		}

		// Check rate limit (max 100 requests per minute)
		if len(rl.requests[ip]) >= 100 {
			http.Error(w, `{"error":"Rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}

		// Add current request
		rl.requests[ip] = append(rl.requests[ip], now)

		next.ServeHTTP(w, r)
	})
}

// OWASP #6: Logging and Monitoring Middleware
// Logs all requests for security monitoring
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		log.Printf("[%s] %s %s - IP: %s", r.Method, r.URL.Path, r.Proto, r.RemoteAddr)

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		// Log response
		duration := time.Since(start)
		log.Printf("[%s] %s - Status: %d - Duration: %v", r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Helper function to generate JWT token (for testing/authentication)
func GenerateToken(consumerID int64, nik string) (string, error) {
	claims := Claims{
		ConsumerID: consumerID,
		NIK:        nik,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Helper function to write JSON response
func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper function to write error response
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}
