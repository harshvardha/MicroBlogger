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

func S3UrlValidator(fl validator.FieldLevel) bool {
	return true
}

func GithubURLValidator(fl validator.FieldLevel) bool {
	url := fl.Field().String()

	// checking if the url starts with: https://github.com/harshvardha
	pattern := `^(https?:\/\/)?(www\.)?github\.com\/harshvardha\/[a-zA-Z0-9]+(?:-?[a-zA-Z0-9]+)*(\/)?$`
	isURLValid := regexp.MustCompile(pattern).MatchString(url)
	return isURLValid
}

func TagsValidator(fl validator.FieldLevel) bool {
	tags := fl.Field().String()
	pattern := `^[a-zA-z0-9;]`
	return regexp.MustCompile(pattern).MatchString(tags)
}

func UsernameValidator(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	pattern := `^[a-zA-z0-9_]`
	return regexp.MustCompile(pattern).MatchString(username)
}

func BookNameValidator(fl validator.FieldLevel) bool {
	bookName := fl.Field().String()
	pattern := `^[a-zA-z0-9_]`
	return regexp.MustCompile(pattern).MatchString(bookName)
}
