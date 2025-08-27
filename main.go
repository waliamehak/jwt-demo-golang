package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Loading secret from environment and not hardcoding it in the code.
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatalln("JWT_SECRET is not set. Run: export JWT_SECRET=\"your-strong-secret\"")
	}
	jwtKey := []byte(secret)

	// Create a new token object, specifying signing method and the claims
	claims := jwt.RegisteredClaims{
		Subject:   "waliamehak",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	signed, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatalf("could not sign token: %v\n", err)
	}
	fmt.Println("Generated token (copy this to use protected routes):")
	fmt.Println(signed)
	fmt.Println()

	// Parse and validate the token (to be done later)
	parsed := &jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(signed, parsed, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC (safety check)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		fmt.Println("Token validation failed:", err)
		return
	}

	// 5) If we got here, token is valid
	fmt.Println("Token is valid")
	fmt.Println("Subject (username):", parsed.Subject)
	fmt.Println("Expires at:", parsed.ExpiresAt.Time)
}
