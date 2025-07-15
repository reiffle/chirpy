package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Issuer:    "chirpy",
		Subject:   userID.String(),
	}

	// Create a new token with the claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	secretKey := []byte(tokenSecret) // Use a strong, secure key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error signing token:", err)
		return "", err
	}

	return tokenString, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key (for HS256)
		return []byte(tokenSecret), nil
	})

	if err != nil {
		// Handle parsing or validation errors
		return uuid.Nil, err
	}

	if token.Valid {
		newUUID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, fmt.Errorf("Invalid user ID: %w", err)
		}
		// Access your custom claims
		fmt.Printf("Subject claim: %s\n", claims.Subject)
		// Access registered claims
		fmt.Printf("Expiration time: %v\n", claims.ExpiresAt)
		return newUUID, nil
	} else {
		fmt.Println("Invalid token")
		return uuid.Nil, err
	}
}
