package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/harshvardha/artOfSoftwareEngineering/controllers"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
)

func getSubjects(tokenSubjects string) (map[string]string, error) {
	subjects := strings.Split(tokenSubjects, ",")
	if len(subjects) == 0 {
		return nil, errors.New("no subject found")
	}

	subjectsMap := make(map[string]string, len(subjects))
	var temp []string
	for _, subject := range subjects {
		temp = strings.Split(subject, ":")
		if len(temp) == 0 {
			return nil, errors.New("not all claims found")
		}
		subjectsMap[temp[0]] = temp[1]
	}

	return subjectsMap, nil
}

func ValidateJWT(handler authenticatedRequestHandler, tokenSecret string, db *database.Queries, routes map[string][]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extracting JWT token from request header
		authHeader := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authHeader) != 2 {
			utility.RespondWithError(w, http.StatusNotAcceptable, "malformed request auth header")
			return
		}

		// initializing an empty struct to parse claims
		jwtClaims := jwt.RegisteredClaims{}
		token, parseError := jwt.ParseWithClaims(authHeader[1], &jwtClaims, func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})

		// extracting userID from token claims
		subjects, err := token.Claims.GetSubject()
		if err != nil {
			utility.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		parsedSubjects, err := getSubjects(subjects)
		if err != nil {
			utility.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userID, err := uuid.Parse(parsedSubjects["user_id"])
		if err != nil {
			utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
			return
		}
		UserRoleAndId := controllers.IDAndRole{
			ID:   userID,
			Role: parsedSubjects["role"],
		}

		if parseError != nil {
			// checking if the access token is expired or not
			// if access token is expired then we will check if the refresh token is expired or not
			// if refresh token is not expired then we will create a new access token and continue
			// if refresh token is expired then we will ask user to login again
			if errors.Is(parseError, jwt.ErrTokenExpired) {
				refreshTokenExpiresAt, err := db.GetRefreshTokenExpirationTime(r.Context(), userID)
				if err != nil {
					utility.RespondWithError(w, http.StatusUnauthorized, err.Error())
					return
				}

				// checking if refresh token is expired or not
				// if refresh token is expired then login again
				// otherwise create new access token and send it with response
				if time.Now().After(refreshTokenExpiresAt) {
					utility.RespondWithError(w, http.StatusUnauthorized, "Please login again")
					return
				} else {
					// creating new access token
					newAccessToken, err := controllers.MakeJWT(struct {
						UserId string
						Role   string
					}{
						UserId: userID.String(),
						Role:   parsedSubjects["role"],
					}, tokenSecret, time.Hour)
					if err != nil {
						utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
						return
					}

					// calling the authorization middleware to check whether the user is authorized to access this endpoint
					userAuthorization(w, r, handler, routes, &UserRoleAndId, newAccessToken)
				}
			}

			utility.RespondWithError(w, http.StatusUnauthorized, parseError.Error())
			return
		}

		// calling the authorization middleware to check whether the user is authorized to access this endpoint
		userAuthorization(w, r, handler, routes, &UserRoleAndId, "")
	}
}
