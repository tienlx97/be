package main

import (
	"be/http"
	"be/log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const defaultPort = ":8080"

func main() {
	// init log
	log.InitLog()

	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := gin.Default()

	r.GET("/", http.PlaygroundHandler())
	r.POST("/query", http.GraphqlHandler())

	r.Run(port)
}
