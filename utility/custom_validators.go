package utility

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func CustomPasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// check for atleast on uppercase letter
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)

	// check for atleast one lowercase letter
	hasLower := regexp.MustCompile(`[a=z]`).MatchString(password)

	// check for atleast on number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)

	// check for atleast one special character
	hasSpecialCharacter := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	return hasUpper && hasLower && hasNumber && hasSpecialCharacter
}
