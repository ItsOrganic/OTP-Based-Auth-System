package models

type User struct {
	PhoneNumber string `json:"phone_number"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Age         int    `json:"age"`
}

type Login struct {
	PhoneNumber string `json:"phone_number"`
}

type OTPResponse struct {
	Response *Login `json:"number"`
	Code     int    `json:"code"`
}
