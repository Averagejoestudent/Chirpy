package auth
import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTFlow(t *testing.T) {
	secret := "my-super-secret-key"
	userID := uuid.New()
	duration := time.Hour

	// 1. Test Creation
	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}
	if token == "" {
		t.Fatal("Token should not be empty")
	}

	// 2. Test Validation (Happy Path)
	parsedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if parsedID != userID {
		t.Errorf("Expected userID %v, got %v", userID, parsedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	secret := "secret"
	userID := uuid.New()
	// Create a token that expired 1 hour ago
	duration := -time.Hour 

	token, err := MakeJWT(userID, secret, duration)
	if err != nil {
		t.Fatalf("Error creating token: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Expected error for expired token, but got nil")
	}
}

func TestInvalidSecret(t *testing.T) {
	secret := "correct-secret"
	wrongSecret := "wrong-secret"
	userID := uuid.New()

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("Expected error when validating with wrong secret, but got nil")
	}
}