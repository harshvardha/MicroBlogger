package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (apiConfig *ApiConfig) HandleSendOTP(w http.ResponseWriter, r *http.Request, registrationOrLoginRequest *RegistrationOrLoginRequest) {

}

func (apiConfig *ApiConfig) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {

}

func (apiConfig *ApiConfig) HandleResendOTP(w http.ResponseWriter, r *http.Request) {

}

func (apiConfig *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request, registrationOrLoginRequest *RegistrationOrLoginRequest) {

}

func MakeJWT(tokenClaims struct{ UserId string }, tokenSecret string, expiresIn time.Duration) (string, error) {
	// creating the signing key to be used for signing the token
	signingKey := []byte(tokenSecret)

	// creating token claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "http://localhost:8080",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   tokenClaims.UserId,
	}

	// signing the access token with the signing key
	accessToken := jwt.NewWithClaims(jwt.SigningMethodES512, claims)
	signedAccessToken, err := accessToken.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return signedAccessToken, nil
}

func generateRefreshToken() (string, error) {
	refreshToken := make([]byte, 32)
	_, err := rand.Read(refreshToken)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(refreshToken), nil
}
