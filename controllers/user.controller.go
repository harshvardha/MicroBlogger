package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
	"golang.org/x/crypto/bcrypt"
)

func (apiConfig *ApiConfig) HandleUpdateEmail(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type UpdateEmailRequest struct {
		VerificationToken string `json:"verificationToken"`
		OTP               string `json:"otp"`
		NewEmail          string `json:"newEmail"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := UpdateEmailRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking if otp is valid and correct
	if err = apiConfig.OtpCache.VerifyOTP(params.VerificationToken, params.OTP); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking if the email is valid or not
	if err = apiConfig.DataValidator.Var(params.NewEmail, "required,email"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	// updating new email
	if err = apiConfig.DB.UpdateEmail(r.Context(), database.UpdateEmailParams{
		Email: params.NewEmail,
		ID:    IDAndRole.ID,
	}); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// removing existing refresh token
	if err = apiConfig.DB.RemoveRefreshToken(r.Context(), IDAndRole.ID); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, nil)
}

func (apiConfig *ApiConfig) HandleUpdatePassword(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type UpdatePasswordRequest struct {
		VerificationToken string `json:"verificationToken"`
		OTP               string `json:"otp"`
		NewPassword       string `json:"newPassword"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := UpdatePasswordRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking if the otp is valid or not
	if err = apiConfig.OtpCache.VerifyOTP(params.VerificationToken, params.OTP); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// checking if the password is valid or not
	if err = apiConfig.DataValidator.Var(params.NewPassword, "required,min=6,max=64,password"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	// hashing the new password
	newPassword, err := bcrypt.GenerateFromPassword([]byte(params.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// resetting the password
	if err = apiConfig.DB.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		Password: string(newPassword),
		ID:       IDAndRole.ID,
	}); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// removing the existing refresh token
	if err = apiConfig.DB.RemoveRefreshToken(r.Context(), IDAndRole.ID); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, nil)
}

func (apiConfig *ApiConfig) HandleUpdateOtherDetails(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type UpdateUsernameOrProfilePic struct {
		Username      string `json:"username,omitempty"`
		ProfilePicURL string `json:"profilePicUrl,omitempty"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := UpdateUsernameOrProfilePic{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking which field to update
	existingInformation, err := apiConfig.DB.GetUserByID(r.Context(), IDAndRole.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	updateUsernameOrProfilePic := database.UpdateOtherDetailsParams{}
	if apiConfig.DataValidator.Var(params.Username, "required,username") == nil {
		updateUsernameOrProfilePic.Username = params.Username
	} else {
		updateUsernameOrProfilePic.Username = existingInformation.Username
	}

	if apiConfig.DataValidator.Var(params.ProfilePicURL, "required,url") == nil {
		updateUsernameOrProfilePic.ProfilePicUrl = params.ProfilePicURL
	} else {
		updateUsernameOrProfilePic.ProfilePicUrl = existingInformation.ProfilePicUrl
	}
	updateUsernameOrProfilePic.ID = IDAndRole.ID
	// updating user information
	updatedUserInformation, err := apiConfig.DB.UpdateOtherDetails(r.Context(), updateUsernameOrProfilePic)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type UpdatedUser struct {
		UsernameAndProfilePic UpdateUsernameOrProfilePic `json:"updatedUsernameAndProfilePic"`
		AccessToken           string                     `json:"accessToken"`
	}
	utility.RespondWithJson(w, http.StatusOK, UpdatedUser{
		UsernameAndProfilePic: UpdateUsernameOrProfilePic{
			Username:      updatedUserInformation.Username,
			ProfilePicURL: updatedUserInformation.ProfilePicUrl,
		},
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleRemoveUserAccount(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type RemoveUserAccountRequest struct {
		VerificationToken string `json:"verificationToken"`
		OTP               string `json:"otp"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := RemoveUserAccountRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// verifying if the otp is correct or valid
	if err = apiConfig.OtpCache.VerifyOTP(params.VerificationToken, params.OTP); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// removing user account
	if err := apiConfig.DB.RemoveUser(r.Context(), IDAndRole.ID); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, nil)
}

func (apiConfig *ApiConfig) HandleGetUserByID(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	user, err := apiConfig.DB.GetUserByID(r.Context(), IDAndRole.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	type User struct {
		Email         string    `json:"email"`
		Username      string    `json:"username"`
		ProfilePicURL string    `json:"profilePicUrl"`
		AccessToken   string    `json:"accessToken"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
	}

	utility.RespondWithJson(w, http.StatusOK, User{
		Email:         user.Email,
		Username:      user.Username,
		ProfilePicURL: user.ProfilePicUrl,
		AccessToken:   newAccessToken,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	})
}
