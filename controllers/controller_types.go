package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/cache"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
)

type ApiConfig struct {
	DB            *database.Queries
	JwtSecret     string
	OtpCache      *cache.OtpCache
	DataValidator *validator.Validate
}

type IDAndRole struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

type registrationRequest struct {
	Email             string `json:"email"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	ProfilePicUrl     string `json:"profile_pic_url"`
	VerificationToken string `json:"verification_token"`
	OTP               string `json:"otp"`
}

type VerifyOTPRequest struct {
	VerificationToken string `json:"verification_token"`
	OTP               string `json:"otp"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Username      string `json:"username"`
	ProfilePicUrl string `json:"profile_pic_url"`
	AccessToken   string `json:"access_token"`
}
