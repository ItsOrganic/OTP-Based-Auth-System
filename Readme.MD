# OTP Authentication System

A secure phone authentication system built with Go, featuring OTP verification, device fingerprinting, and location tracking.

## Features
- 📱 Phone number authentication with OTP
- 🔐 JWT-based session management
- 📍 Device fingerprinting and location tracking
- 🚨 New device detection alerts
- 📦 MongoDB for data persistence
- 🔄 Twilio integration for OTP delivery

## Requirements
- Go 1.22.5
- MongoDB
- Twilio Account

## Environment Setup
Create a `.env` file in the root directory:
env
```
TWILIO_ACCOUNT_SID=your_sid
TWILIO_AUTHTOKEN=your_token
TWILIO_SERVICES_ID=your_service_id
```
## Installation
Clone the repository
git clone https://github.com/yourusername/otp-auth-system.git && cd otp-auth-system
Install dependencies
```
go mod tidy
```
## API Endpoints

### Public Routes

#### 1. Create User
http
POST /create-user
Content-Type: application/json
{
"phone_number": "+1234567890",
"name": "John Doe",
"email": "john@example.com",
"age": 25
}

#### 2. Login (Request OTP)
http
POST /login
Content-Type: application/json
{
"phone_number": "+1234567890"
}
#### 3. Verify OTP
POST /verify-otp
Content-Type: application/json
{
"number": {
"phone_number": "+1234567890"
},
"code": 123456
}
### Protected Routes

#### 4. Get User Profile
http
GET /profile
Authorization: Bearer <your_jwt_token>


## Project Structure
otp-auth-system/
├── config/
│ └── config.go # Configuration management
├── database/
│ └── database.go # MongoDB connection
├── handlers/
│ └── handler.go # Request handlers
├── middlewares/
│ └── middleware.go # JWT authentication
├── models/
│ └── model.go # Data models
├── service/
│ ├── service.go # Twilio OTP service
│ └── fingerprint.go # Device fingerprinting
├── utils/
│ └── jwt.go # JWT utilities
├── main.go # Application entry point
├── go.mod
└── .env

## Authentication Flow
1. User registers with phone number and details
2. System captures initial device information
3. Login request triggers OTP generation
4. New device detection during login
5. OTP verification generates JWT token
6. Protected routes require valid JWT
7. Device history maintained for security