package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/username/xyz-multifinance/internal/handler"
	"github.com/username/xyz-multifinance/internal/infrastructure/db"
	"github.com/username/xyz-multifinance/internal/infrastructure/security"
	"github.com/username/xyz-multifinance/internal/repository"
	"github.com/username/xyz-multifinance/internal/usecase"
)

func main() {
	// ✅ Load .env file (WAJIB di Go)
	_ = godotenv.Load()

	// Load database configuration from environment variables
	dbConfig := db.Config{
		Host:   getEnv("DB_HOST", "localhost"),
		Port:   getEnv("DB_PORT", "3306"),
		User:   getEnv("DB_USER", "root"),
		DBName: getEnv("DB_NAME", "xyz_multifinance"),
	}

	// ✅ Handle password dengan BENAR (root tanpa password)
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword != "" {
		dbConfig.Password = dbPassword
	}

	// Debug (boleh hapus setelah sukses)
	log.Printf("DB_USER=%s DB_PASSWORD=[%s]", dbConfig.User, dbConfig.Password)

	// Initialize database connection
	database, err := db.NewMySQLConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	consumerRepo := repository.NewConsumerRepository(database)
	limitRepo := repository.NewLimitRepository(database)
	transactionRepo := repository.NewTransactionRepository(database)

	// Initialize use cases
	transactionUsecase := usecase.NewTransactionUsecase(
		transactionRepo,
		limitRepo,
		consumerRepo,
		database,
	)

	// Initialize handlers
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)
	authHandler := handler.NewAuthHandler(consumerRepo)

	// Setup routers
	publicMux := http.NewServeMux()
	protectedMux := http.NewServeMux()

	// Public endpoints (no authentication required)
	publicMux.HandleFunc("/health", handler.HealthCheck)
	publicMux.HandleFunc("/api/v1/auth/login", authHandler.Login)

	// Protected endpoints (authentication required)
	protectedMux.HandleFunc("/api/v1/transactions", transactionHandler.CreateTransaction)
	protectedMux.HandleFunc("/api/v1/transactions/detail", transactionHandler.GetTransaction)
	protectedMux.HandleFunc("/api/v1/transactions/consumer", transactionHandler.GetConsumerTransactions)
	protectedMux.HandleFunc("/api/v1/limits", transactionHandler.GetConsumerLimits)

	// Apply middleware to protected routes
	rateLimiter := security.NewRateLimiter()
	protectedHandler := security.AuthMiddleware(
		security.InputValidationMiddleware(protectedMux))

	publicHandler := security.InputValidationMiddleware(publicMux)

	// Combine both handlers
	mainMux := http.NewServeMux()
	mainMux.Handle("/health", publicHandler)
	mainMux.Handle("/api/v1/auth/", publicHandler)
	mainMux.Handle("/api/v1/", protectedHandler)

	// Apply common middleware to all routes
	handlerChain := security.LoggingMiddleware(
		security.SecurityHeadersMiddleware(
			rateLimiter.Middleware(mainMux)))

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("🚀 PT XYZ Multifinance API Server starting on port %s", port)

	if err := http.ListenAndServe(":"+port, handlerChain); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
