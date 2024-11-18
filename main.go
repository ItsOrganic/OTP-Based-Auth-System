package main

import (
	"otp-auth-system/database"
	"otp-auth-system/handlers"
	"otp-auth-system/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	database.ConnectDB()
	router := gin.Default()
	router.POST("/create-user", handlers.CreateUser)
	router.POST("/login", handlers.Login)
	router.POST("/verify-otp", handlers.VerifyOTP)
	protected := router.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		router.GET("/profile", handlers.GetUserDetails)
	}
	router.Run(":8080")
}
