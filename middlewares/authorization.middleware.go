package middlewares

import (
	"net/http"
	"slices"

	"github.com/harshvardha/artOfSoftwareEngineering/controllers"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
)

type Routes struct {
	UserRoutes  []string
	AdminRoutes []string
}

func (routes *Routes) userAuthorization(w http.ResponseWriter, r *http.Request, handler authenticatedRequestHandler, IDAndRole *controllers.IDAndRole, db *database.Queries, newAccessToken *string) {
	// extracting the endpoint user trying to access
	endpoint := r.URL.Path

	// checking if the user is authorized to access the endpoint according to the role
	userRole, err := db.GetRoleById(r.Context(), IDAndRole.Role)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	switch userRole {
	case "admin":
		if !slices.Contains(routes.AdminRoutes, endpoint) {
			utility.RespondWithError(w, http.StatusBadRequest, "endpoint not available")
			return
		}
		break
	case "user":
		if !slices.Contains(routes.UserRoutes, endpoint) {
			utility.RespondWithError(w, http.StatusBadRequest, "endpoint not available")
			return
		}
		break
	default:
		utility.RespondWithError(w, http.StatusNotFound, "endpoint does not exist")
		return
	}

	handler(w, r, IDAndRole, newAccessToken)
}
