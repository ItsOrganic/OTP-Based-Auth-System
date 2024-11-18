
# OTP Authentication System

A **secure phone authentication system** built with Go, featuring **OTP verification** and **device fingerprinting**.

## Features
- 📱 **Phone number authentication** with OTP
- 🔐 **JWT-based session management**
- 📍 **Device fingerprinting** and **location tracking**
- 🚨 **New device detection alerts**
- 📦 **MongoDB** for data persistence
- 🔄 **Twilio integration** for OTP delivery

---

## Requirements
- **Go** 1.22
- **MongoDB**
- **Twilio Account**

---

## Environment Setup
---

## Installation

1. **Clone the repository**:
    ```plaintext
    git clone https://github.com/itsorganic/otp-based-auth-system.git && cd otp-based-auth-system
    ```

2. **Install dependencies**:
    ```plaintext
    go mod tidy
    ```
3. Create a `.env` file in the root directory with the following configuration:

```plaintext
TWILIO_ACCOUNT_SID=your_sid
TWILIO_AUTHTOKEN=your_token
TWILIO_SERVICES_ID=your_service_id
```

4. Build and run the project. It will create executable binary in bin which will run the server
```
make run
```
---

## Architecture of the OTP-Auth-System 
<img src="https://github.com/ItsOrganic/OTP-Based-Auth-System/assets/diagram.png" align="center">

## API Endpoints

### Public Routes

#### 1. Create User
```plaintext
POST /create-user
Content-Type: application/json

{
  "phone_number": "+91XXXXXXXXX",
  "name": "Any name",
  "email": "john@example.com",
  "age": 25
}
```

#### 2. Login (Request OTP)
```plaintext
POST /login
Content-Type: application/json

{
  "phone_number": "+91XXXXXXXXX"
}
```

#### 3. Verify OTP
```plaintext
POST /verify-otp
Content-Type: application/json

{
  "number": {
    "phone_number": "+91XXXXXXXXX"
  },
  "code": 123456
}
```

### Protected Routes

#### 4. Get User Profile
```plaintext
GET /profile
Authorization: Bearer <your_jwt_token>
```

---

## Project Structure

```plaintext
otp-auth-system/
├── config/
│   └── config.go           # Configuration management
├── database/
│   └── database.go         # MongoDB connection
├── handlers/
│   └── handler.go          # Request handlers
├── middlewares/
│   └── middleware.go       # JWT authentication
├── models/
│   └── model.go            # Data models
├── service/
│   ├── service.go          # Twilio OTP service
│   └── fingerprint.go      # Device fingerprinting
├── utils/
│   └── jwt.go              # JWT utilities
├── main.go                 # Application entry point
├── go.mod
└── .env
```

---

## Authentication Flow

1. User registers with phone number and details.
2. System captures initial device information.
3. Login request triggers OTP generation.
4. System detects new devices during login.
5. OTP verification generates a JWT token.
6. Protected routes require a valid JWT token.
7. Device history is maintained for security.
