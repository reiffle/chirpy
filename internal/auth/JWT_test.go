package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTCorrectly(t *testing.T) {
	user := uuid.New()
	token := "MyPassword"
	duration := time.Minute
	_, err := MakeJWT(user, token, duration)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
}

func TestValidateJWT(t *testing.T) {
	user := uuid.New()
	token := "MyPassword"
	duration := time.Minute
	NewJWT, err := MakeJWT(user, token, duration)
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
	duration := time.Minute
	NewJWT, err := MakeJWT(user, token, duration)
	if err != nil {
		t.Errorf("Couldn't make JWT: %v", err)
	}
	_, err = ValidateJWT(NewJWT, token+"s")
	if err == nil {
		t.Errorf("Wrong Password Allowed Access: %v", err)
	}
}
