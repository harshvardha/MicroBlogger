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
type categoryResponse struct {
	Category    string    `json:"category"`
	AccessToken string    `json:"accessToken"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (apiConfig *ApiConfig) HandleCreateCategory(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	// request struct
	type createCategoryRequest struct {
		Category string `json:"category"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := createCategoryRequest{}
	if err := decoder.Decode(&params); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if params.Category == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid category")
		return
	}

	// creating new category
	newCategory, err := apiConfig.DB.CreateCategory(r.Context(), params.Category)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, categoryResponse{
		Category:    newCategory.Category,
		AccessToken: newAccessToken,
		CreatedAt:   newCategory.CreatedAt,
		UpdatedAt:   newCategory.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleUpdateCategory(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	// request struct
	type updateCategoryRequest struct {
		CategoryID uuid.UUID `json:"categoryID"`
		Category   string    `json:"category"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := updateCategoryRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if params.CategoryID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid category id")
		return
	}
	if params.Category == "" {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid category")
		return
	}

	// updating category
	updatedCategory, err := apiConfig.DB.UpdateCategory(r.Context(), database.UpdateCategoryParams{
		Category: params.Category,
		ID:       params.CategoryID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, categoryResponse{
		Category:    updatedCategory.Category,
		AccessToken: newAccessToken,
		CreatedAt:   updatedCategory.CreatedAt,
		UpdatedAt:   updatedCategory.UpdatedAt,
	})
}

func (apiConfig *ApiConfig) HandleRemoveCategory(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	// request struct
	type removeCategoryRequest struct {
		CategoryID uuid.UUID `json:"categoryID"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := removeCategoryRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if params.CategoryID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	// removing category
	if err = apiConfig.DB.RemoveCategory(r.Context(), params.CategoryID); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

func (apiConfig *ApiConfig) HandleGetAllCategories(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	// response struct
	type AllCategories struct {
		Categories  []database.Category
		AccessToken string `json:"accessToken"`
	}

	allCategories, err := apiConfig.DB.GetAllCategories(r.Context())
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, AllCategories{
		Categories:  allCategories,
		AccessToken: newAccessToken,
	})
}
