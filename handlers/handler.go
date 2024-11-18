package handlers

import (
	"fmt"
	"log"
	"otp-auth-system/database"
	"otp-auth-system/models"
	"otp-auth-system/service"
	"otp-auth-system/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Login(c *gin.Context) {
	// Create a struct to receive the phone number from request body
	var loginRequest models.Login
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format", "error": err.Error()})
		return
	}

	// Check if phone number is empty
	if loginRequest.PhoneNumber == "" {
		c.JSON(400, gin.H{"message": "Phone number is required"})
		return
	}

	// First check if the user exists
	var user models.User
	err := database.MI.DB.Collection("users").FindOne(c, bson.M{"phonenumber": loginRequest.PhoneNumber}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"message": "User with that number not found", "error": err.Error()})
		return
	}

	// Generate and send OTP
	sid, err := service.TwilioSendOTP(loginRequest.PhoneNumber)
	log.Println("Sending OTP to", loginRequest.PhoneNumber)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to send OTP", "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "OTP sent successfully to registered number", "sid": sid})
}

func VerifyOTP(c *gin.Context) {
	var otp models.OTPResponse
	if err := c.BindJSON(&otp); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format", "error": err.Error()})
		return
	}

	if otp.Response == nil || otp.Response.PhoneNumber == "" {
		c.JSON(400, gin.H{"message": "Phone number is required"})
		return
	}

	err := service.TwilioVerifyOTP(otp.Response.PhoneNumber, fmt.Sprintf("%d", otp.Code))
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to verify OTP", "error": err.Error()})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(otp.Response.PhoneNumber)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to generate token", "error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "OTP verified successfully",
		"token":   token,
	})
}

func CreateUser(c *gin.Context) {
	var newUser models.User
	c.BindJSON(&newUser)

	// First check if the user already exists
	var existingUser models.User
	err := database.MI.DB.Collection("users").FindOne(c, bson.M{"phonenumber": newUser.PhoneNumber}).Decode(&existingUser)
	if err == nil {
		c.JSON(409, gin.H{"message": "User already exists"})
		return
	}

	_, err = database.MI.DB.Collection("users").InsertOne(c, newUser)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create user"})
		return
	}
	c.JSON(201, gin.H{"message": "User created successfully"})
}

func FindUser(c *gin.Context) {
	var number string
	c.BindJSON(&number)
	_, err := database.MI.DB.Collection("users").Find(c, bson.M{"phone_number": number})
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

}

func GetUserDetails(c *gin.Context) {
	// Get phone number from JWT token (set by middleware)
	phoneNumber, exists := c.Get("phone_number")
	if !exists {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	// Find user in database
	var user models.User
	err := database.MI.DB.Collection("users").FindOne(
		c,
		bson.M{"phone_number": phoneNumber.(string)},
	).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "User not found",
			"error":   err.Error(),
		})
		return
	}
	log.Println(phoneNumber)

	c.JSON(200, gin.H{
		"message": "User details retrieved successfully",
		"user":    user,
	})
}
