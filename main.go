package main

import (
	"asset-dairy/db"
	"asset-dairy/handlers"
	"asset-dairy/middleware"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading .env file")
	}

	dbConn := db.InitDB()
	authHandler := handlers.NewAuthHandler(dbConn)

	r := gin.Default()

	// CORS middleware for frontend on port 5173
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

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
		protected.POST("/auth/change-password", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"}) })
		protected.GET("/profile", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"}) })
		protected.PUT("/profile", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"}) })
		protected.GET("/accounts", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"}) })
		protected.POST("/accounts", func(c *gin.Context) { c.JSON(http.StatusNotImplemented, gin.H{"message": "Not implemented"}) })
		// Add more routes as needed
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
