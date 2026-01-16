package api

import (
	"otpgo/internals/mail"
	"otpgo/internals/models"
	"otpgo/internals/otp"
	"otpgo/internals/redis"
	"time"

	"github.com/gin-gonic/gin"
)

func Health(cntx *gin.Context) {
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
			"error": "error generating OTP: " + err1.Error(),
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
	errSet := redis.RDB.Set(ctx, "otp:email:"+email.Email, hash, 300*time.Second).Err()
	if errSet != nil {
		cntx.JSON(500, gin.H{
			"error": "error saving OTP: " + errSet.Error(),
		})
		return
	}

	errMail := mail.SendOTPEmail(email.Email, generateOTP)
	if errMail != nil {
		cntx.JSON(500, gin.H{
			"error": "failed to send email: " + errMail.Error(),
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
