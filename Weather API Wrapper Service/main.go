package main

import (
    "os"

    "github.com/RohithBN/handler"
    "github.com/RohithBN/redis"
    "github.com/gin-gonic/gin"
)

func init() {
    redis.ConnectToRedis()
}

func main() {
    r := gin.Default()

    r.POST("/api/weather", handler.FetchWeatherData)
    
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    r.Run(":" + port)
  
}