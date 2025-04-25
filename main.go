package main

import (
	"asset-dairy/db"
	"asset-dairy/handlers"
	"asset-dairy/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	authHandler := handlers.NewAuthHandler(dbConn)

	r := gin.Default()

	// CORS middleware for frontend on port 5173
	r.Use(func(c *gin.Context) {
		allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
		originList := strings.Split(allowedOrigins, ",")
		requestOrigin := c.Request.Header.Get("Origin")
		allowed := false
		for _, o := range originList {
			log.Println("Request origin:", requestOrigin)
			if strings.TrimSpace(o) == requestOrigin {
				allowed = true
				log.Println("Allowed origin:", o)

				break
			}
		}
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", requestOrigin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ngrok-skip-browser-warning")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	profileHandler := handlers.NewProfileHandler(dbConn)
	accountHandler := handlers.NewAccountHandler(dbConn)
	tradeHandler := handlers.NewTradeHandler(dbConn)

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
