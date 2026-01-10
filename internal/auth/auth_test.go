package auth

import (
	"errors"
	"fmt"
	"testing"
)

func TestCreateUser(t *testing.T) {
	auth := NewAuthorizationService()
	userEmail := "user@example.com"
	providerKey := "user@example.com"
	user, err := auth.CreateUser(userEmail, "email", providerKey, "password123")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.Email != providerKey {
		t.Errorf("Expected email %s, got %s", userEmail, user.Email)
	}
}

func TestCreateUserMakesUniqueIDs(t *testing.T) {
	auth := NewAuthorizationService()
	email := "user@example.com"
	providerType := "email"
	providerKey := "hello"
	credential := "password123"
	userIDs := make(map[string]bool)
	for i := 0; i < 100; i++ {
		uniqueEmail := fmt.Sprintf("%s-%d@example.com", email, i)
		uniqueKey := fmt.Sprintf("%s-%d", providerKey, i)
		user, err := auth.CreateUser(uniqueEmail, providerType, uniqueKey, credential)
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
	auth := NewAuthorizationService()
	email := "dup@example.com"

	auth.CreateUser(email, "email", "key1", "pass")
	_, err := auth.CreateUser(email, "email", "key2", "pass")

	var dupErr *ErrDuplicateField
	if !errors.As(err, &dupErr) {
		t.Fatalf("Expected ErrDuplicateField, got %T (%v)", err, err)
	}
	if dupErr.Field != "email" {
		t.Errorf("Expected error on field 'email', got %s", dupErr.Field)
	}
}
func TestCreateUserForbidsDuplicateProviderKeys(t *testing.T) {
	auth := NewAuthorizationService()
	providerType := "email"
	credential := "password123"
	duplicateKey := "key-1"
	auth.CreateUser("email-1@email.com", providerType, duplicateKey, credential)
	_, err := auth.CreateUser("email-2@email.com", providerType, duplicateKey, credential)
	if err == nil {
		t.Errorf("Expected duplicate email registration to raise an error")
	}
}

func TestCreateUserForbidsAddingNewProviderTypes(t *testing.T) {
	auth := NewAuthorizationService()
	email := "email-1@email.com"
	auth.CreateUser("email-1@email.com", "firstProviderType", "randomkey", "credential")
	providers := map[string]string{"A": "1", "B": "2", "C": "3"}
	for pType, pKey := range providers {
		_, err := auth.CreateUser(email, pType, pKey, "credential")
		if err == nil {
			t.Errorf("adding a new provider type using create user function should fail %s: %v", pType, err)
		}
	}
}

func TestCreateUserRequiresNonEmptyFields(t *testing.T) {
	auth := NewAuthorizationService()
	_, err := auth.CreateUser("", "email", "key-1", "pass")
	if err == nil {
		t.Error("Expected error for empty email, but user was created")
	}
	_, err = auth.CreateUser("user@ex.com", "email", "", "pass")
	if err == nil {
		t.Error("Expected error for empty provider key, but user was created")
	}
}

func TestCreateUserIsAtomic(t *testing.T) {
	auth := NewAuthorizationService()
	auth.CreateUser("original@ex.com", "email", "key-conflict", "pass")
	newEmail := "new-potential-user@ex.com"
	_, err := auth.CreateUser(newEmail, "email", "key-conflict", "pass")
	if err == nil {
		t.Fatal("Expected error due to duplicate provider key, but got nil")
	}
	_, err = auth.CreateUser(newEmail, "email", "valid-key", "pass")
	if err != nil {
		t.Errorf("Atomicity failure: %v. The email was 'locked' even though registration failed.", err)
	}
}

func TestGetUserRetrievesCreatedUser(t *testing.T) {
	auth := NewAuthorizationService()
	providerType := "email"
	emails := []string{"1@ex.com", "2@email.com", "3@email.com"}
	for _, userEmail := range emails {
		auth.CreateUser(userEmail, providerType, userEmail, "auth-credential")
		user, err := auth.GetUser(userEmail)
		if err != nil {
			t.Fatalf("failed to fetch user that was just created")
		}
		if user.Email != userEmail {
			t.Fatalf("fetched the wrong user %s when expecting %s", user.Email, userEmail)
		}
	}
}

func TestLoginReturnsCorrect(t *testing.T) {
	auth := NewAuthorizationService()
	email := "me@email.com"
	providerType := "google"
	providerKey := "my-provider-key"
	credential := "my-credentials"
	auth.CreateUser(email, providerType, providerKey, credential)
	user, err := auth.Login(providerType, providerKey, credential)
	if err != nil {
		t.Errorf("Unable to locate user using %s, %s, %s", providerType, providerKey, credential)
	}
	if user.Email != email {
		t.Errorf("fetched %s when expecting %s", email, user.Email)
	}
	user, err = auth.Login(providerType, providerKey, "invalid-credentials")
	if err == nil {
		t.Errorf("fetched user %s, but expected failure due to invalid credentials", user.Email)
	}
}

func TestLoginHandlesMissingUser(t *testing.T) {
	auth := NewAuthorizationService()
	providerType := "google"
	providerKey := "my-provider-key"
	credential := "my-credentials"
	_, err := auth.Login(providerType, providerKey, credential)
	if err == nil {
		t.Errorf("Should have been unable to locate non-existant user %s, %s, %s", providerType, providerKey, credential)
	}
}

func TestCreateUserHandlesInvalidEmailFormat(t *testing.T) {
	auth := NewAuthorizationService()
	email := "invalid-email"
	_, err := auth.CreateUser(email, "email", "key-conflict", "pass")
	if err == nil {
		t.Errorf("Email address should contain an @, creation with email='%s' should have failed", email)
	}
}

func TestGetUserHandlesMissingUser(t *testing.T) {
	auth := NewAuthorizationService()
	email := "missing-user@gmail.com"
	_, err := auth.GetUser(email)
	if err == nil {
		t.Errorf("Missing email should return an error upon fetch")
	}
}

func TestCreateUserHandlesRandomFailure(t *testing.T) {
	auth := NewAuthorizationService()
	oldReader := randReader
	defer func() { randReader = oldReader }()
	randReader = func(b []byte) (int, error) {
		return 0, errors.New("randomness failed")
	}
	_, err := auth.CreateUser("test@example.com", "email", "key", "pass")
	if err == nil {
		t.Error("Expected error when random identifier generation fails, but got nil")
	}
}

func TestErrDuplicateFieldFormatting(t *testing.T) {
	// 1. Create the error manually
	customErr := &ErrDuplicateField{
		Field: "email",
		Value: "bob@example.com",
	}

	// 2. Test the Error() string output (exercises the code coverage)
	expected := "duplicate email found: bob@example.com"
	if customErr.Error() != expected {
		t.Errorf("Expected string '%s', got '%s'", expected, customErr.Error())
	}

	// 3. Test how it behaves when wrapped (simulating real-world usage)
	wrappedErr := fmt.Errorf("context: %w", customErr)

	var target *ErrDuplicateField
	if !errors.As(wrappedErr, &target) {
		t.Fatal("Failed to recover ErrDuplicateField using errors.As")
	}

	if target.Field != "email" {
		t.Errorf("Expected field 'email', got '%s'", target.Field)
	}
}
