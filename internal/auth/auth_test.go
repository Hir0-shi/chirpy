package auth

import (
	"github.com/google/uuid"
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		authorization string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "valid bearer token",
			authorization: "Bearer mytoken123",
			expectedToken: "mytoken123",
			expectError:   false,
		},
		{
			name:          "missing authorization header",
			authorization: "",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "invalid format",
			authorization: "Bearer",
			expectedToken: "",
			expectError:   true,
		},
		{
			name:          "wrong scheme",
			authorization: "Basic mytoken123",
			expectedToken: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.authorization != "" {
				headers.Set("Authorization", tt.authorization)
			}

			token, err := GetBearerToken(headers)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if token != tt.expectedToken {
					t.Fatalf("expected %s, got %s", tt.expectedToken, token)
				}
			}
		})
	}
}
