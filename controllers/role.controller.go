package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
)

// response struct
type roleResponse struct {
	Role        string    `json:"role"`
	AccessToken string    `json:"access_token"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (apiConfig *ApiConfig) HandleCreateRole(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	// request struct
	type role struct {
		Role string `json:"role"`
	}

	// extracting role name from request body
	decoder := json.NewDecoder(r.Body)
	params := role{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if params.Role == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "Invalid Role Name")
		return
	}

	// creating new role
	newRole, err := apiConfig.DB.CreateRole(r.Context(), params.Role)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, roleResponse{
		Role:        newRole.RoleName,
		AccessToken: newAccessToken,
		CreatedAt:   newRole.CreatedAt,
		UpdatedAt:   newRole.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleRemoveRole(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	// request struct
	type roleRequest struct {
		RoleID uuid.UUID `json:"role_id"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := roleRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// removing the role
	if params.RoleID == uuid.Nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, "invalid role id")
		return
	}
	if err = apiConfig.DB.RemoveRole(r.Context(), params.RoleID); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleUpdateRole(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	// request struct
	type updateRoleRequest struct {
		Role   string    `json:"role"`
		RoleID uuid.UUID `json:"role_id"`
	}

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := updateRoleRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// updating role
	if params.Role == "" || params.RoleID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid role")
		return
	}
	roleUpdate, err := apiConfig.DB.UpdateRoleById(r.Context(), database.UpdateRoleByIdParams{
		RoleName: params.Role,
		ID:       params.RoleID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, roleResponse{
		Role:        roleUpdate.RoleName,
		AccessToken: newAccessToken,
		CreatedAt:   roleUpdate.CreatedAt,
		UpdatedAt:   roleUpdate.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleGetAllRoles(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	// response struct
	type AllRoles struct {
		Roles       []database.GetAllRolesRow `json:"roles"`
		AccessToken string                    `json:"access_token"`
	}
	allRoles, err := apiConfig.DB.GetAllRoles(r.Context())
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, AllRoles{
		Roles:       allRoles,
		AccessToken: newAccessToken,
	})
}
