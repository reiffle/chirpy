package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestMakeJWTCorrectly(t *testing.T) {
	user := uuid.New()
	token := "MyPassword"
	_, err := MakeJWT(user, token)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	user := uuid.New()
	token := "MyPassword"
	NewJWT, err := MakeJWT(user, token)
	if err != nil {
		t.Errorf("Couldn't make JWT: %v", err)
	}
	NewUuid, err := ValidateJWT(NewJWT, token)
	if err != nil {
		t.Errorf("Couldn't validate JWT: %v", err)
	}
	if NewUuid != user {
		t.Errorf("JWT does not match user: %v", err)
	}
}

func TestValidateJWTWrongPassword(t *testing.T) {
	user := uuid.New()
	token := "MyPassword"
	NewJWT, err := MakeJWT(user, token)
	if err != nil {
		t.Errorf("Couldn't make JWT: %v", err)
	}
	_, err = ValidateJWT(NewJWT, token+"s")
	if err == nil {
		t.Errorf("Wrong Password Allowed Access: %v", err)
	}
}

func TestGetBearerToken(t *testing.T) {
	pwd := "Bearer this is a password"
	pwd2 := "this is a password"
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Errorf("Error creating request: %v", err)
		return
	}

	// Set a single header
	req.Header.Set("Authorization", pwd)
	result, err := GetBearerToken(req.Header)
	if err != nil {
		t.Errorf("Couldn't get token: %v", err)
	}
	if result != pwd2 {
		t.Error("Token and password do not match")
	}
}
