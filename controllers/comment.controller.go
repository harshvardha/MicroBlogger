package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
)

// user
func (apiConfig *ApiConfig) HandleCreateComment(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		BlogID      uuid.UUID `json:"blogID"`
		Description string    `json:"description"`
	}

	type Response struct {
		ID          uuid.UUID `json:"id"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
		AccessToken string    `json:"accessToken"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := Request{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if params.BlogID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid blog id")
		return
	}

	if err = apiConfig.DataValidator.Var(params.Description, "required,min=3,max=200"); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "comment is either to small or to large")
		return
	}

	// adding new comment
	newComment, err := apiConfig.DB.CreateComment(r.Context(), database.CreateCommentParams{
		Description: params.Description,
		UserID:      IDAndRole.ID,
		BlogID:      params.BlogID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, Response{
		ID:          newComment.ID,
		Description: newComment.Description,
		CreatedAt:   newComment.CreatedAt,
		UpdatedAt:   newComment.UpdatedAt,
		AccessToken: newAccessToken,
	})
}

// user
func (apiConfig *ApiConfig) HandleUpdateComment(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		ID          uuid.UUID `json:"id"`
		Description string    `json:"description"`
	}

	type Response struct {
		Description string    `json:"description"`
		UpdatedAt   time.Time `json:"updatedAt"`
		AccessToken string    `json:"accessToken"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := Request{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if params.ID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid comment id")
		return
	}

	if err = apiConfig.DataValidator.Var(params.Description, "required,min=3,max=200"); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, "comment is either to small or to large")
		return
	}

	// updating the comment
	updatedComment, err := apiConfig.DB.UpdateCommentByID(r.Context(), database.UpdateCommentByIDParams{
		Description: params.Description,
		ID:          params.ID,
		UserID:      IDAndRole.ID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, Response{
		Description: updatedComment.Description,
		UpdatedAt:   updatedComment.UpdatedAt,
		AccessToken: newAccessToken,
	})
}

// user
func (apiConfig *ApiConfig) HandleRemoveComment(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		ID uuid.UUID `json:"id"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := Request{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if params.ID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid comment id")
		return
	}

	// removing comment
	if err = apiConfig.DB.RemoveComment(r.Context(), database.RemoveCommentParams{
		ID:     params.ID,
		UserID: IDAndRole.ID,
	}); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

// user
func (apiConfig *ApiConfig) HandleGetAllCommentsByBlogID(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		ID uuid.UUID `json:"id"`
	}

	type Response struct {
		Comments    []database.GetCommentByBlogIDRow `json:"comments"`
		AccessToken string                           `json:"accessToken"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := Request{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if params.ID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid blog id")
		return
	}

	allComments, err := apiConfig.DB.GetCommentByBlogID(r.Context(), params.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, Response{
		Comments:    allComments,
		AccessToken: newAccessToken,
	})
}
