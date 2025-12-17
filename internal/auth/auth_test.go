package auth

import (
	"errors"
	"fmt"
	"testing"
)

func setup() {
	RegisteredProviders = make(map[string]bool)
}

func TestCreateUser(t *testing.T) {
	setup()
	userEmail := "user@example.com"
	providerKey := "user@example.com"
	user, err := CreateUser(userEmail, "email", providerKey, "password123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.Email != providerKey {
		t.Errorf("Expected email %s, got %s", userEmail, user.Email)
	}
}

func TestCreateUserMakesUniqueIDs(t *testing.T) {
	setup()
	email := "user@example.com"
	providerType := "email"
	providerKey := "hello"
	credential := "password123"
	userIDs := make(map[string]bool)
	for i := 0; i < 100; i++ {
		uniqueEmail := fmt.Sprintf("%s-%d@example.com", email, i)
		uniqueKey := fmt.Sprintf("%s-%d", providerKey, i)
		user, err := CreateUser(uniqueEmail, providerType, uniqueKey, credential)
		if err != nil {
			t.Fatalf("Error creating user: %v", err)
		}
		if userIDs[user.ID] {
			t.Errorf("Duplicate ID found: %s", user.ID)
		}
		userIDs[user.ID] = true
	}
}

func TestCreateUserForbidsDuplicateEmails(t *testing.T) {
	setup()
	email := "dup@example.com"

	CreateUser(email, "email", "key1", "pass")
	_, err := CreateUser(email, "email", "key2", "pass")

	var dupErr *ErrDuplicateField
	if !errors.As(err, &dupErr) {
		t.Fatalf("Expected ErrDuplicateField, got %T (%v)", err, err)
	}
	if dupErr.Field != "registrationKey" {
		t.Errorf("Expected error on field 'email', got %s", dupErr.Field)
	}
}
func TestCreateUserForbidsDuplicateProviderKeys(t *testing.T) {
	setup()
	providerType := "email"
	credential := "password123"
	duplicateKey := "key-1"
	_, err := CreateUser("email-1@email.com", providerType, duplicateKey, credential)
	_, err = CreateUser("email-2@email.com", providerType, duplicateKey, credential)
	if err == nil {
		t.Errorf("Expected duplicate email registration to raise an error")
	}
}

func TestCreateUserAllowsMultipleProviderTypesPerEmail(t *testing.T) {
	email := "email-1@email.com"
	providers := map[string]string{"A": "1", "B": "2", "C": "3"}
	for pType, pKey := range providers {
		_, err := CreateUser(email, pType, pKey, "password123")
		if err != nil {
			t.Errorf("Failed to register %s: %v", pType, err)
		}
	}
}
