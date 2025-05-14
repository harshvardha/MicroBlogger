package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// admin
func (apiConfig *ApiConfig) HandleCreateBlog(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		Title        string            `json:"title"`
		Brief        string            `json:"brief,omitempty"`
		ContentURL   string            `json:"contentUrl,omitempty"`
		Images       map[string]string `json:"images,omitempty"`
		ThumbnailURL string            `json:"thumbnailUrl"`
		CodeRepoLink sql.NullString    `json:"codeRepoLink,omitempty"`
		Tags         []string          `json:"tags"`
		Category     string            `json:"category"`
	}

	type Response struct {
		ID          uuid.UUID `json:"id"`
		Blog        Request   `json:"blog"`
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

	// validating request params
	if err = apiConfig.DataValidator.Var(params.Title, "required,min=5,max=70"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	if err = apiConfig.DataValidator.Var(params.Brief, "required,max=200"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	if err = apiConfig.DataValidator.Var(params.ContentURL, "required,url"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	// validating and creating json data for image urls
	var imagesJson []byte
	if params.Images != nil {
		for _, url := range params.Images {
			if err = apiConfig.DataValidator.Var(url, "required,url"); err != nil {
				utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
				return
			}
		}
		imagesJson, err = json.Marshal(params.Images)
		if err != nil {
			utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		utility.RespondWithError(w, http.StatusNotAcceptable, "invalid image urls")
		return
	}

	if err = apiConfig.DataValidator.Var(params.ThumbnailURL, "required,url"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	if err = apiConfig.DataValidator.Var(params.CodeRepoLink, "required,github_url"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	if err = apiConfig.DataValidator.Var(params.Tags, "required,min=1,tags"); err != nil {
		utility.RespondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	// creating new blog
	categoryID, err := apiConfig.DB.GetCategoryIDByName(r.Context(), cases.Title(language.English).String(params.Category))
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	newBlog, err := apiConfig.DB.CreateBlog(r.Context(), database.CreateBlogParams{
		Title:        params.Title,
		Brief:        params.Brief,
		ContentUrl:   params.ContentURL,
		Images:       imagesJson,
		ThumbnailUrl: params.ThumbnailURL,
		CodeRepoLink: params.CodeRepoLink,
		Tags:         params.Tags,
		Author:       IDAndRole.ID,
		Category:     categoryID,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusCreated, Response{
		ID: newBlog.ID,
		Blog: Request{
			Title:        newBlog.Title,
			Brief:        newBlog.Brief,
			ThumbnailURL: newBlog.ThumbnailUrl,
			Tags:         newBlog.Tags,
			Category:     params.Category,
		},
		CreatedAt:   newBlog.CreatedAt,
		UpdatedAt:   newBlog.UpdatedAt,
		AccessToken: newAccessToken,
	})
}

// admin
func (apiConfig *ApiConfig) HandleUpdateBlog(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		ID           uuid.UUID         `json:"id"`
		Title        string            `json:"title"`
		Brief        string            `json:"brief"`
		ContentURL   string            `json:"contentUrl"`
		Images       map[string]string `json:"images"`
		ThumbnailURL string            `json:"thumbnailUrl"`
		CodeRepoLink sql.NullString    `json:"codeRepoLink"`
		Tags         []string          `json:"tags"`
	}

	type Response struct {
		Title        string            `json:"title,omitempty"`
		ContentURL   string            `json:"contentUrl,omitempty"`
		Images       map[string]string `json:"images,omitempty"`
		ThumbnailURL string            `json:"thumbnailUrl,omitempty"`
		CodeRepoLink sql.NullString    `json:"codeRepoLink,omitempty"`
		Tags         []string          `json:"tags,omitempty"`
		CreatedAt    time.Time         `json:"createdAt"`
		UpdatedAt    time.Time         `json:"updatedAt"`
		AccessToken  string            `json:"accessToken"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := Request{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// use custom validation to validate id field then fetch the existing information
	// checking which fields to update
	existingInformation, err := apiConfig.DB.GetBlogByID(r.Context(), params.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// validate the param fields which are updatable
	updateBlog := database.UpdateBlogParams{}
	updateBlog.ID = params.ID

	// validating title
	if apiConfig.DataValidator.Var(params.Title, "required,min=5,max=70") == nil {
		updateBlog.Title = params.Title
	} else {
		updateBlog.Title = existingInformation.Title
	}

	// validating brief
	if apiConfig.DataValidator.Var(params.Brief, "required,max=200") == nil {
		updateBlog.Brief = params.Brief
	} else {
		updateBlog.Brief = existingInformation.Brief
	}

	// validating content url
	if apiConfig.DataValidator.Var(params.ContentURL, "required,url") == nil {
		updateBlog.ContentUrl = params.ContentURL
	} else {
		updateBlog.ContentUrl = existingInformation.ContentUrl
	}

	// validating images url
	if params.Images != nil {
		for _, url := range params.Images {
			if err = apiConfig.DataValidator.Var(url, "required,url"); err != nil {
				utility.RespondWithError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		images, err := json.Marshal(params.Images)
		if err != nil {
			utility.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		updateBlog.Images = images
	} else {
		updateBlog.Images = existingInformation.Images
	}

	// validating thumbnail url
	if apiConfig.DataValidator.Var(params.ThumbnailURL, "required,url") == nil {
		updateBlog.ThumbnailUrl = params.ThumbnailURL
	} else {
		updateBlog.ThumbnailUrl = existingInformation.ThumbnailUrl
	}

	// validating code repo link
	if apiConfig.DataValidator.Var(params.CodeRepoLink, "required,github_url") == nil {
		updateBlog.CodeRepoLink = params.CodeRepoLink
	} else {
		updateBlog.CodeRepoLink = existingInformation.CodeRepoLink
	}

	// validating tags
	if apiConfig.DataValidator.Var(params.Tags, "required,min=1,tags") == nil {
		updateBlog.Tags = params.Tags
	} else {
		updateBlog.Tags = existingInformation.Tags
	}

	// updating blogs
	updatedBlog, err := apiConfig.DB.UpdateBlog(r.Context(), updateBlog)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var updatedImages map[string]string
	if err = json.Unmarshal(updateBlog.Images, &updatedImages); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utility.RespondWithJson(w, http.StatusOK, Response{
		Title:        updateBlog.Title,
		ContentURL:   updateBlog.ContentUrl,
		Images:       updatedImages,
		ThumbnailURL: updateBlog.ThumbnailUrl,
		CodeRepoLink: updateBlog.CodeRepoLink,
		Tags:         updateBlog.Tags,
		CreatedAt:    updatedBlog.CreatedAt,
		UpdatedAt:    updatedBlog.UpdatedAt,
		AccessToken:  newAccessToken,
	})
}

// admin
func (apiConfig *ApiConfig) HandleRemoveBlog(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type RemoveBlogRequest struct {
		ID uuid.UUID `json:"id"`
	}

	// decoding request body
	decoder := json.NewDecoder(r.Body)
	params := RemoveBlogRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if params.ID == uuid.Nil {
		utility.RespondWithError(w, http.StatusBadRequest, "invalid blog id")
		return
	}

	// removing blog
	if err = apiConfig.DB.RemoveBlog(r.Context(), params.ID); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, EmptyResponse{
		AccessToken: newAccessToken,
	})
}

// both
func (apiConfig *ApiConfig) HandleGetBlogsByCategory(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		CreatedAt time.Time `json:"createdAt"`
		Category  string    `json:"category"`
		Limit     int32     `json:"limit"`
	}

	type Response struct {
		Blogs       []database.GetAllBlogsByCategoryRow `json:"blogs"`
		AccessToken string                              `json:"accessToken"`
	}

	// decoding the request body
	decoder := json.NewDecoder(r.Body)
	params := Request{}
	err := decoder.Decode(&params)
	if err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fetching all the blogs
	categoryID, err := apiConfig.DB.GetCategoryIDByName(r.Context(), params.Category)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	blogs, err := apiConfig.DB.GetAllBlogsByCategory(r.Context(), database.GetAllBlogsByCategoryParams{
		Category:  categoryID,
		CreatedAt: params.CreatedAt,
		Limit:     params.Limit,
	})
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, Response{
		Blogs:       blogs,
		AccessToken: newAccessToken,
	})
}

// both
func (apiConfig *ApiConfig) HandleFilterBlogs(w http.ResponseWriter, r *http.Request, newAccessToken string) {

}

// both
func (apiConfig *ApiConfig) HandleGetBlogByID(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		ID uuid.UUID `json:"id"`
	}

	type Response struct {
		Title        string            `json:"title"`
		ContentURL   string            `json:"contentUrl"`
		Images       map[string]string `json:"images"`
		CodeRepoLink sql.NullString    `json:"codeRepoLink,omitempty"`
		Views        int32             `json:"views"`
		Likes        int64             `json:"likes"`
		Tags         []string          `json:"tags"`
		Author       string            `json:"author"`
		CreatedAt    time.Time         `json:"createdAt"`
		HasUserLiked bool              `json:"hasUserLiked"`
		AccessToken  string            `json:"accessToken"`
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

	blog, err := apiConfig.DB.GetBlogByID(r.Context(), params.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var images map[string]string
	if err = json.Unmarshal(blog.Images, &images); err != nil {
		utility.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// fetching number of likes for the blog and checking if user has liked the blog or not
	noOfLikes, err := apiConfig.DB.GetNumberOfLikes(r.Context(), params.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	hasUserLikedThisBlog, err := apiConfig.DB.HasUserLikedBlog(r.Context(), database.HasUserLikedBlogParams{
		UserID: IDAndRole.ID,
		BlogID: params.ID,
	})
	if err != nil && err != sql.ErrNoRows {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response := Response{
		Title:        blog.Title,
		ContentURL:   blog.ContentUrl,
		Images:       images,
		CodeRepoLink: blog.CodeRepoLink,
		Views:        blog.Views,
		Likes:        noOfLikes,
		Tags:         blog.Tags,
		Author:       blog.Username,
		CreatedAt:    blog.CreatedAt,
		AccessToken:  newAccessToken,
	}
	if hasUserLikedThisBlog == 1 {
		response.HasUserLiked = true
	} else {
		response.HasUserLiked = false
	}

	utility.RespondWithJson(w, http.StatusOK, response)
}

// user
func (apiConfig *ApiConfig) HandleLikeOrDislike(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		ID uuid.UUID `json:"id"`
	}

	type Response struct {
		LikesCount  int64  `json:"likesCount"`
		AccessToken string `json:"accessToken"`
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

	// checking if blog is liked then disliking it
	isBlogLiked, err := apiConfig.DB.HasUserLikedBlog(r.Context(), database.HasUserLikedBlogParams{
		UserID: IDAndRole.ID,
		BlogID: params.ID,
	})

	if err != nil && err != sql.ErrNoRows {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if isBlogLiked == 1 {
		if err = apiConfig.DB.DislikeBlog(r.Context(), database.DislikeBlogParams{
			UserID: IDAndRole.ID,
			BlogID: params.ID,
		}); err != nil {
			utility.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	} else {
		if err = apiConfig.DB.LikeBlog(r.Context(), database.LikeBlogParams{
			UserID: IDAndRole.ID,
			BlogID: params.ID,
		}); err != nil {
			utility.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	// getting likes count
	likesCount, err := apiConfig.DB.GetNumberOfLikes(r.Context(), params.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, Response{
		LikesCount:  likesCount,
		AccessToken: newAccessToken,
	})
}

// user
func (apiConfig *ApiConfig) HandleIncrementView(w http.ResponseWriter, r *http.Request, IDAndRole *IDAndRole, newAccessToken string) {
	type Request struct {
		ID uuid.UUID `json:"id"`
	}

	type Response struct {
		Views       int32  `json:"views"`
		AccessToken string `json:"accessToken"`
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

	if err = apiConfig.DB.IncrementViews(r.Context(), params.ID); err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	views, err := apiConfig.DB.GetViewCount(r.Context(), params.ID)
	if err != nil {
		utility.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	utility.RespondWithJson(w, http.StatusOK, Response{
		Views:       views,
		AccessToken: newAccessToken,
	})
}
