package middlewares

import (
	"net/http"
	"slices"

	"github.com/harshvardha/artOfSoftwareEngineering/controllers"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
)

func userAuthorization(w http.ResponseWriter, r *http.Request, handler authenticatedRequestHandler, routes map[string][]string, IDAndRole *controllers.IDAndRole, newAccessToken string) {
	// extracting the endpoint user trying to access
	endpoint := r.URL.Path
	var isValidRole bool

	// checking if the user is authorized to access the endpoint according to the role
	switch IDAndRole.Role {
	case "admin":
		isValidRole = true
	case "user":
		if !slices.Contains(routes[IDAndRole.Role], endpoint) {
			utility.RespondWithError(w, http.StatusBadRequest, "endpoint not available")
			return
		}
		isValidRole = true
	default:
		utility.RespondWithError(w, http.StatusNotFound, "endpoint does not exist")
		return
	}

	if isValidRole {
		if slices.Contains(routes["nil_IDAndRole"], endpoint) {
			handler(w, r, nil, newAccessToken)
			isValidRole = false
			return
		}

		if slices.Contains(routes["nil_accessToken"], endpoint) {
			handler(w, r, IDAndRole, "")
			isValidRole = false
			return
		}

		handler(w, r, IDAndRole, newAccessToken)
		isValidRole = false
	}
}
