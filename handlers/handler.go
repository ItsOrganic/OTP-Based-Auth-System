package handlers

import (
	"fmt"
	"log"
	"otp-auth-system/database"
	"otp-auth-system/models"
	"otp-auth-system/service"
	"otp-auth-system/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Create user with details number, name, age, email,  Method=POST
func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}

	err := database.MI.DB.Collection("users").FindOne(c, bson.M{"phonenumber": newUser.PhoneNumber}).Decode(&newUser)
	if err == nil {
		c.JSON(409, gin.H{"message": "User already exists"})
		return
	}

	// Get initial device info
	fingerprintService := service.NewFingerprintService()
	ip := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	deviceID := fingerprintService.GenerateDeviceID(userAgent, ip)

	ipInfo, _ := fingerprintService.GetLocationInfo(ip)
	location := "Unknown"
	if ipInfo != nil {
		location = fmt.Sprintf("%s, %s, %s", ipInfo.City, ipInfo.Region, ipInfo.Country)
	}

	// Initialize known devices with the current device
	newUser.KnownDevices = []models.DeviceInfo{{
		IP:        ip,
		UserAgent: userAgent,
		Location:  location,
		LastUsed:  time.Now(),
		DeviceID:  deviceID,
	}}

	//Verify the email format
	if !utils.VerifyEmail(newUser.Email) {
		c.JSON(400, gin.H{"message": "Invalid email format"})
		return
	}

	// Insert user into database
	_, err = database.MI.DB.Collection("users").InsertOne(c, newUser)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create user"})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully"})
}

// Login with phone number Method=POST
func Login(c *gin.Context) {
	var loginRequest models.Login
	if err := c.BindJSON(&loginRequest); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request format"})
		return
	}

	// Get fingerprint information
	fingerprintService := service.NewFingerprintService()

	// Get IP address (considering proxy headers)
	ip := c.ClientIP()

	// Get user agent
	userAgent := c.GetHeader("User-Agent")

	// Generate device ID
	deviceID := fingerprintService.GenerateDeviceID(userAgent, ip)

	// Get location info
	ipInfo, err := fingerprintService.GetLocationInfo(ip)
	if err != nil {
		// Log the error but continue
		log.Printf("Error getting location info: %v", err)
	}

	// Create device info
	currentDevice := models.DeviceInfo{
		IP:        ip,
		UserAgent: userAgent,
		Location:  fmt.Sprintf("%s, %s, %s", ipInfo.City, ipInfo.Region, ipInfo.Country),
		LastUsed:  time.Now(),
		DeviceID:  deviceID,
	}

	// Find user and their known devices
	var user models.User
	err = database.MI.DB.Collection("users").FindOne(c, bson.M{"phonenumber": loginRequest.PhoneNumber}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"message": "User with that number not found"})
		return
	}

	// Check if this is a new device
	isKnownDevice := fingerprintService.IsKnownDevice(currentDevice, user.KnownDevices)

	// Generate and send OTP
	sid, err := service.TwilioSendOTP(loginRequest.PhoneNumber)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to send OTP", "error": err.Error()})
		return
	}

	// Update user's known devices if this is a new device
	if !isKnownDevice {
		update := bson.M{
			"$push": bson.M{
				"known_devices": currentDevice,
			},
		}
		_, err = database.MI.DB.Collection("users").UpdateOne(
			c,
			bson.M{"phone_number": loginRequest.PhoneNumber},
			update,
		)
		if err != nil {
			log.Printf("Error updating known devices: %v", err)
		}
	}

	c.JSON(200, gin.H{
		"message":             "OTP sent successfully",
		"sid":                 sid,
		"new_device_detected": !isKnownDevice,
		"location":            currentDevice.Location,
	})
}

// Verify OTP and generate JWT token Method=POST
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

func FindUser(c *gin.Context) {
	var number string
	c.BindJSON(&number)
	_, err := database.MI.DB.Collection("users").Find(c, bson.M{"phonenumber": number})
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
