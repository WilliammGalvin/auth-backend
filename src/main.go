package main

import (
	"backend/data"
	"backend/routes/auth"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

// protectedEndpoint ensures API key is needed to interact with the API
func protectedEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Api-Key")

		if apiKey == "" || apiKey != os.Getenv("API_SECRET_KEY") {
			http.Error(w, "Invalid or missing API key.", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// initEnv initializes the environment variables
func initEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	initEnv()
	data.InitDB()

	defer func() {
		if err := data.CloseDB(); err != nil {
			log.Printf("Error closing database connection: %v\n", err)
		}
	}()

	mux := http.NewServeMux()

	mux.Handle("/api/auth/validate-refresh", protectedEndpoint(http.HandlerFunc(auth.ValidateRefreshToken)))
	mux.Handle("/api/auth/refresh-token", protectedEndpoint(http.HandlerFunc(auth.RefreshAuthToken)))
	mux.Handle("/api/auth/users", protectedEndpoint(http.HandlerFunc(auth.CreateUser)))
	mux.Handle("/api/auth/verify", protectedEndpoint(http.HandlerFunc(auth.VerifyUser)))

	fmt.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
