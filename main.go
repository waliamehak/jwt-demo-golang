package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT payload we sign and verify
type Claims struct {
	Username string `json:"waliamehak"`
	jwt.RegisteredClaims
}

var jwtKey []byte

func loadSecret() []byte {
	secret := os.Getenv("JWT_SECRET") // Load from environment variable
	if secret == "" {
		log.Fatalln("JWT_SECRET is not set. Run: export JWT_SECRET=\"your-strong-secret\"")
	}
	return []byte(secret)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// In a real application, you'd verify username and password here
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Decode JSON body from the client into creds
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Now creds.Username and creds.Password contain client-sent data
	if creds.Username != "waliamehak" || creds.Password != "password123" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   creds.Username,
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// Return the signed token
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	// Get the token from the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		authHeader = authHeader[len("Bearer "):]
	}

	claims := &Claims{}

	// Parse and verify the token
	token, err := jwt.ParseWithClaims(authHeader, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Token is valid; return welcome message
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

// corsMiddleware adds CORS headers so React (localhost:3000) can call the API
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		// Handle preflight request quickly
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	jwtKey = loadSecret()
	http.Handle("/login", corsMiddleware(http.HandlerFunc(loginHandler)))
	http.Handle("/welcome", corsMiddleware(http.HandlerFunc(welcomeHandler)))

	port := "8080"
	fmt.Println("Server starting at :", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
