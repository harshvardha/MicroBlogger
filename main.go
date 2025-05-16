package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/harshvardha/artOfSoftwareEngineering/controllers"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/cache"
	"github.com/harshvardha/artOfSoftwareEngineering/internal/database"
	"github.com/harshvardha/artOfSoftwareEngineering/middlewares"
	"github.com/harshvardha/artOfSoftwareEngineering/utility"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// loading all the required env variables
	godotenv.Load()

	// loading jwt access secret key
	jwtSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if jwtSecret == "" {
		log.Fatal("jwt_secret variable not set")
	}

	// loading port number
	portNo := os.Getenv("PORT")
	if portNo == "" {
		log.Fatal("port variable not set")
	}

	// loading db uri
	dbUri := os.Getenv("DATABASE_URI")
	if dbUri == "" {
		log.Fatal("Database URI not set")
	}

	// loading otp cache required configs
	fromEmail := os.Getenv("FROM_EMAIL")
	if fromEmail == "" {
		log.Fatal("From Email not set")
	}
	emailSubject := os.Getenv("EMAIL_SUBJECT")
	if emailSubject == "" {
		log.Fatal("Email Subject not set")
	}
	emailBody := os.Getenv("EMAIL_BODY")
	if emailBody == "" {
		log.Fatal("Email Body not set")
	}
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		log.Fatal("SMTP HOST not set")
	}
	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		log.Fatal("SMTP PORT not set")
	}
	appPassword := os.Getenv("GMAIL_APP_PASSWORD")
	if appPassword == "" {
		log.Fatal("App Password not set")
	}

	// creating database connection
	dbConnection, err := sql.Open("postgres", dbUri)
	if err != nil {
		log.Fatal("Error Connecting to Database: ", err)
	}

	// setting data validator
	dataValidator := validator.New()

	// registering new custom password validator
	dataValidator.RegisterValidation("password", utility.CustomPasswordValidator)

	// registering new custom github url validator
	dataValidator.RegisterValidation("github_url", utility.GithubURLValidator)

	// registering new custom username validator
	dataValidator.RegisterValidation("username", utility.UsernameValidator)

	// registering new custom bookName validator
	dataValidator.RegisterValidation("bookname", utility.BookNameValidator)

	// registering new tags validator
	dataValidator.RegisterValidation("tags", utility.NoDuplicatesTagsValidator)

	// setting apiConfig
	apiConfig := controllers.ApiConfig{
		DB:            database.New(dbConnection),
		JwtSecret:     jwtSecret,
		OtpCache:      cache.NewOTPCache(fromEmail, emailSubject, emailBody, smtpHost, smtpPort, appPassword),
		DataValidator: dataValidator,
	}

	routes := map[string][]string{
		"user": {
			"/api/v1/book/filter",
			"/api/v1/book/all",
			"/api/v1/book/review",
			"/api/v1/user/update/email",
			"/api/v1/user/update/password",
			"/api/v1/user/update/other",
			"/api/v1/user/account/remove",
			"/api/v1/user/search",
			"/api/v1/user",
			"/api/v1/blog/view/increment",
			"/api/v1/blog/likedislike",
			"/api/v1/blog",
			"/api/v1/blog/category",
			"/api/v1/comment/create",
			"/api/v1/comment/update",
			"/api/v1/comment/remove",
			"/api/v1/comment/all",
		},
		"nil_IDAndRole": {
			"/api/v1/book/add",
			"/api/v1/book/review",
			"/api/v1/book/all",
			"/api/v1/book/filter",
			"/api/v1/book/update",
			"/api/v1/book/remove",
			"/api/v1/blog/update",
			"/api/v1/blog/remove",
			"/api/v1/blog/category",
			"/api/v1/blog/view/increment",
			"/api/v1/comment/all",
		},
		"nil_accessToken": {
			"/api/v1/user/update/email",
			"/api/v1/user/update/password",
			"/api/v1/user/account/remove",
		},
	}

	// creating new request redirecting multiplexer
	mux := http.NewServeMux()

	// creating a healthz endpoint for server status check
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// api endpoints for authentication
	mux.HandleFunc("GET /api/v1/auth/otp/send", apiConfig.HandleSendOTP)
	mux.HandleFunc("POST /api/v1/auth/register", apiConfig.HandleRegisterUser)
	mux.HandleFunc("GET /api/v1/auth/otp/resend", apiConfig.HandleResendOTP)
	mux.HandleFunc("POST /api/v1/auth/login", apiConfig.HandleLogin)

	// api endpoints for books
	mux.HandleFunc("POST /api/v1/book/add", middlewares.ValidateJWT(apiConfig.HandleAddBook, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("PUT /api/v1/book/update", middlewares.ValidateJWT(apiConfig.HandleUpdateBook, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("DELETE /api/v1/book/remove", middlewares.ValidateJWT(apiConfig.HandleRemoveBook, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("GET /api/v1/book/filter", middlewares.ValidateJWT(apiConfig.HandleFilterBooksByLevel, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("GET /api/v1/book/all", middlewares.ValidateJWT(apiConfig.HandleGetAllBooks, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("GET /api/v1/book/review", middlewares.ValidateJWT(apiConfig.HandleGetReviewByBookID, apiConfig.JwtSecret, apiConfig.DB, routes))

	// api endpoints for user
	mux.HandleFunc("PUT /api/v1/user/update/email", middlewares.ValidateJWT(apiConfig.HandleUpdateEmail, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("PUT /api/v1/user/update/password", middlewares.ValidateJWT(apiConfig.HandleUpdatePassword, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("PUT /api/v1/user/update/other", middlewares.ValidateJWT(apiConfig.HandleUpdateOtherDetails, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("DELETE /api/v1/user/account/remove", middlewares.ValidateJWT(apiConfig.HandleRemoveUserAccount, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("GET /api/v1/user", middlewares.ValidateJWT(apiConfig.HandleGetUserByID, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("GET /api/v1/user/search", middlewares.ValidateJWT(apiConfig.HandleUserSearch, apiConfig.JwtSecret, apiConfig.DB, routes))

	// api endpoints for category
	mux.HandleFunc("POST /api/v1/category/create", middlewares.ValidateJWT(apiConfig.HandleCreateCategory, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("PUT /api/v1/category/update", middlewares.ValidateJWT(apiConfig.HandleUpdateCategory, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("DELETE /api/v1/category/remove", middlewares.ValidateJWT(apiConfig.HandleRemoveCategory, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("GET /api/v1/category/all", middlewares.ValidateJWT(apiConfig.HandleGetAllCategories, apiConfig.JwtSecret, apiConfig.DB, routes))

	// api endpoints for blogs
	mux.HandleFunc("POST /api/v1/blog/create", middlewares.ValidateJWT(apiConfig.HandleCreateBlog, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("PUT /api/v1/blog/update", middlewares.ValidateJWT(apiConfig.HandleUpdateBlog, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("DELETE /api/v1/blog/remove", middlewares.ValidateJWT(apiConfig.HandleRemoveBlog, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("GET /api/v1/blog/category", middlewares.ValidateJWT(apiConfig.HandleGetBlogsByCategory, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("GET /api/v1/blog", middlewares.ValidateJWT(apiConfig.HandleGetBlogByID, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("PUT /api/v1/blog/likedislike", middlewares.ValidateJWT(apiConfig.HandleLikeOrDislike, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("PUT /api/v1/blog/views/increment", middlewares.ValidateJWT(apiConfig.HandleIncrementView, apiConfig.JwtSecret, apiConfig.DB, routes))

	// api endpoints for comments
	mux.HandleFunc("POST /api/v1/comment/create", middlewares.ValidateJWT(apiConfig.HandleCreateComment, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("PUT /api/v1/comment/update", middlewares.ValidateJWT(apiConfig.HandleUpdateComment, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("DELETE /api/v1/comment/remove", middlewares.ValidateJWT(apiConfig.HandleRemoveComment, apiConfig.JwtSecret, apiConfig.DB, routes))
	mux.HandleFunc("GET /api/v1/comment/all", middlewares.ValidateJWT(apiConfig.HandleGetAllCommentsByBlogID, apiConfig.JwtSecret, apiConfig.DB, routes))

	// starting the server
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + portNo,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
