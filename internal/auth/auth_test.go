package auth

import (
	"errors"
	"fmt"
	"testing"
)

func password(s string) Credentials {
	return PasswordCredentials{hashedPassword: s}
}

func TestCreateUser(t *testing.T) {
	auth := NewAuthorizationService()
	userEmail := "user@example.com"
	providerKey := "user@example.com"
	user, err := auth.CreateUser(userEmail, "email", providerKey, password("password123"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.Email != providerKey {
		t.Errorf("Expected email %s, got %s", userEmail, user.Email)
	}
}

func TestCreateUserMakesUniqueIDs(t *testing.T) {
	auth := NewAuthorizationService()
	providerType := "email"
	credential := password("password123")
	userIDs := make(map[string]bool)
	for i := 0; i < 100; i++ {
		uniqueEmail := fmt.Sprintf("user-%d@example.com", i)
		user, err := auth.CreateUser(uniqueEmail, providerType, uniqueEmail, credential)
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
	email := "duplicate@example.com"
	auth.CreateUser(email, "email", email, password("pass"))
	_, err := auth.CreateUser(email, "gmail", "fake-key", password("pass"))
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
	credential := password("password123")
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
	auth.CreateUser("email-1@email.com", "firstProviderType", "randomkey", password("credential"))
	providers := map[string]string{"A": "1", "B": "2", "C": "3"}
	for pType, pKey := range providers {
		_, err := auth.CreateUser(email, pType, pKey, password("credential"))
		if err == nil {
			t.Errorf("adding a new provider type using create user function should fail %s: %v", pType, err)
		}
	}
}

func TestCreateUserRequiresNonEmptyFields(t *testing.T) {
	auth := NewAuthorizationService()
	_, err := auth.CreateUser("", "email", "key-1", password("pass"))
	if err == nil {
		t.Error("Expected error for empty email, but user was created")
	}
	_, err = auth.CreateUser("user@ex.com", "email", "", password("pass"))
	if err == nil {
		t.Error("Expected error for empty provider key, but user was created")
	}
}

func TestCreateUserIsAtomic(t *testing.T) {
	auth := NewAuthorizationService()
	pass := password("pass")
	auth.CreateUser("original@ex.com", "gmail", "key-conflict", pass)
	_, err := auth.CreateUser("new-potential-user@ex.com", "gmail", "key-conflict", pass)
	if err == nil {
		t.Fatal("Expected error due to duplicate provider key, but got nil")
	}
	_, err = auth.CreateUser("new-potential-user@ex.com", "gmail", "valid-key", pass)
	if err != nil {
		t.Errorf("Atomicity failure: %v. The email was 'locked' even though registration failed.", err)
	}
}

func TestGetUserRetrievesCreatedUser(t *testing.T) {
	auth := NewAuthorizationService()
	providerType := "email"
	emails := []string{"1@ex.com", "2@email.com", "3@email.com"}
	pass := password("auth-credential")
	for _, userEmail := range emails {
		auth.CreateUser(userEmail, providerType, userEmail, pass)
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
	auth.CreateUser(email, providerType, providerKey, password(credential))
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
	attempt := "my-credentials"
	_, err := auth.Login(providerType, providerKey, attempt)
	if err == nil {
		t.Errorf("Should have been unable to locate non-existant user %s, %s, %s", providerType, providerKey, attempt)
	}
}

func TestCreateUserHandlesInvalidEmailFormat(t *testing.T) {
	auth := NewAuthorizationService()
	email := "invalid-email"
	pass := password("pass")
	_, err := auth.CreateUser(email, "email", "key-conflict", pass)
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

type errorUuidService struct{}

func (g *errorUuidService) Generate() (string, error) {
	return "", errors.New("randomness failed")
}

func TestCreateUserHandlesFailureToCreateRandomIdentifier(t *testing.T) {
	auth := NewAuthorizationService(WithUuidService(&errorUuidService{}))
	pass := password("pass")
	_, err := auth.CreateUser("test@example.com", "email", "key", pass)
	if err == nil {
		t.Error("Expected error when random identifier generation fails, but got nil")
	}
}

func TestIdentityKeyErrorsOnEmptyKey(t *testing.T) {
	auth := NewAuthorizationService()
	pass := password("pass")
	_, err := auth.CreateUser("test@example.com", "google", "", pass)
	if err == nil {
		t.Errorf("Expected ProviderKey validation to error on empty key")
	}
	_, err = auth.GetUserByIdentity("google", "")
	if err == nil {
		t.Errorf("Expected ProviderKey validation to error on empty key")
	}
	_, err = auth.Login("google", "", "credential")
	if err == nil {
		t.Errorf("Expected ProviderKey validation to error on empty key")
	}
}

func TestFindingMissingUserFails(t *testing.T) {
	auth := NewAuthorizationService()
	_, err := auth.findUserById("missing-user")
	if err == nil {
		t.Errorf("Should not be able to locate missing user")
	}
	_, err = auth.findUserById("")
	if err == nil {
		t.Errorf("Should not be able to locate missing user")
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

func TestLinkMultipleProvidersToSameUser(t *testing.T) {
	auth := NewAuthorizationService()
	oldProviderType := "email"
	oldProviderKey := "user@example.com"
	oldProviderCredentials := password("password123")
	newProviderType := "google"
	newProviderKey := "my-google-key"
	newProviderCredentials := password("my-google-token")
	user, _ := auth.CreateUser(oldProviderKey, oldProviderType, oldProviderKey, oldProviderCredentials)
	err := auth.AddIdentity(user.ID, newProviderType, newProviderKey, newProviderCredentials)
	if err != nil {
		t.Errorf("Expected to successfully add a new identity to existing user")
	}
	user2, err := auth.GetUserByIdentity(newProviderType, newProviderKey)
	if err != nil {
		t.Fatalf("Expected tot fetch user by identity")
	} else if user2 == nil {
		t.Fatal("user2 is nil; identity was likely not saved to the registry")
	} else if user2.ID != user.ID {
		t.Errorf("Expected to fetch the same user based on the new identity")
	}
}

func TestAddIdentityFailsForInvalidInput(t *testing.T) {
	auth := NewAuthorizationService()
	pass1 := password("gmail-token")
	pass2 := password("gmail-token-2")
	err := auth.AddIdentity("fakeUserId", "gmail", "gmail-key", pass1)
	if err == nil {
		t.Errorf("Expected to fail to lookup an invalid user when adding new Identity")
	}
	user, err := auth.CreateUser("me-email@gmail.com", "gmail", "gmail-key", pass1)
	fmt.Printf("Trying to add identity\n")
	err = auth.AddIdentity(user.ID, "gmail", "", pass1)
	if err == nil {
		t.Errorf("Expected to fail to add Identity with empty info")
	}
	user, err = auth.CreateUser("second-user@gmail.com", "gmail-2", "gmail-key-2", pass2)
	// Try to use other person's info
	err = auth.AddIdentity(user.ID, "gmail", "gmail-key", pass1)
	if err == nil {
		t.Errorf("Expected to fail when a second user tries to register another user's third-party keys")
	}
}

func TestCreateUserWithPhone(t *testing.T) {
	auth := NewAuthorizationService()

	// We want this to work
	phone := "+11235550123"
	pass := password("pass123")
	_, err := auth.CreateUser("test@ex.com", "phone", phone, pass)
	if err != nil {
		t.Fatalf("Expected phone registration to work, got %v", err)
	}

	// We want this to FAIL (No + sign)
	invalidPhone := "5550123"
	_, err = auth.CreateUser("test2@ex.com", "phone", invalidPhone, pass)
	if err == nil {
		t.Error("Expected error for phone missing +, but it passed")
	}
}

func TestSwapRegistry(t *testing.T) {

	repo1 := InMemoryUserRepository{registry: map[string]*User{}}
	repo2 := InMemoryUserRepository{registry: map[string]*User{}}
	auth := NewAuthorizationService(WithRepository(repo1))
	auth2 := NewAuthorizationService(WithRepository(repo2))
	userEmail := "user@example.com"
	providerKey := "user@example.com"
	user, err := auth.CreateUser(userEmail, "email", providerKey, password("password123"))
	user, err = auth2.findUserById(user.ID)
	if err == nil {
		t.Errorf("Expected to not find any users in second repo")
	}
}
