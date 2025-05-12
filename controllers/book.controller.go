package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
)

// admin
func (apiConfig *ApiConfig) HandleAddBook(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		Name          string `json:"name"`
		CoverImageURL string `json:"coverImageUrl"`
		Review        string `json:"review"`
		Tags          string `json:"tags"`
		Level         string `json:"level"`
	}

	type Response struct {
		ID          uuid.UUID `json:"id"`
		Book        Request   `json:"book"`
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

	// validating all the request params
	if err = apiConfig.DataValidator.Var(params.Name, "required,bookname"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	if err = apiConfig.DataValidator.Var(params.CoverImageURL, "required,url"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	if err = apiConfig.DataValidator.Var(params.Review, "required,min=30,max=200"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	if err = apiConfig.DataValidator.Var(params.Tags, "required,tags"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	// adding new book
	levelID, err := apiConfig.DB.GetLevelIDByName(r.Context(), params.Level)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	newBook, err := apiConfig.DB.CreateBook(r.Context(), database.CreateBookParams{
		Name:          params.Name,
		CoverImageUrl: params.CoverImageURL,
		Review:        params.Review,
		Tags:          params.Tags,
		Level:         levelID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, Response{
		ID: newBook.ID,
		Book: Request{
			Name:          newBook.Name,
			CoverImageURL: newBook.CoverImageUrl,
			Review:        newBook.Review,
			Tags:          newBook.Tags,
			Level:         params.Level,
		},
		AccessToken: newAccessToken,
	})
}

// admin
func (apiConfig *ApiConfig) HandleUpdateBook(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		ID            uuid.UUID `json:"id"`
		Name          string    `json:"name,omitempty"`
		CoverImageURL string    `json:"coverImageUrl,omitempty"`
		Review        string    `json:"review,omitempty"`
		Tags          string    `json:"tags,omitempty"`
		Level         string    `json:"level,omitempty"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := Request{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// checking which field to update
	existingInformation, err := apiConfig.DB.GetBookByID(r.Context(), params.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	updateBook := database.UpdateBookParams{}
	if apiConfig.DataValidator.Var(params.Name, "required,bookname") == nil {
		updateBook.Name = params.Name
	} else {
		updateBook.Name = existingInformation.Name
	}

	if apiConfig.DataValidator.Var(params.CoverImageURL, "required,url") == nil {
		updateBook.CoverImageUrl = params.CoverImageURL
	} else {
		updateBook.CoverImageUrl = existingInformation.CoverImageUrl
	}

	if apiConfig.DataValidator.Var(params.Review, "required,min=30,max=200") == nil {
		updateBook.Review = params.Review
	} else {
		updateBook.Review = existingInformation.Review
	}

	if apiConfig.DataValidator.Var(params.Tags, "required,tags") == nil {
		updateBook.Tags = params.Tags
	} else {
		updateBook.Tags = existingInformation.Tags
	}

	if params.Level != "" {
		levelID, err := apiConfig.DB.GetLevelIDByName(r.Context(), params.Level)
		if err != nil {
			utility.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		updateBook.Level = levelID
	} else {
		updateBook.Level = existingInformation.Level
	}

	// updating the book
	updateBook.ID = params.ID
	err = apiConfig.DB.UpdateBook(r.Context(), updateBook)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, struct {
		Book        Request `json:"book"`
		AccessToken string  `json:"accessToken"`
	}{
		Book: Request{
			ID:            params.ID,
			Name:          params.Name,
			CoverImageURL: params.CoverImageURL,
			Review:        params.Review,
			Tags:          params.Tags,
			Level:         params.Level,
		},
		AccessToken: newAccessToken,
	})
}

// admin
func (apiConfig *ApiConfig) HandleRemoveBook(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
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

	// removing book
	if err = apiConfig.DB.RemoveBook(r.Context(), params.ID); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

// user
func (apiConfig *ApiConfig) HandleFilterBooksByLevel(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	// extracting level from query params
	level := r.URL.Query().Get("level")
	if level == "" {
		// add a custom validator for validating level values
		utility.RespondWithError(w, http.StatusNotAcceptable, "invalid level")
		return
	}

	// filtering the books according to level
	filteredBooks, err := apiConfig.DB.GetBooksByLevel(r.Context(), level)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type FilteredBooksResponse struct {
		Books       []database.GetBooksByLevelRow `json:"books"`
		AccessToken string                        `json:"accessToken"`
	}

	utility.RespondWithJson(w, http.StatusOK, FilteredBooksResponse{
		Books:       filteredBooks,
		AccessToken: newAccessToken,
	})
}

// user
func (apiConfig *ApiConfig) HandleGetAllBooks(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Response struct {
		Books       []database.GetAllBooksRow `json:"books"`
		AccessToken string                    `json:"accessToken"`
	}

	books, err := apiConfig.DB.GetAllBooks(r.Context())
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, Response{
		Books:       books,
		AccessToken: newAccessToken,
	})
}

// user
func (apiConfig *ApiConfig) HandleGetReviewByBookID(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type ReviewRequest struct {
		ID uuid.UUID `json:"id"`
	}

	type ReviewResposne struct {
		CoverImageURL string `json:"coverImageUrl"`
		Review        string `json:"review"`
		AccessToken   string `json:"accessToken"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := ReviewRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if params.ID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid book id")
		return
	}

	// fetching the review
	review, err := apiConfig.DB.GetReviewByBookID(r.Context(), params.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, ReviewResposne{
		CoverImageURL: review.CoverImageUrl,
		Review:        review.Review,
		AccessToken:   newAccessToken,
	})
}
