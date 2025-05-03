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

type RegistrationOrLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyOTPRequest struct {
	VerificationToken string `json:"verification_token"`
	OTP               string `json:"otp"`
}
