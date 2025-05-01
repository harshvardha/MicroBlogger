package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/harshvardha/artOfSoftwareEngineering/controllers"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
)

type Validator struct {
	Validate *validator.Validate
}

type validationAuthHandler func(http.ResponseWriter, *http.Request, *controllers.RegistrationOrLoginRequest)
type validationOtherHandler func(http.ResponseWriter, *http.Request, *controllers.IDAndRole, *string)

func AuthValidation(handler validationAuthHandler, dataValidator *Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extracting email and password from request body
		decoder := json.NewDecoder(r.Body)
		params := controllers.RegistrationOrLoginRequest{}
		err := decoder.Decode(&params)
		if err != nil {
			utility.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		// validating email
		emailError := dataValidator.Validate.Var(params.Email, "required,email")

		// validating password
		passwordError := dataValidator.Validate.Var(params.Password, "required,min=6,max=64,password")

		if emailError != nil {
			utility.RespondWithError(w, http.StatusBadRequest, emailError.Error())
			return
		}

		if passwordError != nil {
			utility.RespondWithError(w, http.StatusBadRequest, passwordError.Error())
			return
		}

		handler(w, r, &params)
	}
}

func OtherHandler(handler validationOtherHandler, dataValidator *Validator, IDAndRole *controllers.IDAndRole, newAccessToken *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extracting email or password from request body based on the request url
		apiEndPoint := r.URL.Path
		if apiEndPoint == "" {
			utility.RespondWithError(w, http.StatusNoContent, "no endpoint provided")
			return
		}

		decoder := json.NewDecoder(r.Body)
		params := controllers.RegistrationOrLoginRequest{}
		err := decoder.Decode(&params)
		if err != nil {
			utility.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if apiEndPoint == "/api/user/updateEmail" {
			err = dataValidator.Validate.Var(params.Email, "required,email")
			if err != nil {
				utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
				return
			}
		} else if apiEndPoint == "/api/user/updatePassword" {
			err = dataValidator.Validate.Var(params.Password, "required,min=6,max=64,password")
			if err != nil {
				utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
				return
			}
		}

		handler(w, r, IDAndRole, newAccessToken)
	}
}
