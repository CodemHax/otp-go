package main

import (
	"otpgo/api"
	"otpgo/internals/redis"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	redis.ConnectRedis()
	gin.SetMode(gin.ReleaseMode)
	client := gin.Default()

	client.POST("/otp", api.SendOTP)
	client.POST("/verify", api.VerifyOTP)
	client.GET("/health", api.Health)
	err = client.Run(":8080")
	if err != nil {
		return
	}
}
