package main

import (
	"asset-dairy/db"
	"asset-dairy/handlers"
	"asset-dairy/middleware"
	"asset-dairy/services"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	dotenvFile := fmt.Sprintf(".env.%s", env)
	if err := godotenv.Load(dotenvFile); err != nil {
		log.Printf("No %s file found or error loading %s", dotenvFile, dotenvFile)
	}

	dbConn := db.InitDB()

	// Run DB migrations before starting the server
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	driver, err := postgres.WithInstance(dbConn, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create DB driver for migrations: %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Database migrated successfully")

	r := gin.Default()

	// CORS middleware for frontend on port 5173
	r.Use(func(c *gin.Context) {
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		originList := strings.Split(allowedOrigins, ",")
		requestOrigin := c.Request.Header.Get("Origin")
		allowed := false
		for _, o := range originList {
			if strings.TrimSpace(o) == requestOrigin {
				allowed = true
				break
			}
		}
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", requestOrigin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	authService := services.NewAuthService(dbConn)
	authHandler := handlers.NewAuthHandler(authService)
	profileService := services.NewProfileService(dbConn)
	profileHandler := handlers.NewProfileHandler(profileService)
	accountHandler := handlers.NewAccountHandler(dbConn)
	tradeService := services.NewTradeService(dbConn)
	holdingService := services.NewHoldingService(tradeService)
	tradeHandler := handlers.NewTradeHandler(dbConn, tradeService)
	holdingHandler := handlers.NewHoldingHandler(holdingService)

	// Public routes
	public := r.Group("/auth")
	{
		public.POST("/sign-in", authHandler.SignIn)
		public.POST("/sign-up", authHandler.SignUp)
		public.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes (JWT middleware to be added)
	protected := r.Group("/")
	{
		protected.Use(middleware.JWTAuthMiddleware())
		protected.POST("/auth/logout", authHandler.Logout)
		protected.POST("/profile/change-password", profileHandler.ChangePassword)
		protected.GET("/profile", profileHandler.GetProfile)
		protected.PUT("/profile", profileHandler.UpdateProfile)
		protected.GET("/accounts", accountHandler.ListAccounts)
		protected.POST("/accounts", accountHandler.CreateAccount)
		protected.PUT("/accounts/:id", accountHandler.UpdateAccount)
		protected.DELETE("/accounts/:id", accountHandler.DeleteAccount)

		// Trade routes
		protected.GET("/trades", tradeHandler.ListTrades)
		protected.POST("/trades", tradeHandler.CreateTrade)
		protected.PUT("/trades/:id", tradeHandler.UpdateTrade)
		protected.DELETE("/trades/:id", tradeHandler.DeleteTrade)

		// Asset routes
		protected.GET("/holdings", holdingHandler.ListHoldings)
	}

	r.GET("/swagger/*any", ginSwaggerHandler()) // Swagger UI placeholder

	r.Run(":3000")
}

// ginSwaggerHandler is a placeholder for Swagger docs
func ginSwaggerHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Swagger docs not implemented yet"})
	}
}
