package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mikecodesun/backend-sql/routes"
)

func main() {
	// .env varible port
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error laoding .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// gin route
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "This is Home Page"})
	})
	// routes user/post
	routes.PostRoute(r)
	routes.UserRoute(r)
	// server running port
	fmt.Printf(" server is running on port:%s ", port)
	r.Run(":" + port)
}