package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/harshvardha/artOfSoftwareEngineering/controllers"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
)

type authHandler func(http.ResponseWriter, *http.Request, *controllers.IDAndRole, *string)

func ValidateJWT(handler authHandler, tokenSecret string, db *database.Queries, dataValidator *Validator) http.HandlerFunc {
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
		userIDString, err := token.Claims.GetSubject()
		if err != nil {
			utility.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		userID, err := uuid.Parse(userIDString)
		if err != nil {
			utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
			return
		}

		// fetching the role of the user
		userRole, err := db.GetUserRole(r.Context(), userID)
		if err != nil {
			utility.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		UserRoleAndId := controllers.IDAndRole{
			ID:   userID,
			Role: userRole,
		}
		urlRequested := r.URL.Path

		if parseError != nil {
			// checking if the access token is expired or not
			tokenExpiresAt, err := token.Claims.GetExpirationTime()
			if err != nil {
				utility.RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}

			// if access token is expired then we will check if the refresh token is expired or not
			// if refresh token is not expired then we will create a new access token and continue
			// if refresh token is expired then we will ask user to login again
			if time.Now().After(tokenExpiresAt.Time) {
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
					newAccessToken, err := controllers.MakeJWT(struct{ UserId string }{UserId: userIDString}, tokenSecret, time.Hour)
					if err != nil {
						utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
						return
					}

					if urlRequested == "/api/user/updateEmail" || urlRequested == "/api/user/updatePassword" {
						OtherHandler(validationOtherHandler(handler), dataValidator, &UserRoleAndId, &newAccessToken)
					}

					handler(w, r, &UserRoleAndId, &newAccessToken)
				}
			}

			utility.RespondWithError(w, http.StatusUnauthorized, parseError.Error())
			return
		}

		if urlRequested == "/api/user/updateEmail" || urlRequested == "/api/user/updatePassword" {
			OtherHandler(validationOtherHandler(handler), dataValidator, &UserRoleAndId, nil)
		}

		handler(w, r, &UserRoleAndId, nil)
	}
}
