package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
	"golang.org/x/crypto/bcrypt"
)

func (apiConfig *ApiConfig) HandleSendOTP(w http.ResponseWriter, r *http.Request) {
	// extracting email from request body
	type email struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := email{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// validating email
	if err = apiConfig.DataValidator.Var(params.Email, "required,email"); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// sending otp
	if err = apiConfig.OtpCache.SendOTP(params.Email); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, nil)
}

func (apiConfig *ApiConfig) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	// extracting email, password, verificationToken and otp from request body
	decoder := json.NewDecoder(r.Body)
	params := registrationRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// validating email
	if err = apiConfig.DataValidator.Var(params.Email, "required,email"); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking if the user with this email already exist
	userExists, _ := apiConfig.DB.UserExist(r.Context(), params.Email)
	if userExists {
		utility.RespondWithError(w, http.StatusBadRequest, "user already exist")
		return
	}

	// validating the otp
	if err = apiConfig.OtpCache.VerifyOTP(params.VerificationToken, params.OTP); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// validating password
	if err = apiConfig.DataValidator.Var(params.Password, "required,min=6,max=64,password"); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.MaxCost)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// registering new user
	roleID, err := apiConfig.DB.GetRoleIdByName(r.Context(), "user")
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = apiConfig.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:         strings.ToLower(params.Email),
		Username:      strings.ToLower(params.Username),
		Password:      string(hashedPassword),
		ProfilePicUrl: params.ProfilePicUrl,
		RoleID:        roleID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, nil)
}

func (apiConfig *ApiConfig) HandleResendOTP(w http.ResponseWriter, r *http.Request) {
	// extracting old verification token and email from request body
	type otpCredentials struct {
		OldVerificationToken string `json:"old_verification_token"`
		Email                string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := otpCredentials{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking if resend is allowed or not
	if isResendAllowed, err := apiConfig.OtpCache.IsResendAllowed(params.OldVerificationToken); !isResendAllowed {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking the validation of email
	if err = apiConfig.DataValidator.Var(params.Email, "required,email"); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// resending the otp
	if err = apiConfig.OtpCache.ResendOTP(params.OldVerificationToken, params.Email); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, nil)
}

func (apiConfig *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// extracting email and password from request body
	decoder := json.NewDecoder(r.Body)
	params := loginRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// validating email and password
	if err = apiConfig.DataValidator.Var(params.Email, "required,email"); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err = apiConfig.DataValidator.Var(params.Password, "required,min=6,max=64,password"); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking if the user exist or not
	user, err := apiConfig.DB.GetUserByEmailID(r.Context(), params.Email)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// comparing passwords
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// creating access token
	accessToken, err := MakeJWT(struct{ UserId string }{UserId: user.ID.String()}, apiConfig.JwtSecret, 2*time.Hour)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// generating refresh token
	refreshToken, err := generateRefreshToken()
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err = apiConfig.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	}); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, loginResponse{
		Username:      user.Username,
		ProfilePicUrl: user.ProfilePicUrl,
		AccessToken:   accessToken,
	})
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
