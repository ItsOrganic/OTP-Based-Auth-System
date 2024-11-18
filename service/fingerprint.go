package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"otp-auth-system/models"
)

type IPInfo struct {
	IP      string `json:"ip"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
}

type FingerprintService struct{}

func NewFingerprintService() *FingerprintService {
	return &FingerprintService{}
}

func (s *FingerprintService) GenerateDeviceID(userAgent string, ip string) string {
	// Create a unique device identifier
	data := fmt.Sprintf("%s-%s", userAgent, ip)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (s *FingerprintService) GetLocationInfo(ip string) (*IPInfo, error) {
	// Using ipapi.co for IP geolocation (free tier)
	resp, err := http.Get(fmt.Sprintf("https://ipapi.co/%s/json/", ip))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ipInfo IPInfo
	if err := json.Unmarshal(body, &ipInfo); err != nil {
		return nil, err
	}

	return &ipInfo, nil
}

func (s *FingerprintService) IsKnownDevice(deviceInfo models.DeviceInfo, knownDevices []models.DeviceInfo) bool {
	for _, known := range knownDevices {
		if known.DeviceID == deviceInfo.DeviceID {
			return true
		}
	}
	return false
}
