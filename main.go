package main

import (
	"OTP-Go/internals/mail"
	"OTP-Go/internals/models"
	"OTP-Go/internals/otp"
	"OTP-Go/internals/redis"
	"time"

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

	client.POST("/otp", SendOTP)
	client.POST("/verify", VerifyOTP)
	client.GET("/health", health)
	err = client.Run(":8080")
	if err != nil {
		return
	}
}

func health(cntx *gin.Context) {
	cntx.JSON(200, gin.H{
		"redis": "connected",
	})
}

func SendOTP(cntx *gin.Context) {
	var email models.EmailModel

	if err := cntx.ShouldBindJSON(&email); err != nil {
		cntx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if email.Email == "" {
		cntx.JSON(400, gin.H{"error": "email is required"})
		return
	}

	generateOTP, err1 := otp.GenerateOTP()
	if err1 != nil {
		cntx.JSON(500, gin.H{
			"error": "error generating generateOTP: " + err1.Error(),
		})
		return
	}
	hash, err2 := otp.HashOTP(generateOTP)

	if err2 != nil {
		cntx.JSON(500, gin.H{
			"error": "error hashing OTP: " + err2.Error(),
		})
		return
	}
	ctx := cntx.Request.Context()
	errs := redis.RDB.Set(ctx, "otp:email:"+email.Email, hash, 300*time.Second).Err()
	if errs != nil {
		cntx.JSON(500, gin.H{
			"error": "error saving OTP: " + errs.Error(),
		})
		return
	}

	err5 := mail.SendOTPEmail(email.Email, generateOTP)
	if err5 != nil {
		cntx.JSON(500, gin.H{
			"error": "failed to send email: " + err5.Error(),
		})
		return
	}

	cntx.JSON(200, gin.H{
		"message": "OTP sent successfully",
	})
}

func VerifyOTP(cntx *gin.Context) {
	var req models.VerifyModel

	if err := cntx.ShouldBindJSON(&req); err != nil {
		cntx.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	if req.Email == "" || req.Code == "" {
		cntx.JSON(400, gin.H{"error": "email and code are required"})
		return
	}

	ctx := cntx.Request.Context()
	key := "otp:email:" + req.Email

	hash, err := redis.RDB.Get(ctx, key).Result()
	if err != nil {
		cntx.JSON(401, gin.H{"error": "otp expired or not found"})
		return
	}

	if !otp.VerifyOTP(hash, req.Code) {
		cntx.JSON(401, gin.H{"error": "invalid otp"})
		return
	}

	_ = redis.RDB.Del(ctx, key).Err()

	cntx.JSON(200, gin.H{
		"message": "otp verified successfully",
	})

}
