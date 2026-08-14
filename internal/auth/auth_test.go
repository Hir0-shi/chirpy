package auth

import (
	"github.com/google/uuid"
	"testing"
	"time"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"

	tokenString, err := MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parsedID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if parsedID != userID {
		t.Fatalf("expected %v, got %v", userID, parsedID)
	}
}

func TestValidateJWTExpired(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret"

	tokenString, err := MakeJWT(userID, tokenSecret, -time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestValidateJWTWrongSecret(t *testing.T) {
	userID := uuid.New()

	tokenString, err := MakeJWT(userID, "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = ValidateJWT(tokenString, "wrong-secret")
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}
