package models

import "time"

type User struct {
	PhoneNumber  string       `json:"phone_number"`
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	Age          int          `json:"age"`
	KnownDevices []DeviceInfo `json:"known_devices"`
}

type Login struct {
	PhoneNumber string `json:"phone_number"`
}

type OTPResponse struct {
	Response *Login `json:"number"`
	Code     int    `json:"code"`
}
type DeviceInfo struct {
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Location  string    `json:"location"`
	LastUsed  time.Time `json:"last_used"`
	DeviceID  string    `json:"device_id"`
}
