package main

import (
	"fmt"
	"log"
	"midterm/db"
	"midterm/routers"
	"midterm/utils"

	"github.com/gin-gonic/gin"
)

func init() {
	_, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} else {
		fmt.Println("Connected to database")
	}

	err = utils.ConnectToSpaces()
	if err != nil {
		log.Fatalf("Failed to connect to space: %v", err)
	} else {
		fmt.Println("Connected to space")
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	route := gin.Default()

	// Apply CORS middleware globally
	route.Use(CORSMiddleware())

	// Initialize routes
	routers.InitRoute(route)

	// Start server
	route.Run(":8080")
}
