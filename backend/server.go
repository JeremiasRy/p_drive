package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, hx-current-url, hx-request, hx-target")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.Static("/static", "./static/")
	r.LoadHTMLFiles("./static/index.html")

	r.Use(CORSMiddleware())

	r.GET("/", func(c *gin.Context) {
		backendURL := os.Getenv("BACKEND_URL")
		c.HTML(200, "index.html", gin.H{
			"BackendURL": backendURL,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
