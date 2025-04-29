package main

import (
	"log"
	"net/http"
	"os"

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

	// apiConfig := ApiConfig{
	// 	JwtSecret: jwtSecret,
	// }

	// creating new request redirecting multiplexer
	mux := http.NewServeMux()

	// creating a healthz endpoint for server status check
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// starting the server
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + portNo,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Unable to start server: ", err)
	}
}
