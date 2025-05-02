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

	// setting apiConfig
	apiConfig := controllers.ApiConfig{
		DB:        database.New(dbConnection),
		JwtSecret: jwtSecret,
		OtpCache:  cache.NewOTPCache(fromEmail, smtpHost, smtpPort, appPassword),
	}

	// setting data validator
	dataValidator := middlewares.Validator{
		Validate: validator.New(),
	}

	// registering new custom password validator
	dataValidator.Validate.RegisterValidation("password", utility.CustomPasswordValidator)

	// creating new request redirecting multiplexer
	mux := http.NewServeMux()

	// creating a healthz endpoint for server status check
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// api endpoints for authentication
	mux.HandleFunc("POST /api/auth/sendOTP", middlewares.AuthValidation(apiConfig.HandleSendOTP, &dataValidator))
	mux.HandleFunc("POST /api/auth/verifyOTP", apiConfig.HandleRegisterUser)
	mux.HandleFunc("GET /api/auth/resendOTP", apiConfig.HandleResendOTP)
	mux.HandleFunc("POST /api/auth/login", middlewares.AuthValidation(apiConfig.HandleLogin, &dataValidator))

	// starting the server
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + portNo,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
