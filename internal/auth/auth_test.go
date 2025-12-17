package auth

import (
	"fmt"
	"testing"
)

func setup() {
	RegisteredUsers = make(map[string]bool)
}

func TestCreateUser(t *testing.T) {
	setup()
	providerType := "email"
	providerKey := "user@example.com"
	credential := "password123"
	user, err := CreateUser(providerKey, providerType, providerKey, credential)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.Email != providerKey {
		t.Errorf("Expected email %s, got %s", providerKey, user.Email)
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
	providerType := "email"
	providerKey := "user@example.com"
	credential := "password123"

	_, err := CreateUser(providerKey, providerType, providerKey, credential)
	_, err = CreateUser(providerKey, providerType, providerKey, credential)
	if err == nil {
		t.Errorf("Expected duplicate email registration to raise an error")
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
